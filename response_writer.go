package lambdahttp

import (
	"bytes"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

func newResponseWriter() *responseWriter {
	return &responseWriter{
		status: 200,
	}
}

type responseWriter struct {
	buf    bytes.Buffer
	header http.Header
	status int
}

func (rw *responseWriter) Header() http.Header {
	if rw.header == nil {
		rw.header = make(http.Header)
	}

	return rw.header
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.buf.Write(b)
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
}

func (rw *responseWriter) generateResponse() events.APIGatewayProxyResponse {
	response := events.APIGatewayProxyResponse{
		Body:       rw.buf.String(),
		StatusCode: rw.status,
		Headers:    make(map[string]string, len(rw.header)),
	}

	var hasContent, hasDate bool

	for k, v := range rw.header {
		if len(v) > 0 {
			response.Headers[k] = v[0]
			if k == "Date" {
				hasDate = true
			} else if k == "Content-Type" {
				hasContent = true
			}
		}
	}
	if !hasContent && len(response.Body) > 0 {
		response.Headers["Content-Type"] = http.DetectContentType(rw.buf.Bytes())
	}
	if !hasDate {
		response.Headers["Date"] = time.Now().UTC().Format(http.TimeFormat)
	}
	return response
}
