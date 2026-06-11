package category

type Category struct {
	Id                 int64  `json:"id"`
	UserId             int64  `json:"userId"`
	Name               string `json:"name"`
	Type               string `json:"type"`
	CategoryGroupId    int64  `json:"categoryGroupId"`
	CategoryGroupCode  string `json:"categoryGroupCode"`
	CategoryGroupName  string `json:"categoryGroupName"`
	CategoryGroupColor string `json:"categoryGroupColor"`
	CategoryGroupIcon  string `json:"categoryGroupIcon"`
	IsActive           bool   `json:"isActive"`
}
