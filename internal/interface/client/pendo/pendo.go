package pendo

import "context"

type PendoLabels map[string]string

type PendoMetric struct {
	VisitorID string            `json:"visitorId"`
	Values    map[string]string `json:"values"`
}

type Pendo interface {
	Observe(ctx context.Context, kind, group string, metrics []PendoMetric) error
}
