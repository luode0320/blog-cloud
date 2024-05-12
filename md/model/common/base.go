package common

type PageCondition[T interface{}] struct {
	Page      Page `json:"page"`
	Condition T    `json:"condition"`
}

type PageResult[T interface{}] struct {
	Records []T `json:"records"`
	Total   int `json:"total"`
}

type Page struct {
	Current int `json:"current"`
	Size    int `json:"size"`
}

type CountResult struct {
	Count int `json:"count" db:"count"`
}

type TokenResult struct {
	Name         string `json:"name"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenCache struct {
	Id string `json:"id"`
	TokenResult
}
