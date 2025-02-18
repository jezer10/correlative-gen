package models

import (
	"encoding/json"
)

type RequestBody struct {
	CheckIDs []json.RawMessage `json:"check_ids"`
}
