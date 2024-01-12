# metal-stack-cloud cli

To work with this CLI, it is first necessary to create a metalstack.cloud api-token. This can be issues through the cloud console.

Once you got the token, you probably want to create a CLI context:

```bash
$ metal ctx add project-xyz --api-token <your-token> --default-project project-xyz --activate
✔ added context "project-xyz"
```

The configuration file is by default written to `~/.metal-stack-cloud/config.yaml`.

## Basic commands

Documentation on how to interact with the CLI (maybe just necessary during development):

### IP

```bash
# list ips
$ bin/metal ip list --project <project-id>
```

```bash
# create ip from cli
$ bin/metal ip create --project <project-id> --name <name> --network <network>
```

```bash
# create ip with file option
#
#
# ip.yaml file:
# name: <ip-name>
# network: <network>
# project: <project-id>
# description: <description>
# type: <ephemeral | static>
#
#

$ bin/metal ip create -f <file-name>
```

```bash
# describe ip
$ bin/metal ip describe --project <project-id> <ip-uuid>
```

```bash
# update command to make the ip static
$ bin/metal ip update --project <project-id> --uuid <ip-uuid>
```

```bash
# update ip with file option
#
#
# ip.yaml file:
# project: <project-id>
# uuid: <ip-uuid>
#
#

$ bin/metal ip update -f <file-name>
```

### Admin

```bash
# list all tenants
$ bin/metal admin tenant list

# admit a tenant
$ bin/metal admin tenant admit <tenant-id>

# revoke a tenant
$ bin/metal admin tenant revoke <tenant-id>

# list all coupons
$ bin/metal admin coupon list
```
