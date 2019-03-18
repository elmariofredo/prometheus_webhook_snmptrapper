package webhook

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	config "github.com/sysincz/prometheus_webhook_snmptrapper/config"
	types "github.com/sysincz/prometheus_webhook_snmptrapper/types"

	logrus "github.com/Sirupsen/logrus"
)

var (
	log      = logrus.WithFields(logrus.Fields{"logger": "Webhook-server"})
	myConfig config.Config
)

func init() {
	// Set the log-level:
	logrus.SetLevel(logrus.DebugLevel)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok!")
}

//Run main function for webhook
func Run(myConfigFromMain config.Config, alertsChannel chan types.Alert, waitGroup *sync.WaitGroup) {

	log.WithFields(logrus.Fields{"address": myConfigFromMain.WebhookAddress}).Info("Starting the Webhook server")

	// Populate the config:
	myConfig = myConfigFromMain

	// Set up a channel to handle shutdown:
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill, os.Interrupt)
	http.HandleFunc("/healthz", healthz)
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/alerts", &Handler{AlertsChannel: alertsChannel})
	// Listen for webhooks:
	http.ListenAndServe(myConfig.WebhookAddress, nil)

	// Wait for shutdown:
	for {
		select {
		case <-signals:
			log.Info("Shutting down the Webhook server")

			// Tell main() that we're done:
			waitGroup.Done()
			return
		}
	}

}
