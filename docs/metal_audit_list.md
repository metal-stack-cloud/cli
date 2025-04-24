## metal audit list

list all audit traces

```
metal audit list [flags]
```

### Options

```
      --body string         filters audit trace body payloads for the given text (full-text search).
      --from string         start of range of the audit traces. e.g. 1h, 10m, 2006-01-02 15:04:05
  -h, --help                help for list
      --limit int           limit the number of audit traces.
      --method string       api method of the audit trace.
      --phase string        the audit trace phase.
      --prettify-body       attempts to interpret the body as json and prettifies it.
      --project string      project id of the audit trace
      --request-id string   request id of the audit trace.
      --result-code int32   gRPC result status code of the audit trace.
      --sort-by strings     sort by (comma separated) column(s), sort direction can be changed by appending :asc or :desc behind the column identifier. possible values: method|project|timestamp|user
      --source-ip string    source-ip of the audit trace.
      --tenant string       tenant of the audit trace.
      --to string           end of range of the audit traces. e.g. 1h, 10m, 2006-01-02 15:04:05
      --user string         user of the audit trace.
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

* [metal audit](metal_audit.md)	 - manage audit trace entities

