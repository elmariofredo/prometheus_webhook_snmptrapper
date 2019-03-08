package main

import (
	flag "flag"
	sync "sync"

	config "github.com/sysincz/prometheus_webhook_snmptrapper/config"
	snmptrapper "github.com/sysincz/prometheus_webhook_snmptrapper/snmptrapper"
	types "github.com/sysincz/prometheus_webhook_snmptrapper/types"
	webhook "github.com/sysincz/prometheus_webhook_snmptrapper/webhook"

	logrus "github.com/Sirupsen/logrus"
)

var (
	//conf       config.Config
	log        = logrus.WithFields(logrus.Fields{"logger": "main"})
	waitGroup  = &sync.WaitGroup{}
	configFile = flag.String("config", "/config/snmptrapper.yaml", "The Snmptrapper configuration file")
	debug      = flag.Bool("debug", false, "Set Log to debug level and print as text")
)

func init() {

	// Process the command-line parameters:
	// flag.StringVar(&conf.SNMPTrapAddress, "snmptrapaddress", "127.0.0.1:162", "Address to send SNMP traps to")
	// flag.StringVar(&conf.SNMPCommunity, "snmpcommunity", "public", "SNMP community string")
	// flag.UintVar(&conf.SNMPRetries, "snmpretries", 1, "Number of times to retry sending SNMP traps")
	// flag.StringVar(&conf.WebhookAddress, "webhookaddress", "0.0.0.0:9099", "Address and port to listen for webhooks on")
	//flag.StringVar(&conf.CongifFile, "configFile", "/config/snmptrapper.yaml", "Oid config file")

	flag.Parse()

	if *debug {
		// The TextFormatter is default, you don't actually have to do this.
		logrus.SetFormatter(&logrus.TextFormatter{})
		// Set the log-level:
		//logrus.SetLevel(logrus.DebugLevel)
	} else {
		// Log as JSON instead of the default ASCII formatter.
		logrus.SetFormatter(&logrus.JSONFormatter{})
		// Set the log-level:
		//logrus.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	log.Infof("Loading configuration file %s", *configFile)
	conf, _, err := config.LoadConfigFile(*configFile)
	if err != nil {
		log.Errorf("Error loading configuration: %s", err)
		//return err
	}

	for _, oid := range conf.Oids {
		log.WithFields(logrus.Fields{"OidName": oid.OidName, "OidNumber": oid.OidNumber, "Template": oid.Template}).Debug("Oids from config")
	}
	// Make sure we wait for everything to complete before bailing out:
	defer waitGroup.Wait()

	// Prepare a channel of events (to feed the digester):
	log.Info("Preparing the alerts channel")
	alertsChannel := make(chan types.Alert)

	// Prepare to have background GoRoutines running:
	waitGroup.Add(1)

	// Start webhook server:
	go webhook.Run(*conf, alertsChannel, waitGroup)

	// Start the SNMP trapper:
	go snmptrapper.Run(*conf, alertsChannel, waitGroup)

}
