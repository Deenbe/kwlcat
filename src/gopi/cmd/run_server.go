package cmd

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mattn/go-colorable"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
)

type Setup func(*mux.Router) error

func getLogger() *zap.Logger {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(colorable.NewColorableStdout()), zapcore.DebugLevel), zap.WithCaller(true))
}

func RunServerWithProfiler(name string, setup Setup) error {
	logger := getLogger()

	if profile {
		logger.Sugar().Debug("starting with profiler")
		fd, err := os.Create("cpu.prof")
		if err != nil {
			return errors.WithStack(err)
		}
		pprof.StartCPUProfile(fd)
		defer pprof.StopCPUProfile()
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	chanSignal := make(chan os.Signal, 1)
	done := make(chan error)
	signal.Notify(chanSignal, os.Interrupt)

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
	}
	r := mux.NewRouter()
	r.Use(otelmux.Middleware(name))

	err := setup(r)
	if err != nil {
		return err
	}

	const defaultHealthCheckPath string = "/healthcheck"
	if r.GetRoute(defaultHealthCheckPath) == nil {
		logger.Sugar().Infof("default healthcheck is configured on %s", defaultHealthCheckPath)
		r.PathPrefix(defaultHealthCheckPath).
			Methods("GET").
			HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(200)
				res.Write([]byte("healthy\n"))
			})
	}

	h := handlers.LoggingHandler(os.Stdout, r)
	server.Handler = h

	go func() {
		logger.Sugar().Infof("listening on %v", server.Addr)
		done <- server.ListenAndServe()
	}()

	select {
	case <-chanSignal:
		logger.Sugar().Info("shutting down")
		server.Shutdown(ctx)
	case err = <-done:
	}

	cancelFn()
	return err
}
