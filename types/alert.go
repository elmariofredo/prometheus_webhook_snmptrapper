package types

import (
	"time"
)

//Alert combination of prometheus alert and webhook (webhook with one alert)
type Alert struct {
	Address      string
	Status       string
	Annotations  map[string]string
	Labels       map[string]string
	StartsAt     time.Time
	EndsAt       time.Time
	GeneratorURL string

	//Prometheus webhook data
	Receiver          string
	GroupLabels       map[string]string
	CommonLabels      map[string]string
	CommonAnnotations map[string]string
	ExternalURL       string
}
