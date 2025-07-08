## metal context

manage cli contexts

### Synopsis

you can switch back and forth contexts with "-"

```
metal context [flags]
```

### Options

```
  -h, --help   help for context
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
* [metal context add](metal_context_add.md)	 - add a cli context
* [metal context list](metal_context_list.md)	 - list the configured cli contexts
* [metal context remove](metal_context_remove.md)	 - remove a cli context
* [metal context set-project](metal_context_set-project.md)	 - sets the default project to act on for cli commands
* [metal context show-current](metal_context_show-current.md)	 - prints the current context name
* [metal context switch](metal_context_switch.md)	 - switch the cli context
* [metal context update](metal_context_update.md)	 - update a cli context

