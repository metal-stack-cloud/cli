## metal cluster

manage cluster entities

### Synopsis

manage kubernetes clusters

### Options

```
  -h, --help   help for cluster
```

### Options inherited from parent commands

```
      --api-ca-file string     the path to the ca file of the api server
      --api-token string       the token used for api requests
      --api-url string         the url to the metalstack.cloud api (default "https://api.metalstack.cloud")
  -c, --config string          alternative config file path, (default is ~/.metal-stack-cloud/config.yaml)
      --debug                  debug output
      --force-color            force colored output even without tty
  -o, --output-format string   output format (table|wide|markdown|json|yaml|template|jsonraw|yamlraw), wide is a table with more columns, jsonraw and yamlraw do not translate proto enums into string types but leave the original int32 values intact. (default "table")
      --template string        output template for template output-format, go template format. For property names inspect the output of -o json or -o yaml for reference.
      --timeout duration       request timeout used for api requests
```

### SEE ALSO

* [metal](metal.md)	 - cli for managing entities in metal-stack-cloud
* [metal cluster apply](metal_cluster_apply.md)	 - applies one or more clusters from a given file
* [metal cluster create](metal_cluster_create.md)	 - creates the cluster
* [metal cluster delete](metal_cluster_delete.md)	 - deletes the cluster
* [metal cluster describe](metal_cluster_describe.md)	 - describes the cluster
* [metal cluster edit](metal_cluster_edit.md)	 - edit the cluster through an editor and update
* [metal cluster kubeconfig](metal_cluster_kubeconfig.md)	 - fetch kubeconfig of a cluster
* [metal cluster list](metal_cluster_list.md)	 - list all clusters
* [metal cluster monitoring](metal_cluster_monitoring.md)	 - fetch endpoints and access credentials to cluster monitoring
* [metal cluster update](metal_cluster_update.md)	 - updates the cluster

