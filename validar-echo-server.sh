#!/bin/bash

DIRECCION_SERVER=$(grep SERVER_IP server/config.ini | cut -d ' ' -f 3)
PUERTO_SERVER=$(grep SERVER_PORT server/config.ini | cut -d ' ' -f 3)
MENSAJE="HOLA MUNDO"

make docker-compose-up

docker run --name client-netcat -it --network tp0_testing_net ubuntu bash -c "apt-get update && apt-get install --assume-yes netcat-openbsd"
docker start client-netcat

RESPUESTA=$(docker exec client-netcat bash -c "echo $MENSAJE | netcat $DIRECCION_SERVER $PUERTO_SERVER")

docker stop client-netcat
docker rm client-netcat

make docker-compose-down

if [ "$RESPUESTA" = "$MENSAJE" ];
then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi