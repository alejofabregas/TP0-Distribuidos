import sys
import yaml

def generar_docker_compose(nombre_archivo, cantidad_clientes):
    docker_compose = {
        "name": "tp0",
        
        "services": {
            "server": {
                "container_name": "server",
                "image": "server:latest",
                "entrypoint": "python3 /main.py",
                "environment": [
                    "PYTHONUNBUFFERED=1",
                    "LOGGING_LEVEL=DEBUG"
                ],
                "networks": ["testing_net"],
                "volumes": ["./server/config.ini:/config.ini"]
            },
        },
        
        "networks": {
            "testing_net": {
                "ipam": {
                    "driver": "default",
                    "config": [
                        {"subnet": "172.25.125.0/24"}
                    ]
                }
            }
        }
            
    }

    for i in range(1, cantidad_clientes + 1):
        nombre_servicio = f"client{i}"
        docker_compose["services"][nombre_servicio] = {
            "container_name": f"client{i}",
                "image": "client:latest",
                "entrypoint": "/client",
                "environment": [
                    f"CLI_ID={i}",
                    "CLI_LOG_LEVEL=DEBUG"
                ],
                "networks": ["testing_net"],
                "depends_on": ["server"],
                "volumes": ["./client/config.yaml:/config.yaml"]

        }

    with open(nombre_archivo, "w") as archivo:
        yaml.dump(docker_compose, archivo, default_flow_style=False, sort_keys=False)

    # print("\nSubscript de Python")
    # print("Archivo de Docker Compose generado:", nombre_archivo)
    # print("Cantidad de clientes:", cantidad_clientes)


if __name__ == "__main__":
    
    nombre_archivo = sys.argv[1]
    cantidad_clientes = int(sys.argv[2])

    generar_docker_compose(nombre_archivo, cantidad_clientes)