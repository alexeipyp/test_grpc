package httpclient

import (
	"net/http"
	"net/http/httputil"
	"time"

	"go.uber.org/zap"
)

func New(username string, password string, timeout time.Duration, hhtpTraceLogger *zap.Logger) *http.Client {
	client := &http.Client{
		Timeout: timeout,
		Transport: &CustomTransport{
			httpTraceLogger: hhtpTraceLogger,
			username:        username,
			password:        password,
		},
	}

	return client
}

type CustomTransport struct {
	httpTraceLogger *zap.Logger
	username        string
	password        string
}

func (t *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(t.username, t.password)

	reqBytes, _ := httputil.DumpRequestOut(req, true)
	t.httpTraceLogger.Info("request", zap.ByteString("contains", reqBytes))

	resp, err := http.DefaultTransport.RoundTrip(req)

	if resp != nil {
		respBytes, _ := httputil.DumpResponse(resp, true)
		t.httpTraceLogger.Info("response", zap.ByteString("contains", respBytes))
	}

	return resp, err
}
