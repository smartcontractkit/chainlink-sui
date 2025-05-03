# Chainlink Sui

- [Chainlink Sui](#chainlink-sui)
  - [Development and Contribution](#development-and-contribution)
    - [Prerequisites](#prerequisites)
    - [Developing in a stable dev environment](#developing-in-a-stable-dev-environment)
    - [Running Tasks](#running-tasks)
  - [Get something running](#get-something-running)
    - [Prerequisites](#prerequisites-1)
    - [Running Sui (Local) Dev Net](#running-sui-local-dev-net)
    - [Deploying Sample Contracts](#deploying-sample-contracts)
  - [Relayer](#relayer)


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

To make sure that you're `sui` environment is ready, you can first check the active
address using the cli

```bash
(nix:nix-shell-env) $ sui client active-address
```

Then you can proceed to make sure that you can the `local` RPC available in the
list of sui cli environments

```bash
(nix:nix-shell-env) $ sui client envs
```

And if you don't see "local", you can add it as follows

```bash
sui client new-env --alias local --rpc http://127.0.0.1:9000
sui client switch --env local
```

### Running Tasks

We use [Task](https://taskfile.dev/) to execute development tasks. You can find every task referenced in the [Taskfile](./Taskfile.yml)

```bash
(nix:nix-shell-env) $ task lint
(nix:nix-shell-env) $ task lint:fix
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

## Relayer

For detailed documentation about using the Sui Relayer Plugin, including how to configure and use the ChainReader and ChainWriter components, see [RELAYER.md](./RELAYER.md).
