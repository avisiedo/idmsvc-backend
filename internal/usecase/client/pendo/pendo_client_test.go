package pendo

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"github.com/podengo-project/idmsvc-backend/internal/interface/client/pendo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8031/pendo/v1"

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func helperPendoConfig() *config.Config {
	return &config.Config{
		Clients: config.Clients{
			PendoBaseURL:            baseURL,
			PendoAPIKey:             "test-api-key",
			PendoRequestTimeoutSecs: 1,
		},
	}
}

func helperNewPendo(cfg *config.Config, fn RoundTripFunc) pendo.Pendo {
	client := newClient(cfg)
	client.Client.Transport = RoundTripFunc(fn)
	return client
}

func TestNewPendo(t *testing.T) {
	assert.PanicsWithValue(t, "'cfg' is nil", func() {
		NewClient(nil)
	})

	assert.PanicsWithValue(t, "'PendoBaseURL' is empty", func() {
		NewClient(&config.Config{
			Clients: config.Clients{
				PendoBaseURL:            "",
				PendoAPIKey:             "",
				PendoRequestTimeoutSecs: 0,
			},
		})
	})

	assert.PanicsWithValue(t, "'PendoAPIKey' is empty", func() {
		NewClient(&config.Config{
			Clients: config.Clients{
				PendoBaseURL:            baseURL,
				PendoAPIKey:             "",
				PendoRequestTimeoutSecs: 0,
			},
		})
	})

	client := NewClient(&config.Config{
		Clients: config.Clients{
			PendoBaseURL:            baseURL,
			PendoAPIKey:             "kiudsahfq84radihfa",
			PendoRequestTimeoutSecs: 0,
		},
	})
	require.NotNil(t, client)
}

func TestGuardObserver(t *testing.T) {
	cfg := helperPendoConfig()
	client := newClient(cfg)

	require.EqualError(t, client.guardObserve("", "", nil), "'kind' is an empty string")
	require.EqualError(t, client.guardObserve("mykind", "", nil), "'group' is an empty string")
	require.EqualError(t, client.guardObserve("mykind", "mygroup", nil), "'metrics' is nil")
	pendoMetrics := make([]pendo.PendoMetric, 0, 1)
	require.NoError(t, client.guardObserve("mykind", "mygroup", pendoMetrics), "an empty slice does not report error")
	pendoMetrics = append(pendoMetrics, pendo.PendoMetric{})
	require.EqualError(t, client.guardObserve("mykind", "mygroup", pendoMetrics), "'metrics[0].VisitorID' is an empty string")
	pendoMetrics[0].VisitorID = "my-test-visitor-id"

	// Success case
	require.NoError(t, client.guardObserve("mykind", "mygroup", pendoMetrics))
}

func TestObserver(t *testing.T) {
	// https://hassansin.github.io/Unit-Testing-http-client-in-Go#2-by-replacing-httptransport
	cfg := helperPendoConfig()
	kind := "mykind"
	group := "mygroup"
	client := helperNewPendo(cfg, func(req *http.Request) *http.Response {
		assert.Equal(t, baseURL+"/metadata/"+kind+"/"+group+"/value", req.URL.String())
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
			Header:     make(http.Header),
		}
	})

	// Panic when context is nil
	assert.PanicsWithValue(t, "'ctx' is nil", func() {
		client.Observe(nil, "", "", nil)
	})

	// Error when wrong argument
	ctx := app_context.CtxWithLog(context.TODO(), slog.Default())
	require.EqualError(t, client.Observe(ctx, "", "", nil), "'kind' is an empty string")

	// Success
	metrics := make([]pendo.PendoMetric, 0, 1)
	metrics = append(metrics, pendo.PendoMetric{
		VisitorID: "test-visitor-id",
	})
	err := client.Observe(ctx, kind, group, metrics)
	require.NoError(t, err)
}
