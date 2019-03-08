FROM alpine:3.6

RUN apk --no-cache add ca-certificates tini curl bash net-snmp-tools


EXPOSE 9099

COPY bin/prometheus_webhook_snmptrapper /prometheus_webhook_snmptrapper
COPY sample-alert.json /

ENTRYPOINT ["/bin/bash", "-c", "/prometheus_webhook_snmptrapper \"$@\"", "--"]