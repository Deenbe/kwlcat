package cmd

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
)

type Setup func() (http.Handler, error)

func RunServerWithProfiler(setup Setup) error {
	fd, err := os.Create("cpu.prof")
	if err != nil {
		return errors.WithStack(err)
	}
	pprof.StartCPUProfile(fd)
	defer pprof.StopCPUProfile()

	ctx, cancelFn := context.WithCancel(context.Background())
	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, os.Interrupt)

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", host, port),
	}
	go func() {
		h, err := setup()
		if err != nil {
			log.Fatal(err)
		}
		server.Handler = h
		log.Printf("listening on %s", server.Addr)
		server.ListenAndServe()
	}()

	select {
	case <-chanSignal:
	}
	cancelFn()
	server.Shutdown(ctx)
	return nil
}

