#!/bin/bash

DIRECCION_SERVER=$(awk -F' = ' '/SERVER_IP/ {print $2}' server/config.ini | tr -d '[:space:]')
PUERTO_SERVER=$(awk -F' = ' '/SERVER_PORT/ {print $2}' server/config.ini | tr -d '[:space:]')
MENSAJE="HOLA MUNDO"

# make docker-compose-up

RESPUESTA=$(docker run --rm --name client-netcat --network tp0_testing_net busybox sh -c "echo $MENSAJE | nc $DIRECCION_SERVER $PUERTO_SERVER")

# make docker-compose-down

if [ "$RESPUESTA" = "$MENSAJE" ];
then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi