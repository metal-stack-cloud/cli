## metal

cli for managing entities in metal-stack-cloud

### Options

```
      --api-token string       the token used for api requests
      --api-url string         the url to the metalstack.cloud api (default "https://api.metalstack.cloud")
  -c, --config string          alternative config file path, (default is ~/.metal-stack-cloud/config.yaml)
      --debug                  debug output
      --force-color            force colored output even without tty
  -h, --help                   help for metal
  -o, --output-format string   output format (table|wide|markdown|json|yaml|template|jsonraw|yamlraw), wide is a table with more columns, jsonraw and yamlraw do not translate proto enums into string types but leave the original int32 values intact. (default "table")
      --template string        output template for template output-format, go template format. For property names inspect the output of -o json or -o yaml for reference.
      --timeout duration       request timeout used for api requests
```

### SEE ALSO

* [metal api-methods](metal_api-methods.md)	 - show available api-methods of the metalstack.cloud api
* [metal asset](metal_asset.md)	 - show asset
* [metal cluster](metal_cluster.md)	 - manage cluster entities
* [metal completion](metal_completion.md)	 - Generate the autocompletion script for the specified shell
* [metal context](metal_context.md)	 - manage cli contexts
* [metal health](metal_health.md)	 - print the client and server health information
* [metal ip](metal_ip.md)	 - manage ip entities
* [metal markdown](metal_markdown.md)	 - create markdown documentation
* [metal payment](metal_payment.md)	 - manage payment of the metalstack.cloud
* [metal project](metal_project.md)	 - manage project entities
* [metal storage](metal_storage.md)	 - storage commands
* [metal tenant](metal_tenant.md)	 - manage tenant entities
* [metal token](metal_token.md)	 - manage token entities
* [metal user](metal_user.md)	 - manage user entities
* [metal version](metal_version.md)	 - print the client and server version information

