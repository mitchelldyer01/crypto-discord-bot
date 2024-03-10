package config

import (
	"gopkg.in/yaml.v2"
	"io"
)

type Config struct {
	ChannelsWithCoins []ChannelWithCoin `yaml:"channels_with_coins"`
}

type ChannelWithCoin struct {
	ChannelID string `yaml:"channel_id"`
	Coins     []Coin `yaml:"coins"`
}

type Coin struct {
	Ticker  string `yaml:"ticker"`
	Address string `yaml:"address"`
}

func (c *Config) Load(file io.Reader) error {
	stream, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(stream, c)
}
