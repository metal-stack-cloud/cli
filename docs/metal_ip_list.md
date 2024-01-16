## metal ip list

list all ips

```
metal ip list [flags]
```

### Options

```
  -h, --help              help for list
  -p, --project string    project from where ips should be listed
      --sort-by strings   sort by (comma separated) column(s), sort direction can be changed by appending :asc or :desc behind the column identifier. possible values: ip|name|network|project|type|uuid
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

* [metal ip](metal_ip.md)	 - manage ip entities

