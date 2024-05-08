package unmarshal

import (
	"encoding/json"
	"strings"
)

type TrimmedString string

func (self *TrimmedString) UnmarshalJSON(data []byte) error {
	var result string

	err := json.Unmarshal(data, &result)

	if err != nil {
		return err
	}

	*self = TrimmedString(strings.TrimSpace(result))

	return nil
}
