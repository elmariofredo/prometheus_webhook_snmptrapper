package snmptrapper

import (
	"os"
	"os/signal"
	"sync"

	config "github.com/sysincz/prometheus_webhook_snmptrapper/config"
	types "github.com/sysincz/prometheus_webhook_snmptrapper/types"

	logrus "github.com/Sirupsen/logrus"
)

var (
	log      = logrus.WithFields(logrus.Fields{"logger": "SNMP-trapper"})
	myConfig config.Config
)

func init() {
	// Set the log-level:
	logrus.SetLevel(logrus.DebugLevel)

}

//Run main function for send traps call sendTrap
func Run(myConfigFromMain config.Config, alertsChannel chan types.Alert, waitGroup *sync.WaitGroup) {

	log.WithFields(logrus.Fields{"address": myConfigFromMain.SNMPTrapAddress}).Info("Starting the SNMP trapper")

	// Populate the config:
	myConfig = myConfigFromMain

	// Set up a channel to handle shutdown:
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill, os.Interrupt)

	// Handle incoming alerts:
	go func() {
		for {
			select {

			case alert := <-alertsChannel:

				// Send a trap based on this alert:
				log.WithFields(logrus.Fields{"status": alert.Status}).Debug("Received an alert")
				sendTrap(alert)
			}
		}
	}()

	// Wait for shutdown:
	for {
		select {
		case <-signals:
			log.Warn("Shutting down the SNMP trapper")

			// Tell main() that we're done:
			waitGroup.Done()
			return
		}
	}

}
