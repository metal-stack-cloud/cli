## metal tenant invite

manage tenant invites

### Options

```
  -h, --help   help for invite
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
* [metal tenant invite delete](metal_tenant_invite_delete.md)	 - deletes a pending invite
* [metal tenant invite generate-join-secret](metal_tenant_invite_generate-join-secret.md)	 - generate an invite secret to share with the new member
* [metal tenant invite join](metal_tenant_invite_join.md)	 - join a tenant of someone who shared an invite secret with you
* [metal tenant invite list](metal_tenant_invite_list.md)	 - lists the currently pending invites

