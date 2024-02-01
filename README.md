# metal-stack-cloud cli

[![Markdown Docs](https://img.shields.io/badge/markdown-docs-blue?link=https%3A%2F%2Fgithub.com%2Fmetal-stack-cloud%2Fcli%2Fdocs)](./docs/metal.md)

This is the official CLI for accessing the API of [metalstack.cloud](https://metalstack.cloud).

To work with this CLI, it is first necessary to create an api-token. This can be issued through the [cloud console](https://console.metalstack.cloud/token).

Once you got the token, you probably want to create a CLI context. The project's ID can be copied from the UI, the button is located right next to the poject title in the project dashboard.

```bash
$ metal ctx add devel --activate --default-project <project-uuid> --api-token <your-token>
✔ added context "devel"
```

The configuration file is by default written to `~/.metal-stack-cloud/config.yaml`.

The generated markdown documentation of all the commands can be found [here](./docs/metal.md).

## Installation

Download locations:

- [metal-linux-amd64](https://github.com/metal-stack-cloud/cli/releases/latest/download/metal-linux-amd64)
- [metal-darwin-amd64](https://github.com/metal-stack-cloud/cli/releases/latest/download/metal-darwin-amd64)
- [metal-darwin-arm64](https://github.com/metal-stack-cloud/cli/releases/latest/download/metal-darwin-arm64)
- [metal-windows-amd64](https://github.com/metal-stack-cloud/cli/releases/latest/download/metal-windows-amd64)

### Installation on Linux

```bash
curl -LO https://github.com/metal-stack-cloud/cli/releases/latest/download/metal-linux-amd64
chmod +x metal-linux-amd64
sudo mv metal-linux-amd64 /usr/local/bin/metal
```

### Installation on MacOS

For x86 based Macs:

```bash
curl -LO https://github.com/metal-stack-cloud/cli/releases/latest/download/metal-darwin-amd64
chmod +x metal-darwin-amd64
sudo mv metal-darwin-amd64 /usr/local/bin/metal
```

For Apple Silicon (M1) based Macs:

```bash
curl -LO https://github.com/metal-stack-cloud/cli/releases/latest/download/metal-darwin-arm64
chmod +x metal-darwin-arm64
sudo mv metal-darwin-arm64 /usr/local/bin/metal
```

### Installation on Windows

```bash
curl -LO https://github.com/metal-stack-cloud/cli/releases/latest/download/metal-windows-amd64
copy metal-windows-amd64 metal.exe
```
