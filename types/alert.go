package types

import (
	"time"
)

type Alert struct {
	Address      string
	Status       string
	Annotations  map[string]string
	Labels       map[string]string
	StartsAt     time.Time
	EndsAt       time.Time
	GeneratorURL string

	//Prometheus data
	Receiver string
	//Status   string
	//Alerts   Alerts `json:"alerts"`

	GroupLabels       map[string]string
	CommonLabels      map[string]string
	CommonAnnotations map[string]string

	ExternalURL string
}
