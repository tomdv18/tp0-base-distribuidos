#!/bin/bash

# Archivo de configuracion del servidor
CONFIG_FILE="server/config.ini"

SERVER=$(grep "^SERVER_IP" "$CONFIG_FILE" | sed -E 's/^SERVER_IP *= *//')
PORT=$(grep "^SERVER_PORT" "$CONFIG_FILE" | sed -E 's/^SERVER_PORT *= *//')

MENSAJE="Mensaje de Prueba"


if [ -z "$SERVER" ] || [ -z "$PORT" ]; then
  echo "Error: No se pudo obtener SERVER_IP o SERVER_PORT de $CONFIG_FILE"
  exit 1
fi

SERVER_DIR="${SERVER}:${PORT}"
echo "Probando conexión con $SERVER_DIR"

RESPONSE=$(echo "$MENSAJE" | docker run --rm --network host busybox nc "$SERVER" "$PORT")

if [ "$RESPONSE" = "$MENSAJE" ]; then
  echo "action: test_echo_server | result: success | address: $SERVER_DIR"
else
  echo "action: test_echo_server | result: fail | address: $SERVER_DIR"
fi
