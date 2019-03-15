package snmptrapper

import (
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	config "github.com/sysincz/prometheus_webhook_snmptrapper/config"
	template "github.com/sysincz/prometheus_webhook_snmptrapper/template"
	types "github.com/sysincz/prometheus_webhook_snmptrapper/types"

	logrus "github.com/Sirupsen/logrus"
	snmpgo "github.com/k-sone/snmpgo"
)

var (
	alertsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "snmptrapper_processed_alerts_total",
		Help: "The total number of processed events",
	})

	alertsForwarded = promauto.NewCounter(prometheus.CounterOpts{
		Name: "snmptrapper_forwarded_alerts_total",
		Help: "The total number of processed events",
	})

	alertsFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "snmptrapper_failed_alerts_total",
		Help: "The total number of failed events",
	})
)

func toOid(oidNumber string) *snmpgo.Oid {
	oid, _ := snmpgo.NewOid(oidNumber)
	return oid
}

func sendTrap(alert types.Alert) {

	alertsProcessed.Inc()
	// Prepare an SNMP handler:
	snmp, err := snmpgo.NewSNMP(snmpgo.SNMPArguments{
		Version:   snmpgo.V2c,
		Address:   myConfig.SNMPTrapAddress,
		Retries:   myConfig.SNMPRetries,
		Community: myConfig.SNMPCommunity,
	})
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Failed to create snmpgo.SNMP object")
		alertsFailed.Inc()
		return
	}
	log.WithFields(logrus.Fields{"address": myConfig.SNMPTrapAddress, "retries": myConfig.SNMPRetries, "community": myConfig.SNMPCommunity}).Debug("Created snmpgo.SNMP object")
	RunTemplate("{{ printf \"%#v\" . }}", alert)
	// Build VarBind list:
	var varBinds snmpgo.VarBinds
	varBinds = append(varBinds, snmpgo.NewVarBind(snmpgo.OidSysUpTime, snmpgo.NewTimeTicks(1000)))
	// The "enterprise OID" for the trap (rising/firing or falling/recovery):
	if alert.Status == "firing" {
		varBinds = append(varBinds, snmpgo.NewVarBind(snmpgo.OidSnmpTrap, toOid(myConfig.FiringTrap)))
	} else {
		varBinds = append(varBinds, snmpgo.NewVarBind(snmpgo.OidSnmpTrap, toOid(myConfig.RecoveryTrap)))
	}

	// Insert the AlertManager variables:
	for _, oid := range myConfig.Oids {
		ret := RunTemplate(oid.Template, alert)
		if !notEmpty(oid, ret) {
			alertsFailed.Inc()
			return
		}

		if oid.Type == "int32" {
			varBinds = append(varBinds, snmpgo.NewVarBind(toOid(oid.OidNumber), snmpgo.NewInteger(strToInt32(ret))))
		} else {
			varBinds = append(varBinds, snmpgo.NewVarBind(toOid(oid.OidNumber), snmpgo.NewOctetString([]byte(ret))))
		}

	}

	//fmt.Printf("%+v\n", varBinds)
	// Create an SNMP "connection":
	if err = snmp.Open(); err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Failed to open SNMP connection")
		alertsFailed.Inc()
		return
	}
	defer snmp.Close()

	// Send the trap:
	if err = snmp.V2Trap(varBinds); err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Failed to send SNMP trap")
		alertsFailed.Inc()
		return
	}
	log.WithFields(logrus.Fields{"status": alert.Status}).Info("It's a trap!")
	alertsForwarded.Inc()

}
func notEmpty(oid *config.OidConfig, text string) bool {
	if oid.NotEmpty {
		if text == "" {
			log.WithFields(logrus.Fields{"error": "Value is empty", "oid": oid.OidName, "Template": oid.Template}).Error("Failed to create snmpgo.SNMP object")
			return false
		}
	}
	return true
}
func strToInt32(text string) int32 {
	//convert string to int32
	i1, err := strconv.Atoi(text)
	if err == nil {

		return int32(i1)
	}
	log.Errorf("Unabele to convert string '%s' to int32. %s", text, err)
	return int32(i1)
}

//RunTemplate translate template string to string + trimSpace
func RunTemplate(text string, data interface{}) string {
	tmpl := template.Init()

	value, err := tmpl.Execute(text, data)
	if err != nil {
		log.Errorf("Error loading templates from %s: %s", text, err)
		return ""
	}
	value = strings.TrimSpace(value)
	return value
}
