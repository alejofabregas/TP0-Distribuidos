# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker

### Ejercicio N°1:
Para ejecutar el primer ejercicio simplemente hay que correr el script de bash en una terminal en el root del proyecto:

```
./generar-compose.sh ${NOMBRE_ARCHIVO_SALIDA.yaml} ${CANTIDAD_CLIENTES}
```

Este script invoca al subscript `subscript-generar-compose.py` que toma los argumentos y genera el archivo de Docker Compose.