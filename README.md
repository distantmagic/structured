# Instructor Go (work in progress)

Go converstion of https://github.com/jxnl/instructor/

Same features, Go-like API. Model agnostic - maps data from arbitrary JSON
schema to arbitrary Go struct.

Focused on llama.cpp.

## Usage

API can probably change with time until all features are implemented.

### Initializing the Mapper

```go
import (
	"fmt"
	"net/http"
	"testing"

	"github.com/distantmagic/paddler/llamacpp"
	"github.com/distantmagic/paddler/netcfg"
)

var jsonSchemaMapper *JsonSchemaMapper = &JsonSchemaMapper{
	LlamaCppClient: &llamacpp.LlamaCppClient{
		HttpClient: http.DefaultClient,
		LlamaCppConfiguration: &llamacpp.LlamaCppConfiguration{
			HttpAddress: &netcfg.HttpAddressConfiguration{
				Host:   "127.0.0.1",
				Port:   8081,
				Scheme: "http",
			},
		},
	},
}
```

### Extracting Structured Data from String

```go
responseChannel := make(chan JsonSchemaMapperResult)

go jsonSchemaMapper.MapToSchema(
	responseChannel,
	map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]string{
				"type": "string",
			},
			"surname": map[string]string{
				"type": "string",
			},
			"age": map[string]string{
				"description": "Age in years.",
				"type":        "integer",
			},
		},
	},
	"I am John Doe - living for 40 years and I still like to play chess.",
)

for result := range responseChannel {
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	// map[name:John, surname:Doe, age:40]
	fmt.Print(result.Result)
}
```