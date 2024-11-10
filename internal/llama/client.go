package llama

type ProviderConfig struct {
	baseURL string
	model   string
}

func NewOlamaProviderConfig(baseURL, model string) ProviderConfig {
	return ProviderConfig{
		baseURL: baseURL,
		model:   model,
	}
}

func (o ProviderConfig) BaseURL() string {
	return o.baseURL
}

func (o ProviderConfig) WithBaseURL(baseURL string) ProviderConfig {
	o.baseURL = baseURL
	return o
}

func (o ProviderConfig) Model() string {
	return o.model
}

func (o ProviderConfig) WithModel(model string) ProviderConfig {
	o.model = model
	return o
}
