package error

import "encoding/json"

type ServiceError struct {
	Type   string      `json:"type"`
	Errors interface{} `json:"error"`
}

func (e ServiceError) Error() string {
	return e.Type
}

func (e ServiceError) MarshalJSON() ([]byte, error) {
	type Alias ServiceError
	return json.Marshal(struct {
		Error string `json:"error"`
		*Alias
	}{
		Error: e.Type,
		Alias: (*Alias)(&e),
	})
}
