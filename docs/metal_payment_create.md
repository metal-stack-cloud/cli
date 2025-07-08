## metal payment create

creates the payment

```
metal payment create [flags]
```

### Options

```
      --bulk-output                  when used with --file (bulk operation): prints results at the end as a list. default is printing results intermediately during the operation, which causes single entities to be printed in a row.
      --email string                 the email of the tenant
  -f, --file string                  filename of the create or update request in yaml format, or - for stdin.
                                     
                                     Example:
                                     $ metal payment describe payment-1 -o yaml > payment.yaml
                                     $ vi payment.yaml
                                     $ # either via stdin
                                     $ cat payment.yaml | metal payment create -f -
                                     $ # or via file
                                     $ metal payment create -f payment.yaml
                                     
                                     the file can also contain multiple documents and perform a bulk operation.
                                     	
  -h, --help                         help for create
      --name string                  the name of the tenant
      --phone string                 the phone number of the tenant
      --skip-security-prompts        skips security prompt for bulk operations
      --stripe-public-token string   the stripe public token (default "pk_live_51LyJ2zKXUtoWqdO7FyjC4mrAsVZ4tatWY4mCv8tMRd8FgG1Wzn3RxGrwPQdE9Ic4qQMp5bcPiPBHoxJwWUXpcPz700f1NQLmNa")
      --timestamps                   when used with --file (bulk operation): prints timestamps in-between the operations
      --vat string                   the vat of the tenant
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

* [metal payment](metal_payment.md)	 - manage payment entities

