package dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
)

type MonotonicIDGenerator struct {
	tableName string
	config    *aws.Config
	client    *dynamodb.Client
	tracer    trace.Tracer
	tp        trace.TracerProvider
}

func (g *MonotonicIDGenerator) Generate(ctx context.Context, scope string) (uint64, error) {
	ctx, span := g.tracer.Start(ctx, "MonotonicIDGenerator.Generate")
	defer span.End()
	err := g.ensureScopeExists(ctx, scope)
	if err != nil {
		return 0, err
	}
	return g.atomicIncrement(ctx, scope)
}

func (g *MonotonicIDGenerator) ensureScopeExists(ctx context.Context, scope string) error {
	ctx, put := g.tracer.Start(ctx, "MonotoinicIDGenerator.ensureScopeExists")
	defer put.End()
	_, err := g.client.PutItem(
		ctx,
		&dynamodb.PutItemInput{
			TableName: aws.String(g.tableName),
			Item: map[string]types.AttributeValue{
				"scope_id": &types.AttributeValueMemberS{Value: scope},
				"value":    &types.AttributeValueMemberN{Value: "0"},
			},
			ConditionExpression: aws.String("attribute_not_exists(scope_id)"),
			ReturnValues:        types.ReturnValueNone,
		})

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() != "ConditionalCheckFailedException" {
				return errors.WithStack(err)
			}
		}
	} else {
		put.AddEvent(fmt.Sprintf("new scope item created: %v", scope))
	}

	return nil
}

func (g *MonotonicIDGenerator) atomicIncrement(ctx context.Context, scope string) (uint64, error) {
	ctx, increment := g.tracer.Start(ctx, "MonotoinicIDGenerator.atomicIncrement")
	defer increment.End()

	out, err := g.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(g.tableName),
		Key: map[string]types.AttributeValue{
			"scope_id": &types.AttributeValueMemberS{Value: scope},
		},
		UpdateExpression: aws.String("SET #v = #v + :incr"),
		ExpressionAttributeNames: map[string]string{
			"#v": "value",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":incr": &types.AttributeValueMemberN{Value: "1"},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	})

	if err != nil {
		return 0, errors.WithStack(err)
	}

	v, err := strconv.ParseUint(out.Attributes["value"].(*types.AttributeValueMemberN).Value, 10, 64)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return v, nil
}

type Option func(generator *MonotonicIDGenerator)

func WithConfig(config *aws.Config) Option {
	return func(generator *MonotonicIDGenerator) {
		generator.config = config
	}
}

func WithTracer(tracer trace.Tracer) Option {
	return func(generator *MonotonicIDGenerator) {
		generator.tracer = tracer
	}
}

func WithTraceProvider(tp trace.TracerProvider) Option {
	return func(generator *MonotonicIDGenerator) {
		generator.tp = tp
	}
}

func NewMonotonicIDGenerator(tableName string, options ...Option) (*MonotonicIDGenerator, error) {
	g := &MonotonicIDGenerator{tableName: tableName}
	for _, o := range options {
		o(g)
	}
	if g.tp == nil {
		g.tp = trace.NewNoopTracerProvider()
	}
	g.tracer = g.tp.Tracer("awsap/dynamodb/monotonic_id_generator")
	if g.config == nil {
		config, err := config.LoadDefaultConfig(context.Background(), config.WithHTTPClient(&http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport, otelhttp.WithTracerProvider(g.tp), otelhttp.WithSpanOptions(trace.WithAttributes(semconv.PeerServiceKey.String("dynamodb")))),
		}))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		g.config = &config
	}

	g.client = dynamodb.NewFromConfig(*g.config)
	return g, nil
}
