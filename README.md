# TP0: Docker + Comunicaciones + Concurrencia

## Parte 2: Repaso de Comunicaciones

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