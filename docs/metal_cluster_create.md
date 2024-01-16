## metal cluster create

creates the cluster

```
metal cluster create [flags]
```

### Options

```
      --bulk-output                     when used with --file (bulk operation): prints results at the end as a list. default is printing results intermediately during the operation, which causes single entities to be printed in a row.
  -f, --file string                     filename of the create or update request in yaml format, or - for stdin.
                                        
                                        Example:
                                        $ metal cluster describe cluster-1 -o yaml > cluster.yaml
                                        $ vi cluster.yaml
                                        $ # either via stdin
                                        $ cat cluster.yaml | metal cluster create -f -
                                        $ # or via file
                                        $ metal cluster create -f cluster.yaml
                                        
                                        the file can also contain multiple documents and perform a bulk operation.
                                        	
  -h, --help                            help for create
      --kubernetes-version string       kubernetes version of the cluster
      --maintenance-duration duration   duration in which cluster maintenance is allowed to take place (default 2h0m0s)
      --maintenance-hour int32          hour in which cluster maintenance is allowed to take place
      --maintenance-minute int32        minute in which cluster maintenance is allowed to take place
      --maintenance-timezone string     timezone used for the maintenance time window (default "Local")
      --name string                     name of the cluster
      --partition string                partition of the cluster
  -p, --project string                  project of the cluster
      --skip-security-prompts           skips security prompt for bulk operations
      --timestamps                      when used with --file (bulk operation): prints timestamps in-between the operations
      --worker-group string             the name of the initial worker group (default "group-0")
      --worker-max uint32               the maximum amount of worker nodes of the worker group (default 3)
      --worker-max-surge uint32         the maximum amount of new worker nodes added to the worker group during a rolling update (default 1)
      --worker-max-unavailable uint32   the maximum amount of worker nodes removed from the worker group during a rolling update
      --worker-min uint32               the minimum amount of worker nodes of the worker group (default 1)
      --worker-type string              the worker type of the initial worker group
```

### Options inherited from parent commands

```
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

