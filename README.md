# bed
Don't forget to set your spawn!

Bed is automated daily backups for Minecraft servers utilizing Minio as the storage bucket

# Disclaimer
this repo is unfinished, the goal is for this to be a docker image, I recommend waiting for that release. This disclaimer will be replaced by usage docs and the improved compose file

<hr />

# Dev notes

# Example compose file with bed
## Do not use this yet, the image is not built (yet)
```
version: "3.8"

services:
  mc-container:
    image: itzg/minecraft-server
    tty: true
    std_open: true
    ports: 
      - "25565:25565"
      - "61695:61695"
    volumes:
      - mc-container:/data
    environment:
      EULA: "TRUE"
      MEMORY: 2G
      TYPE: FORGE
      VERSION: "1.20.1"
      FORGE_VERSION: "47.2.0"
      ENABLE_RCON: "true"
      RCON_PASSWORD: "fakepassword"
      RCON_PORT: "61695"
      DIFFICULTY: "hard"

  bed:
    image: jydv/bed
    volumes:
      - mc-container:/data
    environment:
      RCON_MC_CONTAINER: "mc-container"
      RCON_PASSWORD: "fakepassword"
      RCON_PORT: "61696"
      MINIO_URL: "minio url"
      MINIO_BUCKET: "bucket to save backup"
      MINIO_KEY: "minio access key"
      MINIO_ID: "minio client id"
      SERVER_RESTART: "yes|no"

volumes:
  mc-container:
    driver: local

networks:
  default:
    external: true
    name: mc-container_default
```