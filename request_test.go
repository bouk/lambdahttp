package lambdahttp

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

var _ http.ResponseWriter = &responseWriter{}

var getEvent = `{
  "resource": "/{proxy+}",
  "path": "/pets/tobi",
  "httpMethod": "GET",
  "headers": {
    "Accept": "*/*",
    "CloudFront-Forwarded-Proto": "https",
    "CloudFront-Is-Desktop-Viewer": "true",
    "CloudFront-Is-Mobile-Viewer": "false",
    "CloudFront-Is-SmartTV-Viewer": "false",
    "CloudFront-Is-Tablet-Viewer": "false",
    "CloudFront-Viewer-Country": "CA",
    "Host": "apex-ping.com",
    "User-Agent": "curl/7.48.0",
    "Via": "2.0 a44b4468444ef3ee67472bd5c5016098.cloudfront.net (CloudFront)",
    "X-Amz-Cf-Id": "VRxPGF8rOXD7xpRjAjseXfRrFD3wg-QPUHY6chzB9bR7pXlct1NTpg==",
    "X-Amzn-Trace-Id": "Root=1-59554c99-4375fc8705ccb554008b3aad",
    "X-Forwarded-For": "207.102.57.26, 54.182.214.69",
    "X-Forwarded-Port": "443",
    "X-Forwarded-Proto": "https"
  },
  "queryStringParameters": {
    "format": "json"
  },
  "pathParameters": {
    "proxy": "pets/tobi"
  },
  "stageVariables": {
    "env": "prod"
  },
  "requestContext": {
    "path": "/pets/tobi",
    "accountId": "111111111",
    "resourceId": "jcl9w3",
    "stage": "prod",
    "requestId": "344b184b-5cfc-11e7-8483-27dbb2d30a77",
    "identity": {
      "cognitoIdentityPoolId": null,
      "accountId": null,
      "cognitoIdentityId": null,
      "caller": null,
      "apiKey": "",
      "sourceIp": "207.102.57.26",
      "accessKey": null,
      "cognitoAuthenticationType": null,
      "cognitoAuthenticationProvider": null,
      "userArn": null,
      "userAgent": "curl/7.48.0",
      "user": null
    },
    "resourcePath": "/{proxy+}",
    "httpMethod": "GET",
    "apiId": "iwcgwgigca"
  },
  "body": null,
  "isBase64Encoded": false
}`

var getEventBasicAuth = `{
  "resource": "/{proxy+}",
  "path": "/pets/tobi",
  "httpMethod": "GET",
  "headers": {
    "Accept": "*/*",
    "CloudFront-Forwarded-Proto": "https",
    "CloudFront-Is-Desktop-Viewer": "true",
    "CloudFront-Is-Mobile-Viewer": "false",
    "CloudFront-Is-SmartTV-Viewer": "false",
    "CloudFront-Is-Tablet-Viewer": "false",
    "CloudFront-Viewer-Country": "CA",
    "Host": "apex-ping.com",
    "User-Agent": "curl/7.48.0",
    "Via": "2.0 a44b4468444ef3ee67472bd5c5016098.cloudfront.net (CloudFront)",
    "X-Amz-Cf-Id": "VRxPGF8rOXD7xpRjAjseXfRrFD3wg-QPUHY6chzB9bR7pXlct1NTpg==",
    "X-Amzn-Trace-Id": "Root=1-59554c99-4375fc8705ccb554008b3aad",
    "X-Forwarded-For": "207.102.57.26, 54.182.214.69",
    "X-Forwarded-Port": "443",
    "X-Forwarded-Proto": "https",
		"Authorization": "Basic dG9iaTpmZXJyZXQ="
  },
  "queryStringParameters": {
    "format": "json"
  },
  "pathParameters": {
    "proxy": "pets/tobi"
  },
  "stageVariables": {
    "env": "prod"
  },
  "requestContext": {
    "path": "/pets/tobi",
    "accountId": "111111111",
    "resourceId": "jcl9w3",
    "stage": "prod",
    "requestId": "344b184b-5cfc-11e7-8483-27dbb2d30a77",
    "identity": {
      "cognitoIdentityPoolId": null,
      "accountId": null,
      "cognitoIdentityId": null,
      "caller": null,
      "apiKey": "",
      "sourceIp": "207.102.57.26",
      "accessKey": null,
      "cognitoAuthenticationType": null,
      "cognitoAuthenticationProvider": null,
      "userArn": null,
      "userAgent": "curl/7.48.0",
      "user": null
    },
    "resourcePath": "/{proxy+}",
    "httpMethod": "GET",
    "apiId": "iwcgwgigca"
  },
  "body": null,
  "isBase64Encoded": false
}`

var postEvent = `{
  "resource": "/{proxy+}",
  "path": "/pets/tobi",
  "httpMethod": "POST",
  "headers": {
    "Accept": "*/*",
    "CloudFront-Forwarded-Proto": "https",
    "CloudFront-Is-Desktop-Viewer": "true",
    "CloudFront-Is-Mobile-Viewer": "false",
    "CloudFront-Is-SmartTV-Viewer": "false",
    "CloudFront-Is-Tablet-Viewer": "false",
    "CloudFront-Viewer-Country": "CA",
    "content-type": "application/json",
    "Host": "apex-ping.com",
    "User-Agent": "curl/7.48.0",
    "Via": "2.0 b790a9f06b09414fec5d8b87e81d4b7f.cloudfront.net (CloudFront)",
    "X-Amz-Cf-Id": "_h1jFD3wjq6ZIyr8be6RS7Y7665jF9SjACmVodBMRefoQCs7KwTxMw==",
    "X-Amzn-Trace-Id": "Root=1-59554cc9-35de2f970f0fdf017f16927f",
    "X-Forwarded-For": "207.102.57.26, 54.182.214.86",
    "X-Forwarded-Port": "443",
    "X-Forwarded-Proto": "https"
  },
  "queryStringParameters": null,
  "pathParameters": {
    "proxy": "pets/tobi"
  },
  "requestContext": {
    "path": "/pets/tobi",
    "accountId": "111111111",
    "resourceId": "jcl9w3",
    "stage": "prod",
    "requestId": "50f6e0ce-5cfc-11e7-ada1-4f5cfe727f01",
    "identity": {
      "cognitoIdentityPoolId": null,
      "accountId": null,
      "cognitoIdentityId": null,
      "caller": null,
      "apiKey": "",
      "sourceIp": "207.102.57.26",
      "accessKey": null,
      "cognitoAuthenticationType": null,
      "cognitoAuthenticationProvider": null,
      "userArn": null,
      "userAgent": "curl/7.48.0",
      "user": null
    },
    "resourcePath": "/{proxy+}",
    "httpMethod": "POST",
    "apiId": "iwcgwgigca"
  },
  "body": "{ \"name\": \"Tobi\" }",
  "isBase64Encoded": false
}`

var postEventBinary = `{
  "resource": "/{proxy+}",
  "path": "/pets/tobi",
  "httpMethod": "POST",
  "headers": {
    "Accept": "*/*",
    "CloudFront-Forwarded-Proto": "https",
    "CloudFront-Is-Desktop-Viewer": "true",
    "CloudFront-Is-Mobile-Viewer": "false",
    "CloudFront-Is-SmartTV-Viewer": "false",
    "CloudFront-Is-Tablet-Viewer": "false",
    "CloudFront-Viewer-Country": "CA",
    "content-type": "text/plain",
    "Host": "apex-ping.com",
    "User-Agent": "curl/7.48.0",
    "Via": "2.0 b790a9f06b09414fec5d8b87e81d4b7f.cloudfront.net (CloudFront)",
    "X-Amz-Cf-Id": "_h1jFD3wjq6ZIyr8be6RS7Y7665jF9SjACmVodBMRefoQCs7KwTxMw==",
    "X-Amzn-Trace-Id": "Root=1-59554cc9-35de2f970f0fdf017f16927f",
    "X-Forwarded-For": "207.102.57.26, 54.182.214.86",
    "X-Forwarded-Port": "443",
    "X-Forwarded-Proto": "https"
  },
  "queryStringParameters": null,
  "pathParameters": {
    "proxy": "pets/tobi"
  },
  "requestContext": {
    "path": "/pets/tobi",
    "accountId": "111111111",
    "resourceId": "jcl9w3",
    "stage": "prod",
    "requestId": "50f6e0ce-5cfc-11e7-ada1-4f5cfe727f01",
    "identity": {
      "cognitoIdentityPoolId": null,
      "accountId": null,
      "cognitoIdentityId": null,
      "caller": null,
      "apiKey": "",
      "sourceIp": "207.102.57.26",
      "accessKey": null,
      "cognitoAuthenticationType": null,
      "cognitoAuthenticationProvider": null,
      "userArn": null,
      "userAgent": "curl/7.48.0",
      "user": null
    },
    "resourcePath": "/{proxy+}",
    "httpMethod": "POST",
    "apiId": "iwcgwgigca"
  },
  "body": "SGVsbG8gV29ybGQ=",
  "isBase64Encoded": true
}`

func TestNewRequest(t *testing.T) {
	t.Run("GET", func(t *testing.T) {
		var in events.APIGatewayProxyRequest
		err := json.Unmarshal([]byte(getEvent), &in)
		assert.NoError(t, err, "unmarshal")

		req, err := NewRequest(context.Background(), &in)
		assert.NoError(t, err, "new request")

		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, "apex-ping.com", req.Host)
		assert.Equal(t, "/pets/tobi", req.URL.Path)
		assert.Equal(t, "format=json", req.URL.Query().Encode())
		assert.Equal(t, "207.102.57.26", req.RemoteAddr)
	})

	t.Run("POST", func(t *testing.T) {
		var in events.APIGatewayProxyRequest
		err := json.Unmarshal([]byte(postEvent), &in)
		assert.NoError(t, err, "unmarshal")

		req, err := NewRequest(context.Background(), &in)
		assert.NoError(t, err, "new request")

		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "apex-ping.com", req.Host)
		assert.Equal(t, "/pets/tobi", req.URL.Path)
		assert.Equal(t, "", req.URL.Query().Encode())
		assert.Equal(t, "207.102.57.26", req.RemoteAddr)

		b, err := ioutil.ReadAll(req.Body)
		assert.NoError(t, err, "read body")

		assert.Equal(t, `{ "name": "Tobi" }`, string(b))
	})

	t.Run("POST binary", func(t *testing.T) {
		var in events.APIGatewayProxyRequest
		err := json.Unmarshal([]byte(postEventBinary), &in)
		assert.NoError(t, err, "unmarshal")

		req, err := NewRequest(context.Background(), &in)
		assert.NoError(t, err, "new request")

		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, "/pets/tobi", req.URL.Path)
		assert.Equal(t, "", req.URL.Query().Encode())
		assert.Equal(t, "207.102.57.26", req.RemoteAddr)

		b, err := ioutil.ReadAll(req.Body)
		assert.NoError(t, err, "read body")

		assert.Equal(t, `Hello World`, string(b))
	})

	t.Run("Basic Auth", func(t *testing.T) {
		var in events.APIGatewayProxyRequest
		err := json.Unmarshal([]byte(getEventBasicAuth), &in)
		assert.NoError(t, err, "unmarshal")

		req, err := NewRequest(context.Background(), &in)
		assert.NoError(t, err, "new request")

		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, "/pets/tobi", req.URL.Path)
		user, pass, ok := req.BasicAuth()
		assert.Equal(t, "tobi", user)
		assert.Equal(t, "ferret", pass)
		assert.True(t, ok)
	})
}
