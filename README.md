# TP0: Docker + Comunicaciones + Concurrencia

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
