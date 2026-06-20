package account

type Account struct {
	Id       int64  `json:"id"`
	UserId   int64  `json:"userId"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Icon     string `json:"icon"`
	Color    string `json:"color"`
	IsActive bool   `json:"isActive"`
}
