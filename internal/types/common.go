package types

import (
	"encoding/json"
	"strconv"
)

type IDArray []int64

func (m IDArray) MarshalJSON() ([]byte, error) {
	stringArray := make([]string, len(m))
	for i, v := range m {
		stringArray[i] = strconv.FormatInt(v, 10)
	}
	return json.Marshal(stringArray)
}

func (m *IDArray) UnmarshalJSON(data []byte) error {
	var stringArray []string
	if err := json.Unmarshal(data, &stringArray); err != nil {
		return err
	}

	*m = make(IDArray, len(stringArray))
	for i, s := range stringArray {
		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		(*m)[i] = val
	}
	return nil
}
