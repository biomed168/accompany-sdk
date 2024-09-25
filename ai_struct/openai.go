package ai_struct

type AiConfig struct {
	OpenAiConfig `json:"openAiConfig"`
}

// OpenAiConfig 配置结构体包含与 OpenAI、DALL·E 以及 Fallback 相关的配置选项。
// 该结构体用于控制 OpenAI 和相关服务的启用、API 版本、组织配置、代理设置等功能。
type OpenAiConfig struct {
	// EnableOpenAI 控制是否启用 OpenAI 服务。为 true 时表示启用。
	EnableOpenAI bool `json:"enable_openai" yaml:"enable_openai"`

	// OpenAIAzure 指示是否使用 Azure 提供的 OpenAI 服务。为 true 时表示使用 Azure 版本。
	OpenAIAzure bool `json:"openai_azure" yaml:"openai_azure"`

	// OpenAIAPIVersion 指定要使用的 OpenAI API 版本（如 "v1"）。可根据需要设置特定版本。
	OpenAIAPIVersion string `json:"openai_api_version" yaml:"openai_api_version"`

	// OpenAIAutoProxy 控制是否为 OpenAI 服务启用自动代理，适用于需要通过代理访问 OpenAI 的场景。
	OpenAIAutoProxy bool `json:"openai_auto_proxy" yaml:"openai_auto_proxy"`

	// OpenAIOrganization 指定 OpenAI 服务所属的组织 ID，通常用于多租户或企业账户。
	OpenAIOrganization string `json:"openai_organization" yaml:"openai_organization"`

	// OpenAIServers 定义 OpenAI 服务器的地址列表，可以配置多个服务器以供选择。
	OpenAIServers []string `json:"openai_servers" yaml:"openai_servers"`

	// OpenAIKeys 存储 OpenAI API 的密钥，支持配置多个密钥以便在不同环境下使用。
	OpenAIKeys []string `json:"openai_keys" yaml:"openai_keys"`

	// EnableOpenAIDalle 控制是否启用 DALL·E 服务（用于生成图像的 OpenAI 模型）。为 true 时表示启用。
	EnableOpenAIDalle bool `json:"enable_openai_dalle" yaml:"enable_openai_dalle"`

	// DalleUsingOpenAISetting 指定 DALL·E 是否使用通用的 OpenAI 配置。如果为 true，DALL·E 将使用 OpenAI 的设置。
	DalleUsingOpenAISetting bool `json:"dalle_using_openai_setting" yaml:"dalle_using_openai_setting"`

	// OpenAIDalleAzure 指示是否使用 Azure 提供的 DALL·E 服务。为 true 时表示使用 Azure 版本。
	OpenAIDalleAzure bool `json:"openai_dalle_azure" yaml:"openai_dalle_azure"`

	// OpenAIDalleAPIVersion 指定要使用的 DALL·E API 版本。
	OpenAIDalleAPIVersion string `json:"openai_dalle_api_version" yaml:"openai_dalle_api_version"`

	// OpenAIDalleAutoProxy 控制是否为 DALL·E 启用自动代理，适用于需要代理访问的场景。
	OpenAIDalleAutoProxy bool `json:"openai_dalle_auto_proxy" yaml:"openai_dalle_auto_proxy"`

	// OpenAIDalleOrganization 指定 DALL·E 服务所属的组织 ID，通常用于多租户或企业账户。
	OpenAIDalleOrganization string `json:"openai_dalle_organization" yaml:"openai_dalle_organization"`

	// OpenAIDalleServers 定义 DALL·E 服务的服务器地址列表，支持多个服务器配置。
	OpenAIDalleServers []string `json:"openai_dalle_servers" yaml:"openai_dalle_servers"`

	// OpenAIDalleKeys 存储 DALL·E API 的密钥，支持配置多个密钥。
	OpenAIDalleKeys []string `json:"openai_dalle_keys" yaml:"openai_dalle_keys"`

	// EnableFallbackOpenAI 控制是否启用备用的 OpenAI 服务，当主服务不可用时切换到备用服务。
	EnableFallbackOpenAI bool `json:"enable_fallback_openai" yaml:"enable_fallback_openai"`

	// FallbackOpenAIAzure 指示备用服务是否使用 Azure 提供的 OpenAI 服务。为 true 时表示使用 Azure 版本。
	FallbackOpenAIAzure bool `json:"fallback_openai_azure" yaml:"fallback_openai_azure"`

	// FallbackOpenAIServers 定义备用 OpenAI 服务的服务器地址列表。
	FallbackOpenAIServers []string `json:"fallback_openai_servers" yaml:"fallback_openai_servers"`

	// FallbackOpenAIKeys 存储备用 OpenAI API 的密钥，支持多个密钥配置。
	FallbackOpenAIKeys []string `json:"fallback_openai_keys" yaml:"fallback_openai_keys"`

	// FallbackOpenAIOrganization 指定备用 OpenAI 服务所属的组织 ID。
	FallbackOpenAIOrganization string `json:"fallback_openai_organization" yaml:"fallback_openai_organization"`

	// FallbackOpenAIAPIVersion 指定备用 OpenAI API 的版本。
	FallbackOpenAIAPIVersion string `json:"fallback_openai_api_version" yaml:"fallback_openai_api_version"`

	// FallbackOpenAIAutoProxy 控制是否为备用 OpenAI 启用自动代理，适用于网络访问受限的场景。
	FallbackOpenAIAutoProxy bool `json:"fallback_openai_auto_proxy" yaml:"fallback_openai_auto_proxy"`
}
