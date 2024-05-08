package structured

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type myTestPerson struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     int    `json:"age"`
}

func TestUnmarshalsToStruct(t *testing.T) {
	result := EntityExtractorResult{
		Content: "{\"name\":\"John\",\"surname\":\"Doe\",\"age\":40}",
	}

	var person myTestPerson

	err := UnmarshalToStruct(result, &person)

	assert.Nil(t, err)
	assert.Equal(t, "John", person.Name)
	assert.Equal(t, "Doe", person.Surname)
	assert.Equal(t, 40, person.Age)
}
