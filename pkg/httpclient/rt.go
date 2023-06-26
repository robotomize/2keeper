package httpclient

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"
)

var (
	defaultMaxWait    = 1 * time.Second
	defaultMinWait    = 100 * time.Millisecond
	defaultMaxRetries = 3
)

func ExponentialBackoff(mn, mx time.Duration, n int) time.Duration {
	mult := math.Pow(2, float64(n)) * float64(mn)
	sleep := time.Duration(mult)
	if float64(sleep) != mult || sleep > mx {
		sleep = mx
	}

	return sleep
}

type Option func(*retry)

func WithMaxRetry(n int) Option {
	return func(r *retry) {
		r.maxRetries = n
	}
}

func WithWait(mn, mx time.Duration) Option {
	return func(r *retry) {
		r.maxWait = mx
		r.minWait = mn
	}
}

func NewRetry(rt http.RoundTripper, opts ...Option) http.RoundTripper {
	ret := retry{rt: rt, maxRetries: defaultMaxRetries, maxWait: defaultMaxWait, minWait: defaultMinWait}
	for _, o := range opts {
		o(&ret)
	}

	return &ret
}

type retry struct {
	maxRetries int
	maxWait    time.Duration
	minWait    time.Duration
	rt         http.RoundTripper
}

func (t *retry) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	for i := 0; i <= t.maxRetries; i++ {
		ctx := req.Context()

		resp, err = t.rt.RoundTrip(req)
		if err == nil {
			return resp, nil
		}

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			if s, ok := resp.Header["Retry-After"]; ok {
				if sleep, parseInt := strconv.ParseInt(s[0], 10, 64); parseInt == nil {
					if err = wait(ctx, time.Second*time.Duration(sleep)); err != nil {
						return nil, err
					}
					continue
				}
			}

			return nil, err
		}

		waitTime := ExponentialBackoff(t.minWait, t.maxWait, i)
		if err = wait(ctx, waitTime); err != nil {
			return nil, err
		}
	}

	return nil, err
}

func wait(ctx context.Context, t time.Duration) error {
	ticker := time.NewTicker(t)
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ticker.C:
	}

	return nil
}
