package ai_struct

type AiConfig struct {
	OpenAiConfig
}

type OpenAiConfig struct {
	EnableOpenAI       bool     `json:"enable_openai" yaml:"enable_openai"`
	OpenAIAzure        bool     `json:"openai_azure" yaml:"openai_azure"`
	OpenAIAPIVersion   string   `json:"openai_api_version" yaml:"openai_api_version"`
	OpenAIAutoProxy    bool     `json:"openai_auto_proxy" yaml:"openai_auto_proxy"`
	OpenAIOrganization string   `json:"openai_organization" yaml:"openai_organization"`
	OpenAIServers      []string `json:"openai_servers" yaml:"openai_servers"`
	OpenAIKeys         []string `json:"openai_keys" yaml:"openai_keys"`

	EnableOpenAIDalle       bool     `json:"enable_openai_dalle" yaml:"enable_openai_dalle"`
	DalleUsingOpenAISetting bool     `json:"dalle_using_openai_setting" yaml:"dalle_using_openai_setting"`
	OpenAIDalleAzure        bool     `json:"openai_dalle_azure" yaml:"openai_dalle_azure"`
	OpenAIDalleAPIVersion   string   `json:"openai_dalle_api_version" yaml:"openai_dalle_api_version"`
	OpenAIDalleAutoProxy    bool     `json:"openai_dalle_auto_proxy" yaml:"openai_dalle_auto_proxy"`
	OpenAIDalleOrganization string   `json:"openai_dalle_organization" yaml:"openai_dalle_organization"`
	OpenAIDalleServers      []string `json:"openai_dalle_servers" yaml:"openai_dalle_servers"`
	OpenAIDalleKeys         []string `json:"openai_dalle_keys" yaml:"openai_dalle_keys"`

	// OpenAI Fallback 配置
	EnableFallbackOpenAI       bool     `json:"enable_fallback_openai" yaml:"enable_fallback_openai"`
	FallbackOpenAIAzure        bool     `json:"fallback_openai_azure" yaml:"fallback_openai_azure"`
	FallbackOpenAIServers      []string `json:"fallback_openai_servers" yaml:"fallback_openai_servers"`
	FallbackOpenAIKeys         []string `json:"fallback_openai_keys" yaml:"fallback_openai_keys"`
	FallbackOpenAIOrganization string   `json:"fallback_openai_organization" yaml:"fallback_openai_organization"`
	FallbackOpenAIAPIVersion   string   `json:"fallback_openai_api_version" yaml:"fallback_openai_api_version"`
	FallbackOpenAIAutoProxy    bool     `json:"fallback_openai_auto_proxy" yaml:"fallback_openai_auto_proxy"`
}
