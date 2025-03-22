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

RESPONSE=$(echo "$MENSAJE" | docker run --rm --network tp0-tom_testing_net busybox nc "$SERVER" "$PORT")



if [ "$RESPONSE" = "$MENSAJE" ]; then
  echo "action: test_echo_server | result: success | address: $SERVER_DIR"
else
  echo "action: test_echo_server | result: fail | address: $SERVER_DIR"
fi
