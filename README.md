# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker

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