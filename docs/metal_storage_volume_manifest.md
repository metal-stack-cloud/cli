## metal storage volume manifest

volume manifest

### Synopsis

show detailed info about the storage cluster

```
metal storage volume manifest [flags]
```

### Options

```
  -h, --help               help for manifest
      --name string        name of the PersistentVolume (default "restored-pv")
      --namespace string   namespace for the PersistentVolume (default "default")
  -p, --project string     project
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

* [metal storage volume](metal_storage_volume.md)	 - manage volume entities

