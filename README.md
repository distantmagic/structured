# Structured (work in progress)

The project started as a Go conversion of https://github.com/jxnl/instructor/,
but it is a more general-purpose library.

It also features a language-agnostic HTTP server that you can set up in front
of [llama.cpp](https://github.com/ggerganov/llama.cpp).

Same features, Go-like API. Model agnostic - maps data from arbitrary JSON
schema to arbitrary Go struct (or just plain JSON).

It is focused on [llama.cpp](https://github.com/ggerganov/llama.cpp). Support
for other vendor APIs (like OpenAI or Anthropic) might be added in the future.

## HTTP API

```mermaid
sequenceDiagram
    You->>Structured: JSON schema + data
    Structured->>llama.cpp: extract
    llama.cpp->>Structured: extracted entity
    Structured->>Structured: validates extracted entity (double check)
    Structured-->>llama.cpp: retry if validation fails
    Structured->>You: JSON matching your schema schema
```

Start a server and point it to your local
[llama.cpp](https://github.com/ggerganov/llama.cpp) instance:

```shell
./structured \
	--llamacpp-host 127.0.0.1 \
	--llamacpp-port 8081 \
	--port 8080
```

Structured server connects to
[llama.cpp](https://github.com/ggerganov/llama.cpp) to extract the data.

Now, you can issue requests. Include `schema` and `data` in your POST body.
The server will respond with JSON matching your schema:

```
POST http://127.0.0.1:8080/extract/entity
{
  "schema": {
    "type": "object",
    "properties": {
      "hello": {
        "type": "string"
      }
    },
    "required": ["hello"]
  },
  "data": "Say 'world'"
}

Response:
{
  "hello": "world"
}
```

## Programmatic Usage

API can change with time until all features are implemented.

### Initializing the Mapper

Point it to your local [llama.cpp](https://github.com/ggerganov/llama.cpp)
instance:

```go
import (
	"fmt"
	"net/http"
	"testing"

	"github.com/distantmagic/structured/structured"
	"github.com/distantmagic/paddler/llamacpp"
	"github.com/distantmagic/paddler/netcfg"
)

var entityExtractor *EntityExtractor = &structured.EntityExtractor{
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
	MaxRetries: 3,
}
```

### Extracting Structured Data from String

```go
import "github.com/distantmagic/structured/structured"

responseChannel := make(chan structured.EntityExtractorResult)

go entityExtractor.ExtractFromString(
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
		panic(result.Error)
	}

	// map[name:John, surname:Doe, age:40]
	fmt.Print(result.Result)
}
```

### Mapping Extracted Result onto an Arbitrary Struct

Once you obtain the result:

```go
import "github.com/distantmagic/structured/structured"

type myTestPerson struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     int    `json:"age"`
}

func DoUnmarshalsToStruct(result structured.EntityExtractorResult) {
	var person myTestPerson

	err := structured.UnmarshalToStruct(result, &person)

	if nil != err {
		panic(err)
	}

	person.Name // John
	person.Surname // Doe
}
```

## See Also

[Paddler](https://github.com/distantmagic/paddler) - (work in progress)
	[llama.cpp](https://github.com/ggerganov/llama.cpp) load balancer,
	supervisor and request queue

## Community

- [Discord](https://discord.gg/kysUzFqSCK)
