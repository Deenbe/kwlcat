package cmd

import (
	"encoding/json"
	"github.com/buddyspike/awsap/dynamodb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"net/http"
	"os"
	"otsample/global"
)

var (
	IDGenCmd = &cobra.Command{
		Use:   "idgen",
		Short: "Run id generation api",
		Long:  "Run id generation api",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunServerWithProfiler(setupIDGenServer)
		},
	}
)

func setupIDGenServer() (http.Handler, error) {
	global.InitialiseTrace("idgen")
	monotonicIDGenerator, err := dynamodb.NewMonotonicIDGenerator("ids", dynamodb.WithTraceProvider(global.TracerProvider))
	if err != nil {
		return nil, err
	}
	r := mux.NewRouter()
	r.Use(otelmux.Middleware("idgen"))
	r.PathPrefix("/ids/next").
		Methods("GET").
		HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			id, err := monotonicIDGenerator.Generate(req.Context(), "scope-a")
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			if err = json.NewEncoder(res).Encode(id); err != nil {
				res.WriteHeader(http.StatusInternalServerError)
			}
	})

	rtr := handlers.LoggingHandler(os.Stdout, r)
	return rtr, nil
}
