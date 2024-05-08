package structured

import (
	"github.com/santhosh-tekuri/jsonschema/v5"
)

type EntityJsonSchemaValidator struct {
	CompiledJsonSchema  *jsonschema.Schema
	MarshaledJsonSchema []byte
}
