package accounttemplate

type AccountTemplate struct {
	Id           int64  `json:"id"`
	ProviderCode string `json:"providerCode"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Icon         string `json:"icon"`
	Color        string `json:"color"`
	Sort         int    `json:"sort"`
	Enabled      bool   `json:"enabled"`
}

type CreateTemplateRequest struct {
	ProviderCode string `json:"providerCode"`
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

type InstantiateRequest struct {
	TemplateId int64  `json:"templateId"`
	CustomName string `json:"customName"`
}
