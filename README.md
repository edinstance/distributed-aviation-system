# distributed-aviation-system

This project is a distributed system that re-creates some functionality of existing aviation systems.

# Graphql

This project uses [Graphql](https://graphql.org/), each service exposes its own set of Graphql Application Programming
Interface's (API's), and there is a routing layer that
uses [Apollo Router](https://www.apollographql.com/docs/graphos/routing) for routing to the correct subgraph.

## Apollo Router

The project runs the router in a [Docker](https://www.docker.com/) container, but the configuration of the router is
done in the [router.yml](router.yml) file and it uses this [supergraph](supergraph.graphql).

## Supergraph

The routers use a single federated graph or supergraph which is made from all of the individual services graphs or
sub-graphs. The supergraph can be created using the [Apollo Rover CLI](https://www.apollographql.com/docs/rover).

[!CAUTION]
Before using the router make sure that the supergraph has been created or updated.

### Rover installation

Either follow [this guide](https://www.apollographql.com/docs/rover/getting-started) to get started or follow the
repository-specific steps below.

1. Install rover

You can either use a CLI or npm.

```shell
curl -sSL https://rover.apollo.dev/nix/latest | sh
```

```shell
npm install -g @apollo/rover
```

2. Usage

There are many ways to run an application using rover and some can be
found [here](https://www.apollographql.com/docs/rover). This project uses rover to create a supergraph, to configure the
creation of the supergraph use the [supergraph.yml](supergraph.yml) file. To create the supergraph use this command.

```shell
rover supergraph compose --config ./supergraph.yml --output supergraph.graphql
```

# Buf

[Buf](https://buf.build/home) is used to manage Protobuf definitions and generate client and server stubs for our
services. Buf is configured in the [buf.yml](buf.yaml) and the code generation is in [buf.gen.yml](buf.gen.yaml).

## Installation

Please follow the [official instructions](https://buf.build/docs/cli/installation/) to install buf. Then verify it is
installed correctly using

```shell
  buf --version
```

## Usage

1. Linting of proto files

```shell
  buf lint
```

2. Generating code stubs

```shell
  buf generate --template buf.gen.yaml
```