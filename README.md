# bed
Don't forget to set your spawn!

# Disclaimer
this repo is unfinished, it may work when ran with docker compose but can experience breaking changes until a stable release. The end-game goal for this repo is to be a docker image that can be added to Minecraft server compose files for easy setup for world backups. If you are seeing this disclaimer I wouldn't recommend using this version of this software and to wait for the full release

<hr />

# Dev notes

# TODO
1. allow for users to change how often backups are grabbed (timezone and clock time, maybe), how many to keep (maybe)
2. use iota to create a loop for env validation

# Example compose file with bed
## Do not use this yet, the image is not built
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

volumes:
  mc-container:
    driver: local

networks:
  default:
    external: true
    name: mc-container_default
```