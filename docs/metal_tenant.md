## metal tenant

manage tenant entities

### Synopsis

manage api tenants

### Options

```
  -h, --help   help for tenant
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
* [metal tenant apply](metal_tenant_apply.md)	 - applies one or more tenants from a given file
* [metal tenant create](metal_tenant_create.md)	 - creates the tenant
* [metal tenant delete](metal_tenant_delete.md)	 - deletes the tenant
* [metal tenant describe](metal_tenant_describe.md)	 - describes the tenant
* [metal tenant edit](metal_tenant_edit.md)	 - edit the tenant through an editor and update
* [metal tenant invite](metal_tenant_invite.md)	 - manage tenant invites
* [metal tenant join](metal_tenant_join.md)	 - join a tenant of someone who shared an invite secret with you
* [metal tenant list](metal_tenant_list.md)	 - list all tenants
* [metal tenant member](metal_tenant_member.md)	 - manage tenant members
* [metal tenant update](metal_tenant_update.md)	 - updates the tenant

