# TP0: Docker + Comunicaciones + Concurrencia

## Parte 2: Repaso de Comunicaciones

### Ejercicio N°5:

En este ejercicio se implementó el envío y recepción de apuestas entre cliente y servidor. Se puede observar el comportamiento 
viendo los logs utilizando:
```
make docker-compose-up
make docker-compose-logs
make docker-compose-down
```

También podemos ver cómo se guardan las apuestas en el servidor accediendo al archivo bets.csv dentro de su container.

Se implementan funciones read_all y write_all para leer y escribir del socket sin el error de short read y short write.

### Protocolo de comunicación

Se utiliza un esquema mixto. Los mensajes tienen la siguiente forma:
```
<LARGO_STRING>"<AGENCY_ID>|<FIRST_NAME>|<LAST_NAME>|<DOCUMENT_ID>|<BIRTH_DATE>|<BET_NUMBER>\n"
```
Primero tenemos el largo de la string del mensaje, que es un `uint32` que indica cuántos bytes más hay que leer. 
Esto es lo primero que va en el mensaje, tiene un largo fijo de 4 bytes que se leen siempre primero y se encodean en BigEndian. 
Luego tenemos una string que usa como separador el caracter `|`. Se encuentran todos las partes de una Bet hasta un último `\n`.
