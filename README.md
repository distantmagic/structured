# Instructor Go (work in progress)

Go conversion of https://github.com/jxnl/instructor/

Same features, Go-like API. Model agnostic - maps data from arbitrary JSON
schema to arbitrary Go struct.

It is focused on [llama.cpp](https://github.com/ggerganov/llama.cpp). Support
for other vendor APIs (like OpenAI or Anthropic) might be added in the future.

## Usage

API can change with time until all features are implemented.

### Initializing the Mapper

Point it to your local [llama.cpp](https://github.com/ggerganov/llama.cpp)
instance:

```go
import (
	"fmt"
	"net/http"
	"testing"

	"github.com/distantmagic/instructor-go/instructor"
	"github.com/distantmagic/paddler/llamacpp"
	"github.com/distantmagic/paddler/netcfg"
)

var entityExtractor *EntityExtractor = &instructor.EntityExtractor{
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
import "github.com/distantmagic/instructor-go/instructor"

responseChannel := make(chan instructor.EntityExtractorResult)

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
		t.Fatal(result.Error)
	}

	// map[name:John, surname:Doe, age:40]
	fmt.Print(result.Result)
}
```

### Mapping Extracted Result onto an Arbitrary Struct

Once you obtain the result:

```go
import "github.com/distantmagic/instructor-go/instructor"

type myTestPerson struct {
	Name   string `json:"name"`
	Surname string `json:"surname"`
	Age    int    `json:"age"`
}

func DoUnmarshalsToStruct(result EntityExtractorResult) {
	var person myTestPerson

	err := instructor.UnmarshalToStruct(result, &person)

	if nil != err {
		panic(err)
	}

	person.Name // John
	person.Surname // Doe
}
```