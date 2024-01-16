## metal completion

Generate the autocompletion script for the specified shell

### Synopsis

Generate the autocompletion script for metal for the specified shell.
See each sub-command's help for details on how to use the generated script.


### Options

```
  -h, --help   help for completion
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
* [metal completion bash](metal_completion_bash.md)	 - Generate the autocompletion script for bash
* [metal completion fish](metal_completion_fish.md)	 - Generate the autocompletion script for fish
* [metal completion powershell](metal_completion_powershell.md)	 - Generate the autocompletion script for powershell
* [metal completion zsh](metal_completion_zsh.md)	 - Generate the autocompletion script for zsh

