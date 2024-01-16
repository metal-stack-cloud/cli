## metal ip

manage ip entities

### Synopsis

an ip address of metalstack.cloud

### Options

```
  -h, --help   help for ip
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
* [metal ip apply](metal_ip_apply.md)	 - applies one or more ips from a given file
* [metal ip create](metal_ip_create.md)	 - creates the ip
* [metal ip delete](metal_ip_delete.md)	 - deletes the ip
* [metal ip describe](metal_ip_describe.md)	 - describes the ip
* [metal ip edit](metal_ip_edit.md)	 - edit the ip through an editor and update
* [metal ip list](metal_ip_list.md)	 - list all ips
* [metal ip update](metal_ip_update.md)	 - updates the ip

