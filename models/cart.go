package models

type CancelCartBody struct {
	T       int    `json:"t"`
	OutSkus string `json:"outSkus"`
	Random  int    `json:"random"`
}
