# Development Certificates Generator

devcert takes away the pain of manually creating self-signed certificates for development.

------------

:warning: **Note: These certificates are NOT meant to be used on any server other than your local development machine. These certificates are NOT secure and the generated certificate authority by this tool is NOT trusted by browser vendors.**

-------------


![devcert-photo](https://user-images.githubusercontent.com/489775/167084056-4cf4a8f8-ff49-4ccc-b5de-a3c110ccbd01.png)

## Installation

Grab a [pre-built binary](https://github.com/primalskill/devcert/releases).

OR

Clone this repo and compile from source using Go.

### Install a pre-built binary

1. Download the binary for your platform, example macOS ARM: `curl https://github.com/primalskill/devcert/releases/download/v1.1.2/devcert_darwin_arm64 > /usr/local/bin/devcert`
2. Make it an executable: `chmod +x /usr/local/bin`
3. Generate a certificate for a local domain (see the detailed usage below): `devcert example.test`


### Compile from Source

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

### Supported Platforms

- macOS
- Windows
- Linux (Debian, Ubuntu, OpenSUSE, RHEL, CentOS, Fedora, Arch Linux)


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

## How It Works

All the certificates created by devcert will be placed in the `~/.devcert` directory.

Running devcert for the first time will execute the setup process which will:

1. Create the `~/.devcert` directory
2. Create a local certificate authority (CA) used to sign other domain specific certificates.
3. It will mark the CA as trusted automatically.

Once the setup process is completed it will generate the domain specific certificate. You can generate as many self-signed, trusted, local certificates for development as you like, the `.crt` and `.key` files will be placed in the `~/.devcert` directory.


## Known Issues

### Fixing `SEC_ERROR_REUSED_ISSUER_AND_SERIAL` in Firefox

If you are getting this error, it's most likely Firefox preloaded a previously generated certificate authority (CA) in the default browser profile. This happens if the devcert CA files are manually removed and generated again.

To fix it:

1. Close all instances of Firefox
2. Go in the profile folder
  - Windows: `C:\Users\%userprofile%\AppData\Roaming\Mozilla\Firefox\Profiles\%profile.default%`
  - MacOS: `~/Library/Application Support/Firefox/Profiles/<profile folder>`
3. Remove the files `cert8.db`, `cert9.db`, `cert_override.txt` (Note: some of these files may not exist).
