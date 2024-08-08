package pendo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/podengo-project/idmsvc-backend/internal/api/header"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"github.com/podengo-project/idmsvc-backend/internal/interface/client/pendo"
)

type pendoClient struct {
	Config *config.Config
	Client *http.Client
}

func NewClient(cfg *config.Config) pendo.Pendo {
	return newClient(cfg)
}

func newClient(cfg *config.Config) *pendoClient {
	if cfg == nil {
		panic("'cfg' is nil")
	}
	if cfg.Clients.PendoBaseURL == "" {
		panic("'PendoBaseURL' is empty")
	}
	if cfg.Clients.PendoAPIKey == "" {
		panic("'PendoAPIKey' is empty")
	}
	if cfg.Clients.PendoRequestTimeoutSecs == 0 {
		cfg.Clients.PendoRequestTimeoutSecs = 3
	}
	client := &http.Client{
		Timeout: time.Duration(cfg.Clients.PendoRequestTimeoutSecs),
	}
	return &pendoClient{
		Config: cfg,
		Client: client,
	}
}

func (c *pendoClient) guardObserve(kind, group string, metrics []pendo.PendoMetric) error {
	if kind == "" {
		return fmt.Errorf("'kind' is an empty string")
	}
	if group == "" {
		return fmt.Errorf("'group' is an empty string")
	}
	if metrics == nil {
		return fmt.Errorf("'metrics' is nil")
	}
	for i := range metrics {
		if metrics[i].VisitorID == "" {
			return fmt.Errorf("'metrics[%d].VisitorID' is an empty string", i)
		}
	}
	// TODO check 'kind' and 'group' contains valid characters
	//      we need additional information to do that
	return nil
}

// Observe launch a request to pendo service to register some metric
func (c *pendoClient) Observe(ctx context.Context, kind, group string, metrics []pendo.PendoMetric) error {
	logger := app_context.LogFromCtx(ctx)
	if err := c.guardObserve(kind, group, metrics); err != nil {
		logger.Error("wrong arguments: " + err.Error())
		return err
	}

	// https://github.com/RedHatInsights/cloud-connector/blob/master/internal/pendo_transmitter/pendo_reporter.go#L51
	reqBody, err := json.Marshal(metrics)
	if err != nil {
		logger.Error("marshaling the pendo metric")
		return err
	}

	// Prepare the request
	url := c.Config.Clients.PendoBaseURL + "/metadata/" + url.PathEscape(kind) + "/" + url.PathEscape(group) + "/value"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		logger.Error("creating request to observe pendo metric")
		return err
	}

	// Add headers to the request
	req.Header.Set("content-type", "application/json")
	req.Header.Set(header.HeaderXPendoIntegrationKey, c.Config.Clients.PendoAPIKey)

	// Launch request
	resp, err := c.Client.Do(req)
	switch {
	case err != nil:
		logger.Error("doing request to pendo service")
		return err
	case resp.StatusCode != 200:
		logger.Error("pendo service returned unsuccessful")
		return fmt.Errorf("Pendo request unsuccessful: status: %v", resp.StatusCode)
	}

	// Maybe we don't need the body, but adding general approach
	defer resp.Body.Close()
	_, err = io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		logger.Error("reading pendo response body")
		return err
	}

	logger.Debug("pendo metric observed successfully")
	return nil
}
