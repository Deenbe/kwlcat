package cmd

import (
	"encoding/json"
	"github.com/buddyspike/awsap/dynamodb"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopi/global"
	"net/http"
	"os"
)

var (
	idGenTableName string
	IDGenCmd = &cobra.Command{
		Use:   "idgen",
		Short: "Run id generation api",
		Long:  "Run id generation api",
		RunE: func(cmd *cobra.Command, args []string) error {
			if idGenTableName == "" {
				var ok bool
				idGenTableName, ok = os.LookupEnv("IDGEN_NAME")
				if !ok {
					return errors.Errorf("idgen-table-name must be specified")
				}
			}
			const name string = "idgen"
			logger := getLogger()
			global.InitialiseTrace(name, logger)
			return RunServerWithProfiler(name, setupIDGenServer, logger)
		},
	}
)

func init() {
	IDGenCmd.Flags().StringVar(&idGenTableName, "idgen-table-name", "", "dynamodb table name")
}

func setupIDGenServer(r *mux.Router) error {
	monotonicIDGenerator, err := dynamodb.NewMonotonicIDGenerator(idGenTableName, dynamodb.WithTraceProvider(global.TracerProvider))
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
