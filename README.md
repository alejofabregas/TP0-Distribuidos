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
