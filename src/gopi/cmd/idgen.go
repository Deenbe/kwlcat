package cmd

import (
	"encoding/json"
	"github.com/buddyspike/awsap/dynamodb"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"gopi/global"
	"net/http"
)

var (
	IDGenCmd = &cobra.Command{
		Use:   "idgen",
		Short: "Run id generation api",
		Long:  "Run id generation api",
		RunE: func(cmd *cobra.Command, args []string) error {
			const name string = "idgen"
			global.InitialiseTrace(name)
			return RunServerWithProfiler(name, setupIDGenServer)
		},
	}
)

func setupIDGenServer(r *mux.Router) error {
	monotonicIDGenerator, err := dynamodb.NewMonotonicIDGenerator("ids", dynamodb.WithTraceProvider(global.TracerProvider))
	if err != nil {
		return err
	}
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

	return nil
}
