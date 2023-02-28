# metal-stack-cloud cli

## Admin usage

In order to use the admin commands you must first create a config file `~/.metal-stack-cloud/config.yaml` with this content.

```yaml
# this config works in the mini-lab
# api-url: http://api.172.17.0.1.nip.io:8080
# this config is working for a api-server started locally
api-url: http://localhost:8080
api-token: eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJtZXRhbC1zdGFjay1jbG91ZCIsInN1YiI6ImFkbWluIiwiZXhwIjo0ODEyNjE0OTczLCJyb2xlcyI6eyIqIjoiYWRtaW4ifX0.gsqlaAcvIZFFYZSxrOMIwiZdKb0AZiGhFt4qpS0keC8
api-ca-file: yourdevelopmentfolder/metal-stack-cloud/deployment/files/certs/ca.pem
```

This config contains a api-token with admin permissions for the development and is not suitable for production use.

After that you can see the available admin subcommands with:

```bash
bin/metal admin
```

## Basic commands

Documentation on how to interact with the CLI (maybe just necessary during development):

### IP

```bash
# list ips
$ bin/metal ip list --project <project-id>

# create ip from cli
$ bin/metal ip create --project <project-id> --name <name> --network <network>

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
or
$ bin/metal ip apply -f <file-name>

# describe ip
$ bin/metal ip describe --project <project-id> <ip-uuid>

# update command to make the ip static
$ bin/metal ip update --project <project-id> --uuid <ip-uuid>

# update ip with file option
#
#
# ip.yaml file:
# project: <project-id>
# uuid: <ip-uuid>
#
#

$ bin/metal ip update -f <file-name>

# delete an ip address
$ bin/metal ip delete --project <project-id> <ip-uuid>
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

### Cluster

```bash
# list all clusters
$ bin/metal cluster list --project <project-id>

# describe a cluster
$ bin/metal cluster describe --project <project-id> <cluster-uuid>

# delete a cluster
$ bin/metal cluster delete --project <project-id> <cluster-uuid>

# create a cluster
$ bin/metal cluster create --name <cluster-name> --project <project-id> --partition <partition> --kubernetes <kubernetes-version> --workername <worker-name> --machinetype <machine-type> --minsize <min-worker> --maxsize <max-worker> --maintenancebegin <maintenance-begin> --maintenanceduration <maintenance-duration>

# create a cluster with file option

# cluster.yaml file
# name: <cluster-name>
# project: <project-id>
# partition: <partition>
# kubernetes:
#   version: <kubernetes-version>
# workers:
#   - name: <worker-name>
#     machinetype: <machine-type>
#     minsize: <min-worker>
#     maxsize: <max-worker>
# maintenance:
#   timewindow:
#     begin:
#       seconds: <maintenance-begin>
#     duration:
#       seconds: <maintenance-duration>

$ bin/metal cluster create -f <file-name>
or
$ bin/metal cluster apply -f <file-name>
```
