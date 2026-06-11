package categorygroup

type CategoryGroup struct {
	Id        int64  `json:"id"`
	GroupCode string `json:"groupCode"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Icon      string `json:"icon"`
	Color     string `json:"color"`
	Sort      int    `json:"sort"`
	Enabled   bool   `json:"enabled"`
}
