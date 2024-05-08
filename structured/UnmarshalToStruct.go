package structured

import "encoding/json"

func UnmarshalToStruct[T any](
	result EntityExtractorResult,
	out *T,
) error {
	if result.Error != nil {
		return result.Error
	}

	err := json.Unmarshal([]byte(result.Content), out)

	if err != nil {
		return err
	}

	return nil
}
