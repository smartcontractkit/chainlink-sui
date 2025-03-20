# Chainlink Sui

## Development and Contribution

### Prerequisites

**Nix**

Install Nix using the Determinate Systems [installer](https://github.com/DeterminateSystems/nix-installer) to get Nix with Flakes installed.

```bash
$ curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install
```

### Developing in a stable dev environment

The repository comes with a developer environment (devShell) which can be accessed by running:

```bash
$ nix develop
# to exit the shell enviroment enter "exit" in your shell
```

The devShell provides all the system tools and dependencies required to develop and run the project

```bash
(nix:nix-shell-env) $ go version # go version go1.23.7
(nix:nix-shell-env) $ sui version # sui 1.44.3-615516edb0ed
```

## Get something running

### Prerequisites

- Docker
- Docker Compose
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) (or [OrbStack](https://orbstack.dev/))

### Running Sui (Local) Dev Net

1. Run `docker compose up` to get the `sui` devnet up and running
2. You can now `exec` into the container by running `docker compose exec -it sui bash`
3. Run `sui client envs` to view the available sui environments

> NOTE: You can view the open ports and the commands used by inspecting the `/sui/docker-compose.yml` file. It is a work in-progress and will likely change.

### Deploying Sample Contracts

Once you have the local Sui devnet running, you can deploy the sample contracts using:

```bash
./scripts/deploy_contracts.sh
```

This will build and deploy the contracts in `contracts/test` to your local Sui network.
