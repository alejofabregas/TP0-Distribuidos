#!/bin/bash
MENSAJE="HOLA MUNDO"

docker run --name client-netcat -it --network tp0_testing_net ubuntu bash -c "apt-get update && apt-get install --assume-yes netcat-openbsd"
docker start client-netcat

RESPUESTA=$(docker exec -it client-netcat bash -c "echo HOLA MUNDO | netcat server 12345")

docker stop client-netcat
docker rm client-netcat

echo "Mensaje enviado: $MENSAJE"
echo "Respuesta recibida: $RESPUESTA"