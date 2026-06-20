package categorytemplate

type CategoryTemplate struct {
	Id           int64  `json:"id"`
	TemplateCode string `json:"templateCode"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Icon         string `json:"icon"`
	Color        string `json:"color"`
	Sort         int    `json:"sort"`
	Enabled      bool   `json:"enabled"`
}

type CreateTemplateRequest struct {
	TemplateCode string `json:"templateCode"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Icon         string `json:"icon"`
	Color        string `json:"color"`
	Sort         int    `json:"sort"`
}

type UpdateTemplateRequest struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Icon    string `json:"icon"`
	Color   string `json:"color"`
	Sort    int    `json:"sort"`
	Enabled *bool  `json:"enabled"`
}

type ImportRequest struct {
	TemplateIds []int64 `json:"templateIds"`
}
