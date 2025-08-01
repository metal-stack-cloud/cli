## metal tenant update

updates the tenant

```
metal tenant update <id> [flags]
```

### Options

```
      --accept-terms-and-conditions   can be used to accept the terms and conditions
      --avatar-url string             the avatar url of the tenant to create
      --bulk-output                   when used with --file (bulk operation): prints results at the end as a list. default is printing results intermediately during the operation, which causes single entities to be printed in a row.
      --description string            the description of the tenant to update
      --email string                  the name of the tenant to update
  -f, --file string                   filename of the create or update request in yaml format, or - for stdin.
                                      
                                      Example:
                                      $ metal tenant describe tenant-1 -o yaml > tenant.yaml
                                      $ vi tenant.yaml
                                      $ # either via stdin
                                      $ cat tenant.yaml | metal tenant update <id> -f -
                                      $ # or via file
                                      $ metal tenant update <id> -f tenant.yaml
                                      
                                      the file can also contain multiple documents and perform a bulk operation.
                                      	
  -h, --help                          help for update
      --name string                   the name of the tenant to update
      --skip-security-prompts         skips security prompt for bulk operations
      --tenant string                 the tenant to update
      --timestamps                    when used with --file (bulk operation): prints timestamps in-between the operations
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

