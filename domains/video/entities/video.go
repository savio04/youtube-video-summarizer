package entities

type Video struct {
	Id         *int    `json:"id"`
	Status     *string `json:"status"`
	ExternalId *string `json:"external_id"`
	Url        string  `json:"url"`
	Summary    *string `json:"summary"`
}
