package csvtojson

type config struct {
	separator rune
	fields    map[string]string
	inject    map[string]string
}

func newConfig(options ...option) *config {
	cfg := &config{
		separator: ',',
	}

	for _, opt := range options {
		opt(cfg)
	}

	return cfg
}

type option func(c *config)

// WithSeparator specifies the csv separator
// defaults to ','
func WithSeparator(separator rune) option {
	return func(c *config) {
		c.separator = separator
	}
}

// WithFields sets the field mapping. If not provided, csv headers will be used as output keys
// given fields:
//
//	```go
//	map[string]string{"field1": "id", "field2": "name"}
//	```
//
// and csv input as:
//
//	```csv
//	field1,field2
//	123,john doe
//	456,jane doe
//	```
//
// results in json output:
//
//	```json
//	{"id": "123", "name": "john doe"}
//	{"id": "456", "name": "jane doe"}
//	```
//
// Note: if any fields(s) are missing in the incoming csv headers, Copy will fail with an error specifying the missing headers
func WithFields(fields map[string]string) option {
	return func(c *config) {
		c.fields = fields
	}
}

// WithInject specifies static key/value mapping to inject into the output json objects
// given inject:
//
//	```go
//	map[string]string{"mycustomkey1": "abc123", "mycustomkey2": "20230118"}
//	```
//
// and csv input as:
//
//	```csv
//	id,name
//	123,john doe
//	456,jane doe
//	```
//
// results in json output:
//
//	```json
//	{"id": "123", "name": "john doe", "mycustomkey1": "abc123", "mycustomkey2": "20230118"}
//	{"id": "456", "name": "jane doe", "mycustomkey1": "abc123", "mycustomkey2": "20230118"}
//	```
//
// Note: if any key(s) clash with an incoming key, Copy will fail with an error specifying the clashing keys
//
// Note: inject does respect any specified custom field mapping from config.fields
func WithInject(inject map[string]string) option {
	return func(c *config) {
		c.inject = inject
	}
}
