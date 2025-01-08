# metal-stack-cloud CLI

[![Markdown Docs](https://img.shields.io/badge/markdown-docs-blue?link=https%3A%2F%2Fgithub.com%2Fmetal-stack-cloud%2Fcli%2Fdocs)](./docs/metal.md)
![Go Version](https://img.shields.io/github/go-mod/go-version/metal-stack-cloud/cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/metal-stack-cloud/cli)](https://goreportcard.com/report/github.com/metal-stack-cloud/cli)

This is the official CLI for accessing the API of [metalstack.cloud](https://metalstack.cloud).

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

## Usage

All commands follow a general form:

```bash
metal <entity> [<category>] <command> <argument> [<flags>]
```

For example:

```bash
metal tenant member list --api-token <your-token> --api-url <api-url>
metal cluster list -p <project-id>
metal ctx add <context-name>
```

The `api-token`, `api-url` and `project-id` are defaulted by the context, if one exists, and can be omitted.

In addition to the standard API services, there are also admin services that require an admin token for execution.

You can access help for every service and command by using `--help` or `-h`. If you encounter any issues not covered in the help prompt, or if you have suggestions for improvement, please feel free to [contact us](mailto:support@metalstack.cloud) or open an issue in this repository. Your feedback is greatly appreciated!

A list of all available services (excluding admin topics). For their associated commands, arguments and flags visit the correct [documentation](./docs/metal.md).

| Entity        | Description                                                | Documentation                                    |
| ------------- | ---------------------------------------------------------- | ------------------------------------------------ |
| `api-methods` | show available api-methods of the metalstack.cloud api     | [metal api-methods](./docs/metal_api-methods.md) |
| `asset`       | show asset                                                 | [metal asset](./docs/metal_asset.md)             |
| `cluster`     | manage cluster entities                                    | [metal cluster](./docs/metal_cluster.md)         |
| `completion`  | generate the autocompletion script for the specified shell | [metal completion](./docs/metal_completion.md)   |
| `context`     | manage cli contexts                                        | [metal context](./docs/metal_context.md)         |
| `health`      | print the client and server health information             | [metal health](./docs/metal_health.md)           |
| `ip`          | manage ip entities                                         | [metal ip](./docs/metal_ip.md)                   |
| `markdown`    | create markdown documentation                              | [metal completion](./docs/metal_completion.md)   |
| `payment`     | manage payment of the metalstack.cloud                     | [metal payment](./docs/metal_payment.md)         |
| `project`     | manage project entities                                    | [metal project](./docs/metal_project.md)         |
| `storage`     | storage commands                                           | [metal storage](./docs/metal_storage.md)         |
| `tenant`      | manage tenant entities                                     | [metal tenant](./docs/metal_tenant.md)           |
| `token`       | manage token entities                                      | [metal token](./docs/metal_token.md)             |
| `user`        | manage user entities                                       | [metal user](./docs/metal_user.md)               |
| `version`     | print the client and server version information            | [metal version](./docs/metal_version.md)         |

### Autocompletion

To successfully set up autocompletion follow this [guide](./docs/metal_completion.md).

## Authentication and Configuration

To work with this CLI, it's necessary to create an api-token. This can be issued through the [cloud console](https://console.metalstack.cloud/token), make sure to configure the right permissions you want to use within your context.

The project's ID can be copied from the UI, the button is located right next to the project title in the project dashboard. The default API-URL of metal-stack-cloud is https://api.metalstack.cloud.

```bash
$ metal ctx add <context-name> --activate --default-project <project-uuid> --api-token <your-token>
✔ added context "<context-name>"
```

The configuration file is by default written to `~/.metal-stack-cloud/config.yaml`.
