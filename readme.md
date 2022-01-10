# Development Certificates Generator

devcert takes away the pain of creating self-signed certificates for development manually.

## Usage

```
$ devcert my-project.local api.my-project.local my-other-project.test

Generating certificate...
Generated at:
  Certificate: ~/.devcert/devcert_my-project.local_multi.crt
  Private Key: ~/.devcert/devcert_my-project.local_multi.key

Valid for:
  1. my-project.local
  2. api.my-project.local
  3. my-other-project.test
```

You can move the `.crt` and `.key` files to your desired location. It will be signed with the CA, no need to trust this certificate separately.

## On First Run

When running the program for the first time, it will ask for running the setup process which creates the necessary directory, generate the CA and mark it as trusted.

This is a one time process that needs to be executed before generating domain specific certificates.

Example:

```
$ devcert myapp.local

devcert needs to execute the setup process first.
  - It will create ~/.devcert/ directory.
  - It will create a local certificate authority (CA) to sign future certificates.
  - It will mark the CA as trusted locally.
Do you want to continue? [Y/n]: Y

Creating directory...
Directory ~/.devcert/ created.
Creating certificate authority (CA) files...
Certificate authority (CA) created at
  Certificate: ~/.devcert/devcert_ca.crt
  Private Key: ~/.devcert/devcert_ca.key
Trusting certificate authority...
Certificate authority (CA) marked trusted.
```

**Note: The certificate authority (CA) `.crt` and `.key` files should be left in the `~/.devcert` directory as these files will be loaded when generating a domain specific certificate.**

## Installation

Grab a [pre-built binary](https://github.com/primalskill/devcert/releases).

OR

Clone this repo and compile from source using Go.

## Compile from Source

Prerequisites:

- Go
- Make

Execute `make release-<desired platform and architecture>`. Make will create the binary in `./.bin` directory.

Available make commands:

- `make release-win-amd64`
- `make release-darwin-amd64`
- `make release-darwin-arm64`
- `make release-linux-amd64`
- `make release-linux-arm64`

## Supported Platforms

- macOS
- Windows
- Linux (Debian, Ubuntu, OpenSUSE, RHEL, CentOS, Fedora, Arch Linux)


## How It Works

All the certificates created by devcert will be placed in the `~/.devcert` directory.

Running devcert for the first time will execute the setup process which will:

1. Create the `~/.devcert` directory
2. Create a local certificate authority (CA) used to sign other domain specific certificates.
3. It will mark the CA as trusted automatically.

Once the setup process is completed it will generate the domain specific certificate. You can generate as many self-signed, trusted, local certificates for development as you like, the `.crt` and `.key` files will be placed in the `~/.devcert` directory.
