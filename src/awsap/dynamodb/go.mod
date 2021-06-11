module github.com/buddyspike/awsap/dynamodb

go 1.15

require (
	github.com/aws/aws-sdk-go-v2 v1.4.0
	github.com/aws/aws-sdk-go-v2/config v1.1.7
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.2.3
	github.com/aws/smithy-go v1.4.0
	github.com/pkg/errors v0.9.1
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.20.0
	go.opentelemetry.io/otel/sdk v0.20.0
	go.opentelemetry.io/otel/trace v0.20.0
)
