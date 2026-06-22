package channel

type Channel struct {
	Id                  int64  `json:"id"`
	UserId              int64  `json:"userId"`
	ChannelTemplateId   int64  `json:"channelTemplateId"`
	Name                string `json:"name"`
	ChannelType         string `json:"channelType"`
	ProviderCode        string `json:"providerCode"`
	SupportedEventTypes string `json:"supportedEventTypes"`
	Icon                string `json:"icon"`
	Sort                int    `json:"sort"`
	IsActive            bool   `json:"isActive"`
}

type ChannelQuery struct {
	ChannelType string
	EventType   string
	IsActive    *bool
}

type ImportRequest struct {
	TemplateIds []int64 `json:"templateIds"`
}
