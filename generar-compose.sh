#!/bin/bash

# Chequear parametros
if [ "$#" -ne 2 ]; then
    echo "[ERROR] Usar con exactamente 2 parámetros: ./generar-compose.sh ${NOMBRE_ARCHIVO_SALIDA.yaml} ${CANTIDAD_CLIENTES}"
    exit 1
fi

# echo "Nombre del archivo de salida: $1"
# echo "Cantidad de clientes: $2"

python3 subscript-generar-compose.py $1 $2