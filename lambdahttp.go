package lambdahttp

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Start(handler http.Handler) {
	f := func(ctx context.Context, request events.APIGatewayProxyRequest) (response events.APIGatewayProxyResponse, err error) {
		var req *http.Request

		rw := newResponseWriter()
		req, err = NewRequest(ctx, &request)
		if err != nil {
			return
		}

		handler.ServeHTTP(rw, req)

		response = rw.generateResponse()
		return
	}

	lambda.Start(f)
}
