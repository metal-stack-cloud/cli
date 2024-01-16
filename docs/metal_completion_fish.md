## metal completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	metal completion fish | source

To load completions for every new session, execute once:

	metal completion fish > ~/.config/fish/completions/metal.fish

You will need to start a new shell for this setup to take effect.


```
metal completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
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

* [metal completion](metal_completion.md)	 - Generate the autocompletion script for the specified shell

