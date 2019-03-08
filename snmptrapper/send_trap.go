package snmptrapper

import (
	"bytes"
	"text/template"
	"time"

	types "github.com/sysincz/prometheus_webhook_snmptrapper/types"

	logrus "github.com/Sirupsen/logrus"
	snmpgo "github.com/k-sone/snmpgo"
)

func toOid(oidNumber string) *snmpgo.Oid {
	oid, _ := snmpgo.NewOid(oidNumber)
	return oid
}

func sendTrap(alert types.Alert) {

	// Prepare an SNMP handler:
	snmp, err := snmpgo.NewSNMP(snmpgo.SNMPArguments{
		Version:   snmpgo.V2c,
		Address:   myConfig.SNMPTrapAddress,
		Retries:   myConfig.SNMPRetries,
		Community: myConfig.SNMPCommunity,
	})
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Failed to create snmpgo.SNMP object")
		return
	} else {
		log.WithFields(logrus.Fields{"address": myConfig.SNMPTrapAddress, "retries": myConfig.SNMPRetries, "community": myConfig.SNMPCommunity}).Debug("Created snmpgo.SNMP object")
	}

	// Build VarBind list:
	var varBinds snmpgo.VarBinds

	// The "enterprise OID" for the trap (rising/firing or falling/recovery):
	if alert.Status == "firing" {
		varBinds = append(varBinds, snmpgo.NewVarBind(snmpgo.OidSnmpTrap, toOid(myConfig.FiringTrap)))
		//	varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.TimeStamp, snmpgo.NewOctetString([]byte(alert.StartsAt.Format(time.RFC3339)))))
		varBinds = append(varBinds, snmpgo.NewVarBind(toOid("1.3.6.1.3.1977.1.1.7"), snmpgo.NewOctetString([]byte(alert.StartsAt.Format(time.RFC3339)))))
	} else {
		varBinds = append(varBinds, snmpgo.NewVarBind(snmpgo.OidSnmpTrap, toOid(myConfig.RecoveryTrap)))
		//	varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.TimeStamp, snmpgo.NewOctetString([]byte(alert.EndsAt.Format(time.RFC3339)))))

		varBinds = append(varBinds, snmpgo.NewVarBind(toOid("1.3.6.1.3.1977.1.1.7"), snmpgo.NewOctetString([]byte(alert.EndsAt.Format(time.RFC3339)))))
	}
	runTemplate("DataDump", "{{ printf \"%#v\" . }}", alert)
	// Insert the AlertManager variables:
	for _, oid := range myConfig.Oids {
		ret := runTemplate(oid.OidName, oid.Template, alert)
		varBinds = append(varBinds, snmpgo.NewVarBind(toOid(oid.OidNumber), snmpgo.NewOctetString([]byte(ret))))
	}

	//fmt.Printf("%+v\n", varBinds)
	// Create an SNMP "connection":
	if err = snmp.Open(); err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Failed to open SNMP connection")
		return
	}
	defer snmp.Close()

	// Send the trap:
	if err = snmp.V2Trap(varBinds); err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Failed to send SNMP trap")
		return
	} else {
		log.WithFields(logrus.Fields{"status": alert.Status}).Info("It's a trap!")
	}
}

func runTemplate(name string, templateDef string, data interface{}) string {
	tmpl, err := template.New(name).Parse(templateDef)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	ret := buf.String()
	log.WithFields(logrus.Fields{"Name": name, "Template": templateDef, "Output": ret}).Debug("Template processing")
	if err != nil {
		log.Error(err)
	}
	return ret
}
