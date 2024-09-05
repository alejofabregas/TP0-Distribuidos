# TP0: Docker + Comunicaciones + Concurrencia

## Parte 2: Repaso de Comunicaciones

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


La idea sería agregar a los clientes que terminaron a un set. Cuando el set llega a la cantidad de clientes, todos terminaron.

Los clientes terminan cuando mandan un mensaje con el uint32 de length igual a 0 bytes.

La cantidad de clientes la podemos definir en una variable de entorno en el docker compose.

El cliente cuando termina sale del loop y ahí envia el 0. Se queda bloqueado esperando una respuesta.

El servidor cuando detecta que un cliente termina (le llega el 0), lo mete al set de clientes terminados.
Cuando el len del set llega a la cantidad de clientes, le responde uno por uno a todos con los ganadores.

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