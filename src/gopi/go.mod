module gopi

go 1.15

require (
	github.com/buddyspike/awsap/dynamodb v0.0.0-00010101000000-000000000000
	github.com/go-playground/validator/v10 v10.6.1 // indirect
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/ugorji/go v1.2.5 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.20.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.20.0 // indirect
	go.opentelemetry.io/contrib/propagators/aws v0.20.0
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/exporters/otlp v0.20.0
	go.opentelemetry.io/otel/metric v0.20.0
	go.opentelemetry.io/otel/sdk v0.20.0
	go.opentelemetry.io/otel/sdk/metric v0.20.0
	go.opentelemetry.io/otel/trace v0.20.0
	go.uber.org/atomic v1.8.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/crypto v0.0.0-20210506145944-38f3c27a63bf // indirect
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	google.golang.org/genproto v0.0.0-20210510173355-fb37daa5cd7a // indirect
	google.golang.org/grpc v1.37.1 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace github.com/buddyspike/awsap/dynamodb => ./awsap/dynamodb
