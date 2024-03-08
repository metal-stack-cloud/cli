## metal project

manage project entities

### Synopsis

manage api projects

### Options

```
  -h, --help   help for project
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
* [metal project apply](metal_project_apply.md)	 - applies one or more projects from a given file
* [metal project create](metal_project_create.md)	 - creates the project
* [metal project delete](metal_project_delete.md)	 - deletes the project
* [metal project describe](metal_project_describe.md)	 - describes the project
* [metal project edit](metal_project_edit.md)	 - edit the project through an editor and update
* [metal project invite](metal_project_invite.md)	 - manage project invites
* [metal project join](metal_project_join.md)	 - join a project of someone who shared an invite secret with you
* [metal project list](metal_project_list.md)	 - list all projects
* [metal project remove-member](metal_project_remove-member.md)	 - remove member from a project
* [metal project update](metal_project_update.md)	 - updates the project
* [metal project update-member](metal_project_update-member.md)	 - update member from a project

