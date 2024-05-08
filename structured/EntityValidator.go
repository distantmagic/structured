package structured

import (
	"github.com/santhosh-tekuri/jsonschema/v5"
)

type EntityValidator struct {
	CompiledJsonSchema  *jsonschema.Schema
	MarshaledJsonSchema []byte
}
