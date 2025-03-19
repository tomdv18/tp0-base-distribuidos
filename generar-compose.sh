#!/bin/bash

# Valores por defecto
ARCHIVO="docker-compose-dev-tom.yaml"
CLIENTES=3

if [ -z "$1" ]; then
  echo "No se proporcionó nombre de archivo, usando valor por defecto: $ARCHIVO"
  archivo_salida=$ARCHIVO
else
  archivo_salida=$1
  
  if [[ "$archivo_salida" != *.yaml ]]; then
    archivo_salida="$archivo_salida.yaml"
  fi
fi


if [ -z "$2" ]; then
  echo "No se proporcionó cantidad de clientes, usando valor por defecto: $CLIENTES"
  cantidad_clientes=$CLIENTES

else
  if ! [[ "$2" =~ ^[0-9]+$ ]]; then
    echo "Error: La cantidad de clientes debe ser un número entero positivo."
    exit 1
  fi
  cantidad_clientes=$2
fi


echo "Nombre del archivo de salida: $archivo_salida"
echo "Cantidad de clientes: $cantidad_clientes"


python3 generador.py $archivo_salida $cantidad_clientes
