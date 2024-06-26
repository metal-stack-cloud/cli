## metal tenant member

manage tenant members

### Options

```
  -h, --help   help for member
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

* [metal tenant](metal_tenant.md)	 - manage tenant entities
* [metal tenant member list](metal_tenant_member_list.md)	 - lists members of a tenant
* [metal tenant member remove](metal_tenant_member_remove.md)	 - remove member from a tenant
* [metal tenant member update](metal_tenant_member_update.md)	 - update member from a tenant

