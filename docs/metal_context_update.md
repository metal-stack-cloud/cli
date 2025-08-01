## metal context update

update a cli context

```
metal context update <context-name> [flags]
```

### Options

```
      --activate                 immediately switches to the new context
      --api-token string         sets the api-token for this context
      --api-url string           sets the api-url for this context
      --default-project string   sets a default project to act on
  -h, --help                     help for update
      --provider string          sets the login provider for this context
      --timeout duration         sets a default request timeout
```

### Options inherited from parent commands

```
  -c, --config string          alternative config file path, (default is ~/.metal-stack-cloud/config.yaml)
      --debug                  debug output
      --force-color            force colored output even without tty
  -o, --output-format string   output format (table|wide|markdown|json|yaml|template|jsonraw|yamlraw), wide is a table with more columns, jsonraw and yamlraw do not translate proto enums into string types but leave the original int32 values intact. (default "table")
      --template string        output template for template output-format, go template format. For property names inspect the output of -o json or -o yaml for reference.
```

### SEE ALSO

* [metal context](metal_context.md)	 - manage cli contexts

