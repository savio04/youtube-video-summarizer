package entities

type Video struct {
	Id         string `json:"id"`
	ExternalId string `json:"external_id"`
	Url        string `json:"url"`
	Summary    string `json:"Summary"`
}
