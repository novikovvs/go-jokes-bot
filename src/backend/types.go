package backend

import "time"

type GetKeyResponse struct {
	AuthKey string    `json:"auth_key"`
	Expire  time.Time `json:"expire"`
}
type R struct {
	Ok     bool `json:"ok"`
	Result struct {
		GetKeyResponse
	} `json:"result"`
	Error interface{} `json:"error"`
}
