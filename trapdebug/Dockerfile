FROM alpine:3.6
USER root

ENV LISTEN_PORT=162

EXPOSE 162/udp

COPY bin/prometheus_webhook_snmptrapper-debugger /usr/local/bin/trapdebug

CMD exec /usr/local/bin/trapdebug -listenport=$LISTEN_PORT

# docker build -t "prawn/snmp-trapdebug" .
