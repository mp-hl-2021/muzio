package common

type Link struct {
	ServiceName string `json:"serviceName"`
	Url         string `json:"url"`
	IsAvailable bool   `json:"isAvailable"`
}
