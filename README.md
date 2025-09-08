# distributed-aviation-system

This project is a distributed system that re-creates some functionality of existing aviation systems.

# Buf

[Buf](https://buf.build/home) is used to manage Protobuf definitions and generate client and server stubs for our
services. Buf is configured in the [buf.yml](buf.yml) and the code generation is in [buf.gen.yml](buf.gen.yml).

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
  buf generate
```