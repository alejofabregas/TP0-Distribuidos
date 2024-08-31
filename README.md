# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker

### Ejercicio N°3:

Para ejecutar el primer ejercicio simplemente hay que correr el script de bash en una terminal en el root del proyecto:

```
./validar-echo-server.sh
```

El script primero levanta una instancia del server, luego levanta un container con una imagen a la que le agregamos el comando `netcat`.
Luego ejecutamos el comando netcat con la dirección y el puerto del server (levantados de su archivo de configuración), y con un pipe
le enviamos un mensaje de prueba que esperamos recibir de nuevo del echo server. Una vez hecho esto, se detienen ambos containers.
Finalmente, se compara el mensaje de prueba enviado y la respuesta recibida del server. Si son iguales, se imprime el mensaje de éxito; 
si no lo son, se imprime el mensaje de fracaso.