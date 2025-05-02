## metal storage volume

manage volume entities

### Synopsis

volume related actions of metalstack.cloud

### Options

```
  -h, --help   help for volume
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

* [metal storage](metal_storage.md)	 - storage commands
* [metal storage volume delete](metal_storage_volume_delete.md)	 - deletes the volume
* [metal storage volume describe](metal_storage_volume_describe.md)	 - describes the volume
* [metal storage volume encryptionsecret](metal_storage_volume_encryptionsecret.md)	 - volume encryptionsecret template
* [metal storage volume list](metal_storage_volume_list.md)	 - list all volumes
* [metal storage volume manifest](metal_storage_volume_manifest.md)	 - volume manifest
* [metal storage volume update](metal_storage_volume_update.md)	 - updates the volume

