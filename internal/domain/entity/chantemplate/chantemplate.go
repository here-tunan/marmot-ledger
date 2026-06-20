package chantemplate

type ChannelTemplate struct {
	Id                  int64  `json:"id"`
	ChannelCode         string `json:"channelCode"`
	Name                string `json:"name"`
	ChannelType         string `json:"channelType"`
	ProviderCode        string `json:"providerCode"`
	SupportedEventTypes string `json:"supportedEventTypes"`
	Icon                string `json:"icon"`
	Sort                int    `json:"sort"`
	Enabled             bool   `json:"enabled"`
}

type CreateTemplateRequest struct {
	ChannelCode         string `json:"channelCode"`
	Name                string `json:"name"`
	ChannelType         string `json:"channelType"`
	ProviderCode        string `json:"providerCode"`
	SupportedEventTypes string `json:"supportedEventTypes"`
	Icon                string `json:"icon"`
	Sort                int    `json:"sort"`
}

type UpdateTemplateRequest struct {
	Name                string `json:"name"`
	ChannelType         string `json:"channelType"`
	ProviderCode        string `json:"providerCode"`
	SupportedEventTypes string `json:"supportedEventTypes"`
	Icon                string `json:"icon"`
	Sort                int    `json:"sort"`
	Enabled             *bool  `json:"enabled"`
}
