package doh

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

type blockingRoundTripper struct{}

func (blockingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	<-req.Context().Done()
	return nil, req.Context().Err()
}

func TestRequestTimeout(t *testing.T) {
	u, err := NewUpstream("https://example.com/dns-query", blockingRoundTripper{}, nil, 30*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	_, err = u.ExchangeContext(context.Background(), make([]byte, 12))
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected request deadline, got %v", err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("configured request timeout was not applied: %s", elapsed)
	}

	defaultUpstream, err := NewUpstream("https://example.com/dns-query", blockingRoundTripper{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if defaultUpstream.requestTimeout != defaultDoHTimeout {
		t.Fatalf("default request timeout = %s, want %s", defaultUpstream.requestTimeout, defaultDoHTimeout)
	}
}
