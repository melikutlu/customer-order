package model

type QueryParams struct {
	FirstName      string `json:"first_name"`
	AgeGreaterThan string `json:"agt"`
	AgeLessThan    string `json:"alt"`
	Limit          string `json:"limit"`
	Offset         string `json:"offset"`
}
