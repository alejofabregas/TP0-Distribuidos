# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker

### Ejercicio N°4:

En este ejercicio se hace el graceful shutdown tanto del servidor como del cliente. Para poder verlo en acción corremos el proyecto 
como antes y lo finalizamos antes de que termine:

```
make docker-compose-up
make docker-compose-down
```

Para ver los logs que verifican el funcionamiento del shutdown podemos quitar un comentario del archivo Makefile para que se corran los 
logs cuando se ejecuta docker-compose-down:

```
docker-compose-down:
	docker compose -f docker-compose-dev.yaml stop -t 1
#	make docker-compose-logs
	docker compose -f docker-compose-dev.yaml down
.PHONY: docker-compose-down
```

La otra opción es usar los logs a mano en la terminal viendo primero el ID del container, y luego accediendo a sus logs:

```
docker ps -a
docker logs -t ${CONTAINER_ID}
```