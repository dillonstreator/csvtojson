package csvtojson

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	var buf bytes.Buffer

	r := strings.NewReader(`
id,name
123,john doe
456,jane doe
`)

	err := Copy(ctx, r, &buf)
	assert.NoError(err)

	expected := `{"id":"123","name":"john doe"}
{"id":"456","name":"jane doe"}
`

	assert.Equal(expected, buf.String())
}

func TestCopy_WithSeparator(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	var buf bytes.Buffer

	r := strings.NewReader(`
id|name
123|john doe
456|jane doe
`)

	err := Copy(ctx, r, &buf, WithSeparator('|'))
	assert.NoError(err)

	expected := `{"id":"123","name":"john doe"}
{"id":"456","name":"jane doe"}
`

	assert.Equal(expected, buf.String())
}

func TestCopy_fields(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	var buf bytes.Buffer

	r := strings.NewReader(`
field1,field2,other
123,john doe,1
456,jane doe,2
`)

	err := Copy(ctx, r, &buf, WithFields(map[string]string{"field1": "id", "field2": "name"}))
	assert.NoError(err)

	expected := `{"id":"123","name":"john doe"}
{"id":"456","name":"jane doe"}
`

	assert.Equal(expected, buf.String())
}

func TestCopy_fields_empty(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	var buf bytes.Buffer

	r := strings.NewReader(`
field1,field2,other
123,john doe,1
456,jane doe,2
`)

	err := Copy(ctx, r, &buf, WithFields(map[string]string{}))
	assert.NoError(err)

	expected := `{}
{}
`

	assert.Equal(expected, buf.String())
}

func TestCopy_fields_missing_header(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	var buf bytes.Buffer

	r := strings.NewReader(`
field1,field3,other
123,john doe,1
456,jane doe,2
`)

	err := Copy(ctx, r, &buf, WithFields(map[string]string{"field1": "id", "field2": "name"}))
	assert.ErrorContains(err, "missing header(s) in incoming csv: field2")
}

func TestCopy_inject(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	var buf bytes.Buffer

	r := strings.NewReader(`
id,name
123,john doe
456,jane doe
`)

	err := Copy(ctx, r, &buf, WithInject(map[string]string{"mycustomkey": "mycustomvalue"}))
	assert.NoError(err)

	expected := `{"id":"123","mycustomkey":"mycustomvalue","name":"john doe"}
{"id":"456","mycustomkey":"mycustomvalue","name":"jane doe"}
`

	assert.Equal(expected, buf.String())
}

func TestCopy_inject_clashing_key(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	var buf bytes.Buffer

	r := strings.NewReader(`
id,name,other
123,john doe,1
456,jane doe,2
`)

	err := Copy(ctx, r, &buf, WithInject(map[string]string{"other": "mycustomvalue"}))
	assert.ErrorContains(err, "clashing key(s) in injection: other")
}

func TestCopy_inject_clashing_mapped_key(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	var buf bytes.Buffer

	r := strings.NewReader(`
id,name,other
123,john doe,1
456,jane doe,2
`)

	err := Copy(ctx, r, &buf, WithFields(map[string]string{"other": "someother"}), WithInject(map[string]string{"someother": "mycustomvalue"}))
	assert.ErrorContains(err, "clashing key(s) in injection: someother")
}
