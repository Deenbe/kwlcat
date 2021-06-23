package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
	"gopi/global"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	NameGenCmd = &cobra.Command{
		Use:   "namegen",
		Short: "Run name generation api",
		Long:  "Run name generation api",
		RunE: func(cmd *cobra.Command, args []string) error {
			const name string = "namegen"
			if idGenApiBaseUrl == "" {
				serviceDiscoveryEndpoint, ok := os.LookupEnv("COPILOT_SERVICE_DISCOVERY_ENDPOINT")
				if !ok {
					return errors.Errorf("idgen-api-base-url must be specified")
				}
				if idGenApiPort == 0 {
					return errors.Errorf("idgen-api-port must be specified")
				}

				idGenApiBaseUrl = fmt.Sprintf("http://idgen.%s:%d", serviceDiscoveryEndpoint, idGenApiPort)
			}

			global.InitialiseTrace(name)
			return RunServerWithProfiler(name, setupNameGenServer)
		},
	}
	idGenApiBaseUrl = ""
	idGenApiPort int
)

var firstNames = []string{"Alfred", "Charles", "Victor", "Jean", "Tim", "Sue", "Ada", "David", "John"}
var lastNames = []string{"Aho", "Babbage", "Bahl", "Bartik", "Barners-Lee", "Black", "Lovelace", "Blie", "Carmack"}

func init() {
	NameGenCmd.Flags().StringVar(&idGenApiBaseUrl, "idgen-api-base-url", "", "idgen api base address")
	NameGenCmd.Flags().IntVar(&idGenApiPort, "idgen-api-port", 0, "idgen api port - used only with service discovery")
}

func setupNameGenServer(r *mux.Router) error {
	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport, otelhttp.WithTracerProvider(global.TracerProvider), otelhttp.WithSpanOptions(trace.WithAttributes(semconv.PeerServiceKey.String("idgen")))),
	}
	r.PathPrefix("/names/next").
		Methods("GET").
		HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			randSrc := rand.New(rand.NewSource(time.Now().UnixNano()))
			fni := randSrc.Intn(len(firstNames))
			lni := randSrc.Intn(len(lastNames))
			name := fmt.Sprintf("%s %s", firstNames[fni], lastNames[lni])

			idGenApiBaseUrl = strings.TrimRight(idGenApiBaseUrl, "/")
			getID, _ := http.NewRequestWithContext(req.Context(), "GET", fmt.Sprintf("%s/ids/next", idGenApiBaseUrl), nil)
			response, err := client.Do(getID)
			if err != nil {
				log.Print(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			var id uint
			err = json.NewDecoder(response.Body).Decode(&id)
			if err != nil {
				log.Print(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			encoder := json.NewEncoder(res)
			err = encoder.Encode(map[string]interface{}{
				"id":   id,
				"name": name,
			})
			if err != nil {
				log.Printf("%v\n", err)
				res.WriteHeader(http.StatusInternalServerError)
			}
		})

	return nil
}
