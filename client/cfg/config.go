package cfg

import (
	"flag"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	EngineConfig   `yaml:"engine"`
	BotConfig      `yaml:"bot"`
	DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type EngineConfig struct {
	EngineType     string `yaml:"type"`
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	Account        string `yaml:"account"`
	Password       string `yaml:"password"`
	ConversationId string `yaml:"conversationId"`
}

type BotConfig struct {
	RateLimiterConfig   `yaml:"rateLimiter"`
	BotType             string `yaml:"type"`
	TelegramBotToken    string `yaml:"token"`
	TelegramChannelName string `yaml:"channelName"`
	TelegramGroupName   string `yaml:"groupName"`
	LogChannelID        int64  `yaml:"logChannel"`
	WechatBotName       string `yaml:"botName"`
	WechatLoginType     string `yaml:"loginType"`
	ShouldLimitPrivate  bool   `yaml:"limitPrivate"`
	ShouldLimitGroup    bool   `yaml:"limitGroup"`
	AdminID             int64  `yaml:"admin"`
}

type RateLimiterConfig struct {
	Capacity int64 `yaml:"capacity"`
	Duration int64 `yaml:"duration"`
	Enable   bool  `yaml:"enable"`
}

func NewConfig() *Config {
	return &Config{}
}

func InitConfig() (*Config, error) {
	c := NewConfig()

	path := os.Getenv("CONFIG_PATH")

	if path == "" {
		flag.StringVar(&path, "c", "./config.yaml", "Your config file path")
		flag.Parse()
	}
	err := c.loadYaml(path)
	if err != nil {
		log.Fatalf("load config err: %v", err)
	}
	return c, nil
}

// Load config from config.yaml
func (c *Config) loadYaml(path string) error {
	yamlFile := path
	data, err := os.ReadFile(yamlFile)
	if nil != err {
		log.Printf("load local yaml err:%v path: %v\n", err, yamlFile)
		return err
	}

	err = yaml.Unmarshal([]byte(data), c)
	if nil != err {
		log.Printf("unmarshal local yaml err:%v path: %v\n", err, yamlFile)
		return err
	}
	return nil
}
