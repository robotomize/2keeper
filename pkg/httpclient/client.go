package httpclient

import (
	"net/http"
)

type Do interface {
	Do(req *http.Request) (*http.Response, error)
}

type retryClient struct {
	*http.Client
}

func NewClient(rt http.RoundTripper) Do {
	return &retryClient{Client: &http.Client{Transport: rt}}
}
