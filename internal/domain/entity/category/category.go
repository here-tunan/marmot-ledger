package category

type Category struct {
	Id       int64  `json:"id"`
	UserId   int64  `json:"userId"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Icon     string `json:"icon"`
	Color    string `json:"color"`
	IsActive bool   `json:"isActive"`
	// 所属家庭分组ID列表（返回时填充，不存储）
	GroupIds []int64 `json:"groupIds"`
}
