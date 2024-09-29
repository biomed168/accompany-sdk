package openai

import (
	"accompany-sdk/ai_struct"
	"accompany-sdk/pkg/proxy"
	"accompany-sdk/pkg/ternary"
	"github.com/sashabaranov/go-openai"
	"net"
	"net/http"
	"regexp"
	"time"
)

func NewOpenAi(conf *ai_struct.OpenAiConfig) OpenAi {
	var mainClient, backupClient Client
	var proxyDialer *proxy.Proxy
	if proxy.ShouldLoad(&conf.ProxyConfig) {
		proxyDialer = proxy.NewProxy(&conf.ProxyConfig)
	}
	if conf.EnableOpenAI {
		mainClient = NewOpenAIClient(parseMainConfig(conf), proxyDialer)
	}

	if conf.EnableFallbackOpenAI {
		backupClient = NewOpenAIClient(parseBackupConfig(conf), proxyDialer)
	}

	return NewOpenAIProxy(mainClient, backupClient)
}

// OpenAi 是一个接口 Client
type OpenAi = Client

func NewOpenAIClient(conf *Config, pp *proxy.Proxy) Client {
	clients := make([]*openai.Client, 0)

	// 如果是 Azure API，则每一个 Server 对应一个 Key
	// 否则 Servers 和 Keys 取笛卡尔积
	if conf.OpenAIAzure {
		for i, server := range conf.OpenAIServers {
			clients = append(clients, createOpenAIClient(
				true,
				conf.OpenAIAPIVersion,
				server,
				"",
				conf.OpenAIKeys[i],
				ternary.If(conf.AutoProxy, pp, nil),
			))
		}
	} else {
		for _, server := range conf.OpenAIServers {
			for _, key := range conf.OpenAIKeys {
				clients = append(clients, createOpenAIClient(
					false,
					"",
					server,
					conf.OpenAIOrganization,
					key,
					ternary.If(conf.AutoProxy, pp, nil),
				))
			}
		}
	}

	return New(conf, clients)
}

func createOpenAIClient(isAzure bool, apiVersion string, server, organization, key string, pp *proxy.Proxy) *openai.Client {
	openaiConf := openai.DefaultConfig(key)
	openaiConf.BaseURL = server
	openaiConf.OrgID = organization

	if pp != nil {
		openaiConf.HTTPClient = &http.Client{
			Transport: pp.BuildTransport(),
			Timeout:   180 * time.Second,
		}

	} else {
		openaiConf.HTTPClient = &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: 120 * time.Second,
				}).DialContext,
			},
			Timeout: 180 * time.Second,
		}
	}

	if isAzure {
		openaiConf.APIType = openai.APITypeAzure
		openaiConf.APIVersion = apiVersion
		openaiConf.AzureModelMapperFunc = func(model string) string {
			// TODO 应该使用配置文件配置，注意，这里返回的应该是 Azure 部署名称
			switch model {
			case "gpt-3.5-turbo", "gpt-3.5-turbo-0613":
				return "gpt35-turbo"
			case "gpt-3.5-turbo-16k", "gpt-3.5-turbo-16k-0613":
				return "gpt35-turbo-16k"
			case "gpt-4", "gpt-4-0613":
				return "gpt4"
			case "gpt-4-32k", "gpt-4-32k-0613":
				return "gpt4-32k"
			}

			return regexp.MustCompile(`[.:]`).ReplaceAllString(model, "")
		}
	}

	return openai.NewClientWithConfig(openaiConf)
}

type Config struct {
	Enable             bool
	OpenAIAzure        bool
	OpenAIAPIVersion   string
	OpenAIOrganization string
	OpenAIServers      []string
	OpenAIKeys         []string
	AutoProxy          bool
}

func parseMainConfig(conf *ai_struct.OpenAiConfig) *Config {
	return &Config{
		Enable:             conf.EnableOpenAI,
		OpenAIAzure:        conf.OpenAIAzure,
		OpenAIAPIVersion:   conf.OpenAIAPIVersion,
		OpenAIOrganization: conf.OpenAIOrganization,
		OpenAIServers:      conf.OpenAIServers,
		OpenAIKeys:         conf.OpenAIKeys,
		AutoProxy:          conf.OpenAIAutoProxy,
	}
}

func parseBackupConfig(conf *ai_struct.OpenAiConfig) *Config {
	return &Config{
		Enable:             conf.EnableFallbackOpenAI,
		OpenAIAzure:        conf.FallbackOpenAIAzure,
		OpenAIAPIVersion:   conf.FallbackOpenAIAPIVersion,
		OpenAIOrganization: conf.FallbackOpenAIOrganization,
		OpenAIServers:      conf.FallbackOpenAIServers,
		OpenAIKeys:         conf.FallbackOpenAIKeys,
		AutoProxy:          conf.FallbackOpenAIAutoProxy,
	}
}

func parseDalleConfig(conf *ai_struct.OpenAiConfig) *Config {
	if conf.DalleUsingOpenAISetting {
		return &Config{
			Enable:             conf.EnableOpenAI && conf.EnableOpenAIDalle,
			OpenAIAzure:        conf.OpenAIAzure,
			OpenAIAPIVersion:   conf.OpenAIAPIVersion,
			OpenAIOrganization: conf.OpenAIOrganization,
			OpenAIServers:      conf.OpenAIServers,
			OpenAIKeys:         conf.OpenAIKeys,
			AutoProxy:          conf.OpenAIAutoProxy,
		}
	}

	return &Config{
		Enable:             conf.EnableOpenAIDalle,
		OpenAIAzure:        conf.OpenAIDalleAzure,
		OpenAIAPIVersion:   conf.OpenAIDalleAPIVersion,
		OpenAIOrganization: conf.OpenAIDalleOrganization,
		OpenAIServers:      conf.OpenAIDalleServers,
		OpenAIKeys:         conf.OpenAIDalleKeys,
		AutoProxy:          conf.OpenAIDalleAutoProxy,
	}
}
