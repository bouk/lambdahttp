package lambdahttp

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
)

func encodeQuery(v map[string]string) string {
	if len(v) == 0 {
		return ""
	}

	var buf bytes.Buffer
	for k, value := range v {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(url.QueryEscape(k))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(value))
	}
	return buf.String()
}

func NewRequest(ctx context.Context, request *events.APIGatewayProxyRequest) (*http.Request, error) {
	host, _ := request.Headers["Host"]
	scheme, _ := request.Headers["X-Forwarded-Proto"]

	u := &url.URL{
		Host:     host,
		Path:     request.Path,
		RawPath:  url.PathEscape(request.Path),
		RawQuery: encodeQuery(request.QueryStringParameters),
		Scheme:   scheme,
	}

	req := &http.Request{
		Body:       http.NoBody,
		GetBody:    func() (io.ReadCloser, error) { return http.NoBody, nil },
		Header:     make(http.Header, len(request.Headers)),
		Host:       host,
		Method:     request.HTTPMethod,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		RemoteAddr: request.RequestContext.Identity.SourceIP,
		URL:        u,
	}

	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	body := request.Body
	if request.IsBase64Encoded {
		b, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return nil, errors.Wrap(err, "decoding base64 body")
		}
		body = string(b)
	}

	if len(body) > 0 {
		reader := strings.NewReader(body)
		bodyReader := *reader
		req.Body = ioutil.NopCloser(&bodyReader)
		req.ContentLength = int64(len(body))
		req.GetBody = func() (io.ReadCloser, error) {
			r := *reader
			return ioutil.NopCloser(&r), nil
		}
	}

	return req.WithContext(ctx), nil
}
