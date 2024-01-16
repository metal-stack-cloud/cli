## metal cluster kubeconfig

fetch kubeconfig of a cluster

```
metal cluster kubeconfig [flags]
```

### Options

```
      --expiration duration   kubeconfig will expire after given time (default 8h0m0s)
  -h, --help                  help for kubeconfig
      --kubeconfig string     specify an explicit path for the merged kubeconfig to be written, defaults to default kubeconfig paths if not provided
      --merge                 merges the kubeconfig into default kubeconfig instead of printing it to the console (default true)
  -p, --project string        the project in which the cluster resides for which to get the kubeconfig for
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

* [metal cluster](metal_cluster.md)	 - manage cluster entities

