package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/justinas/alice"
	"github.com/rs/zerolog/log"

	"github.com/labd/terraform-github-registry/internal"
)

var version = "development"

func main() {
	handler := createHandler()
	lambda.Start(handler)
}

func createHandler() func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	internal.InitLogging()

	r := internal.NewRouter()
	handler := alice.New(
		internal.PanicHandler(),
		internal.LoggingHandler(log.Logger),
	).Then(r)

	adapter := httpadapter.NewV2(handler)

	return func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		response, err := adapter.ProxyWithContext(ctx, req)
		return response, err
	}
}
