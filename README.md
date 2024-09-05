# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker

### Ejercicio N°1:
Para ejecutar el primer ejercicio simplemente hay que correr el script de bash en una terminal en el root del proyecto:

```
./generar-compose.sh ${NOMBRE_ARCHIVO_SALIDA.yaml} ${CANTIDAD_CLIENTES}
```

Este script invoca al subscript `subscript-generar-compose.py` que toma los argumentos y genera el archivo de Docker Compose en la ruta indicada.
El script de bash valida los argumentos, sólo funciona si se le pasan los dos indicados.

### Ejercicio N°2:
Se modifica el subscript para generar el archivo de Docker Compose para que los archivos de configuración se inyecten en sus 
respectivos containers. Para ello se hacen bind mounts: uno para el server y uno para cada cliente. Así, no hay que hacer 
build de las imágenes cada vez que se cambia la configuración.

Para ver que efectivamente se hayan hecho los bind mounts, podemos ejecutar el siguiente comando:
```
docker inspect ${CONTAINER_ID}
```
En la sección de `Mounts` veremos que tenemos lo siguiente:
```
"Type": "bind",
...
"Destination": "/config.ini",
```
Lo que nos indica que efectivamente el container tiene el archivo de configuración montado.

Esto también se puede ver en Docker Desktop si entramos a los detalles de alguno de los containers, en la sección `Bind mounts`.

### Ejercicio N°3:

Para ejecutar el tercer ejercicio simplemente hay que correr el script de bash en una terminal en el root del proyecto:

```
./validar-echo-server.sh
```

El script primero levanta una instancia del server, luego levanta un container con una imagen a la que le agregamos el comando `netcat`.
Luego ejecutamos el comando netcat con la dirección y el puerto del server (levantados de su archivo de configuración), y con un pipe
le enviamos un mensaje de prueba que esperamos recibir de nuevo del echo server. Una vez hecho esto, se detienen ambos containers.
Finalmente, se compara el mensaje de prueba enviado y la respuesta recibida del server. Si son iguales, se imprime el mensaje de éxito; 
si no lo son, se imprime el mensaje de fracaso.

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

### Ejercicio N°6:

En este ejercicio se implementó el envío y recepción de apuestas en forma de batch entre el cliente y servidor. 
Se puede observar el comportamiento viendo los logs utilizando:
```
make docker-compose-up
make docker-compose-logs
make docker-compose-down
```

También se puede ver que el servidor recibe el total de las apuestas correspondiente a los archivos `agency-<n>.csv` de cada cliente 
en el archivo `bets.csv` dentro del container del servidor.
El total de apuestas recibidas por el servidor es de `78697`, lo que corresponde a la suma de las apuestas de todos los clientes.

### Protocolo de comunicación

Se utiliza un esquema muy similar al de ejercicio anterior, anteponiendo la cantidad de bytes a leer. 
La diferencia es que ahora se encadenan todas las bets de un batch una detrás de la otra en un mismo buffer de bytes. 
Estas bets van separadas con un \n.

En cuanto a la variable `BatchMaxAmount` del archivo de configuración del cliente, consideré usar `135` bets por batch. 
Esto es porque la bet de mayor largo tiene 59 caracteres = 59 bytes, por lo que si usamos 135 bets de este tamaño 
no nos excedemos de los 8 kb = 8000 bytes de largo máximo de un paquete requerido por la cátedra.

### Ejercicio N°7:

En este ejercicio se modificó tanto al cliente como al servidor para soportar el sorteo de las apuestas. 
```
make docker-compose-up
make docker-compose-logs
make docker-compose-down
```

Los clientes, luego de enviar todos sus batch de apuestas, envían un mensaje de `Finish` al servidor, que consiste 
únicamente en enviar un mensaje que contenga el uint32 que utilizamos antes para indicar la cantidad de bytes a leer, 
pero con valor `0`. Así, el servidor no sólo sabe que no tiene que leer más, sino que le indica que el cliente finalizó. 
A continuación, el cliente abre una nueva conexión con el servidor, pero se queda bloqueado esperando a recibir los resultados 
del sorteo que realiza el servidor. Una vez que le llegan, los notifica.

El servidor implementa una estrategia de `long polling` para hacer el sorteo. Cada vez que le llega un mensaje `Finish` de un cliente, 
lo agrega a un diccionario con su address y socket. Cuando le lleguen tantos mensaje `Finish` como clientes haya (el diccionario tiene
una cantidad de elementos igual a la cantidad de clientes), sabe que todos terminaron y realiza el sorteo. Finalmente, notifica uno 
por uno a los clientes que se habían quedando esperando por el resultado del sorteo.

Ganadores totales:
22737492|23328212|28188111|30876370|34963649|35635602|33791469|24807259|31660107|24813860

Ganadores agencia 1:
30876370|24807259

Ganadores agencia 2:
33791469|31660107|24813860

Ganadores agencia 3:
22737492|23328212|28188111

Ganadores agencia 4:
34963649|35635602

Ganadores agencia 5:
No hay

## Parte 3: Repaso de Concurrencia

### Ejercicio N°8:

En este ejercicio se modificó tanto al cliente como al servidor para soportar el procesamiento paralelo de mensajes en el servidor.
Se puede observar el comportamiento viendo los logs utilizando:
```
make docker-compose-up
make docker-compose-logs
make docker-compose-down
```

Básicamente el servidor utiliza `procesos` para atender a cada cliente. Cada vez que se abre una nueva conexión con un cliente, 
se crea un nuevo proceso que se va a encargar de atenderla. Así, cada cliente se atiende en un proceso distinto de manera tal 
que es un procesamiento en **paralelo**.

Inicialmente comencé haciendo una implementación que usaba una `ThreadPool`, pero luego cambié a utilizar procesos ya que no 
sufren del problema del `GIL`. De cualquier manera, como el servidor realiza un procesamiento muy intensivo de IO pero no tanto 
de CPU, no sería un gran problema si se utilizan threads en lugar de procesos.

Para sincronizar el acceso a las funciones que no son thread-safe para la persistencia de la información de las apuestas, 
se utilizan `locks` para que sólo se pueda acceder de a un proceso. Así, evitamos cualquier tipo de problema.

Para mantener información compartida entre procesos, en nuestro caso el diccionario con clientes que ya finalizaron, se utiliza 
un `Manager` de la biblioteca de multiprocessing. Este `Manager` nos permite compartir el diccionario de forma segura entre procesos, 
mediante el mecanismo de serialización de `pickle` que ofrece Python. Así, podemos tratar al diccionario como memoria compartida, 
casi sin enterarnos que se está compartiendo entre procesos.
