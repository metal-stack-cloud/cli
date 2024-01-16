## metal storage snapshot

manage snapshot entities

### Synopsis

snapshot related actions of metalstack.cloud

### Options

```
  -h, --help   help for snapshot
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

* [metal storage](metal_storage.md)	 - storage commands
* [metal storage snapshot delete](metal_storage_snapshot_delete.md)	 - deletes the snapshot
* [metal storage snapshot describe](metal_storage_snapshot_describe.md)	 - describes the snapshot
* [metal storage snapshot list](metal_storage_snapshot_list.md)	 - list all snapshots

