# Sui Integration

This directory contains the Sui-Chainlink integration.

## Prerequisites

- Docker
- Docker Compose
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (or [OrbStack](https://orbstack.dev/))
- Go (1.23+)

## Getting Started

### Running Sui (Local) Dev Net

1. Pull the repo and change directory into `/sui`
2. Run `docker compose up` to get the `sui` devnet up and running
3. You can now `exec` into the container by running `docker compose exec -it sui bash`
4. Run `sui client envs` to view the available sui environments

> NOTE: You can view the open ports and the commands used by inspecting the `/sui/docker-compose.yml` file. It is a work in-progress and will likely change.


## TODO:

- [ ] Add the chainlink local environment to docker compose along with the necessary Dockerfiles
- [ ] Write a basic script to upload a test contract onto the local devnet
- [ ] Write a basic version of `chainreader` and `chainwriter`