package webhook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	types "github.com/sysincz/prometheus_webhook_snmptrapper/types"

	logrus "github.com/Sirupsen/logrus"
	template "github.com/prometheus/alertmanager/template"
)

//Handler A webhook handler with a "ServeHTTP" method:
type Handler struct {
	AlertsChannel chan types.Alert
}

// Handle webhook requests:
func (webhookHandler *Handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {

	// Read the request body:
	payload, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Failed to read the request body")
		http.Error(responseWriter, "Failed to read the request body", http.StatusBadRequest)
		return
	}

	// Validate the payload:
	alerts, err := validatePayload(payload)
	if err != nil {
		http.Error(responseWriter, "Failed to unmarshal the request-body into an alert", http.StatusBadRequest)
		return
	}

	// Send the alerts to the snmp-trapper:
	for alertIndex, alert := range alerts {
		log.WithFields(logrus.Fields{"index": alertIndex, "status": alert.Status, "labels": alert.Labels}).Debug("Forwarding an alert to the SNMP trapper")

		// Enrich the request with the remote-address:
		alert.Address = request.RemoteAddr

		// Put the alert onto the alerts-channel:
		webhookHandler.AlertsChannel <- alert
	}

}

//validatePayload Validate a webhook payload and return a list of Alerts:
func validatePayload(payload []byte) ([]types.Alert, error) {

	// Make our response:
	alerts := make([]types.Alert, 0)

	// Make a new Prometheus data-structure to unmarshal the request body into:
	prometheusData := &template.Data{}

	// Unmarshal the request body into the alert:
	err := json.Unmarshal(payload, prometheusData)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err, "payload": payload}).Error("Failed to unmarshal the request body into an alert")
		return alerts, err
	}

	log.WithFields(logrus.Fields{"payload": string(payload)}).Debug("Received a valid webhook alert")

	// Iterate over the list of alerts:
	for _, alertDetails := range prometheusData.Alerts {

		// Make a new SNMP alert:
		alerts = append(alerts, types.Alert{
			Status:            prometheusData.Status,
			Labels:            alertDetails.Labels,
			Annotations:       alertDetails.Annotations,
			StartsAt:          alertDetails.StartsAt,
			EndsAt:            alertDetails.EndsAt,
			Receiver:          prometheusData.Receiver,
			GroupLabels:       prometheusData.GroupLabels,
			CommonLabels:      prometheusData.CommonLabels,
			CommonAnnotations: prometheusData.CommonAnnotations,
			ExternalURL:       prometheusData.ExternalURL,
		})

	}

	return alerts, nil
}
