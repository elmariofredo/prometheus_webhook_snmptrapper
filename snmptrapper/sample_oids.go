package snmptrapper

// SNMPv2 "generic" OIDs (including "RMON"): http://www.oidview.com/mibs/0/md-0-1.html

const (
	oidSNMPv2SysDescr                           = "1.3.6.1.2.1.1.1"         // Variable: summary
	oidSNMPv2SysName                            = "1.3.6.1.2.1.1.5"         // Variable: instance
	oidSNMPv2SysLocation                        = "1.3.6.1.2.1.1.6"         // Variable: location
	oidSNMPv2SysServices                        = "1.3.6.1.2.1.1.7"         // Variable: service
	oidSNMPv2SysORLastChange                    = "1.3.6.1.2.1.1.8"         //
	oidRMONAlarmStatus                          = "1.3.6.1.2.1.16.3.1.1.12" //
	oidRMONHostAddress                          = "1.3.6.1.2.1.16.4.2.1.1"  // Variable: instance
	oidRMONEventDescription                     = "1.3.6.1.2.1.16.9.1.1.2"  // Variable: summary
	oidRMONEventType                            = "1.3.6.1.2.1.16.9.1.1.3"  // Variable: severity
	oidSNMPv2LinkDown                           = "1.3.6.1.2.1.11.2"        // Trap: firing
	oidSNMPv2LinkUp                             = "1.3.6.1.2.1.11.3"        // Trap: recovery
	oidIFLinkDown                               = "1.3.6.1.6.3.1.1.5.3"     // Trap: firing
	oidIFLinkUp                                 = "1.3.6.1.6.3.1.1.5.4"     // Trap: recovery
	oidRMONRisingAlarm                          = "1.3.6.1.2.1.16.0.1"      // Notification: firing
	oidRMONFallingAlarm                         = "1.3.6.1.2.1.16.0.2"      // Notification: recovery
	oidPrometheusTrapperFiringNotification      = "1.3.6.1.3.1977.1.0.1"    // Notification: firing
	oidPrometheusTrapperRecoveryNotification    = "1.3.6.1.3.1977.1.0.2"    // Notification: recovery
	oidPrometheusTrapperNotificationInstance    = "1.3.6.1.3.1977.1.1.1"    // Variable: instance
	oidPrometheusTrapperNotificationService     = "1.3.6.1.3.1977.1.1.2"    // Variable: service
	oidPrometheusTrapperNotificationLocation    = "1.3.6.1.3.1977.1.1.3"    // Variable: location
	oidPrometheusTrapperNotificationSeverity    = "1.3.6.1.3.1977.1.1.4"    // Variable: severity
	oidPrometheusTrapperNotificationDescription = "1.3.6.1.3.1977.1.1.5"    // Variable: description
	oidPrometheusTrapperNotificationJob         = "1.3.6.1.3.1977.1.1.6"    // Variable: job
	oidPrometheusTrapperNotificationTimestamp   = "1.3.6.1.3.1977.1.1.7"    // Variable: timestamp

	oidnSvcEventIndex    = ".1.3.6.1.4.1.20006.1.3.1.1"
	oidnSvcHostname      = ".1.3.6.1.4.1.20006.1.3.1.2"
	oidnSvcHostAlias     = ".1.3.6.1.4.1.20006.1.3.1.3"
	oidnSvcHostStateID   = ".1.3.6.1.4.1.20006.1.3.1.4"
	oidnSvcHostStateType = ".1.3.6.1.4.1.20006.1.3.1.5"
	oidnSvcDesc          = ".1.3.6.1.4.1.20006.1.3.1.6"
	oidnSvcStateID       = ".1.3.6.1.4.1.20006.1.3.1.7"
	oidnSvcAttempt       = ".1.3.6.1.4.1.20006.1.3.1.8"
	oidnSvcDurationSec   = ".1.3.6.1.4.1.20006.1.3.1.9"
	oidnSvcGroupName     = ".1.3.6.1.4.1.20006.1.3.1.10"
	oidnSvcLastCheck     = ".1.3.6.1.4.1.20006.1.3.1.11"
	oidnSvcLastChange    = ".1.3.6.1.4.1.20006.1.3.1.12"
	oidnSvcLastOK        = ".1.3.6.1.4.1.20006.1.3.1.13"
	oidnSvcLastWarn      = ".1.3.6.1.4.1.20006.1.3.1.14"
	oidnSvcLastCrit      = ".1.3.6.1.4.1.20006.1.3.1.15"
	oidnSvcLastUnkn      = ".1.3.6.1.4.1.20006.1.3.1.16"
	oidnSvcOutput        = ".1.3.6.1.4.1.20006.1.3.1.17"
	oidnSvcPerfData      = ".1.3.6.1.4.1.20006.1.3.1.18"
	oidnSvcNote          = ".1.3.6.1.4.1.20006.1.3.1.19"
	oidnSvcGrapher       = ".1.3.6.1.4.1.20006.1.3.1.20"
	oidnCIIMPACT         = ".1.3.6.1.4.1.20006.1.3.1.21"
	oidnHostClass        = ".1.3.6.1.4.1.20006.1.3.1.22"
	oidnHostCountry      = ".1.3.6.1.4.1.20006.1.3.1.23"
	oidnSvcLogConfirmURL = ".1.3.6.1.4.1.20006.1.3.1.25"
	oidnSvcSource        = ".1.3.6.1.4.1.20006.1.3.1.26"
	oidnAutoTicket       = ".1.3.6.1.4.1.20006.1.3.1.27"
	oidnOTAssignee       = ".1.3.6.1.4.1.20006.1.3.1.28"
)
