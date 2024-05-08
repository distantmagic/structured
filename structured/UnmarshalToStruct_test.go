package structured

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalsToStruct(t *testing.T) {
	result := EntityExtractorResult{
		Content: "{\"name\":\"John\",\"surname\":\"Doe\",\"age\":40}",
	}

	var person fixtureNamedPersonWithAge

	err := UnmarshalToStruct(result, &person)

	assert.Nil(t, err)
	assert.Equal(t, "John", person.Name)
	assert.Equal(t, "Doe", person.Surname)
	assert.Equal(t, 40, person.Age)
}
