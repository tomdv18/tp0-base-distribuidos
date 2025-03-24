import sys
import yaml

def generar_docker_compose(archivo_salida, cantidad_clientes):
    docker_compose = {
        "name": "tp0",
        "services": {
            "server": {
                "container_name": "server",
                "image": "server:latest",
                "entrypoint": "python3 /main.py",
                "environment": [
                    "PYTHONUNBUFFERED=1",
                ],
                "networks": ["testing_net"],
                "volumes" : [
                "./server/config.ini:/config.ini"
                ],
            }
        },
        "networks": {
            "testing_net": {
                "ipam": {
                    "driver": "default",
                    "config": [{"subnet": "172.25.125.0/24"}]
                }
            }
        }
    }

    for i in range(1, cantidad_clientes + 1):
        nombre = f"client{i}"
        docker_compose["services"][nombre] = {
            "container_name": nombre,
            "image": "client:latest",
            "entrypoint": "/client",
            "environment": [
                f"CLI_ID={i}",
                f"NOMBRE={nombre}",
                f"APELLIDO={nombre}",
                f"DOCUMENTO={i+23322510}",
                f"NACIMIENTO=1985-10-18",
                f"NUMERO={(i+7)*4}",
            ],
            "networks": ["testing_net"],
            "volumes": [
                "./client/config.yaml:/config.yaml"
            ],
            "depends_on": ["server"]
        }

    with open(archivo_salida, 'w') as file:
        yaml.dump(docker_compose, file, default_flow_style=False)

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Error: Se deben proporcionar dos parámetros: el nombre del archivo y la cantidad de clientes.")
        sys.exit(1)

    archivo_salida = sys.argv[1]
    try:
        cantidad_clientes = int(sys.argv[2])
    except ValueError:
        print("Error: La cantidad de clientes debe ser un número entero válido.")
        sys.exit(1)

    generar_docker_compose(archivo_salida, cantidad_clientes)
