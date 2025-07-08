## metal payment

manage payment entities

### Synopsis

manage payment of the metalstack.cloud

### Options

```
  -h, --help   help for payment
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

* [metal](metal.md)	 - cli for managing entities in metal-stack-cloud
* [metal payment apply](metal_payment_apply.md)	 - applies one or more payments from a given file
* [metal payment create](metal_payment_create.md)	 - creates the payment
* [metal payment delete](metal_payment_delete.md)	 - deletes the payment
* [metal payment describe](metal_payment_describe.md)	 - describes the payment
* [metal payment edit](metal_payment_edit.md)	 - edit the payment through an editor and update
* [metal payment show-default-prices](metal_payment_show-default-prices.md)	 - show default prices
* [metal payment update](metal_payment_update.md)	 - updates the payment

