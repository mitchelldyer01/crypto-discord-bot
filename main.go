package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/castillobgr/sententia"
	"github.com/mitchelldyer01/crypto-discord-bot/config"
)

func main() {
	path := flag.String("path", "/config.yaml", "A path to the config file")

	flag.Parse()

	err := run(*path)
	if err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
}

func run(path string) error {
	last_time := time.Now()
	// open config file at path
	fmt.Printf("[INFO] Opening %s...\n", path)
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	fmt.Printf("[INFO] Loading %s...\n", path)
	// load file into config
	cfg := config.Config{}
	err = cfg.Load(file)
	if err != nil {
		return err
	}

	fmt.Print("[INFO] ...Loaded!\n")

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		fmt.Print("[ERROR] DISCORD_TOKEN not found")
	}

	fmt.Print("[INFO] Connecting to discord...\n")
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}
	fmt.Print("[INFO] ...Connected!\n")

	client := &http.Client{}

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		current_time := time.Now()

		if current_time.Before(last_time.Add(31 * time.Second)) {
			return
		}

		last_time = current_time

		if m.Author.ID == s.State.User.ID {
			return
		}

		// only bother with messages in channels in our config
		for _, channel := range cfg.ChannelsWithCoins {
			if channel.ChannelID == m.ChannelID {
				// only bother with messages containing one of our coins
				for _, coin := range channel.Coins {
					if strings.Contains(strings.ToLower(m.Content), strings.ToLower(coin.Ticker)) {
						url := fmt.Sprintf("https://api.ethplorer.io/getTokenInfo/%s?apiKey=freekey", coin.Address)

						res, err := client.Get(url)
						if err != nil {
							fmt.Printf("[ERROR]: %s\n", err)
							return
						}

						pricer := ethplorer{}

						err = json.NewDecoder(res.Body).Decode(&pricer)
						if err != nil {
							fmt.Printf("[ERROR]: %s\n", err)
							return
						}

						price := fmt.Sprintf("%.20f", pricer.Price.Rate)

						adj, err := sententia.Make("{{adjective}}")
						if err != nil {
							fmt.Printf("[ERROR]: %s\n", err)
							return
						}

						msg := fmt.Sprintf("%s is at $%s _(the bull is %s)_", coin.Ticker, price, adj)
						_, err = session.ChannelMessageSend(m.ChannelID, msg)
						if err != nil {
							fmt.Printf("[ERROR]: %s\n", err)
							return
						}
						fmt.Printf("[INFO]: %s\n", msg)
					}
				}
			}
		}

	})
	session.Identify.Intents = discordgo.IntentsGuildMessages

	err = session.Open()
	if err != nil {
		return err
	}
	fmt.Printf("[INFO] Listening for messages...\n")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	session.Close()
	return nil
}

type ethplorer struct {
	Price ethplorer_price `json:"price"`
}

type ethplorer_price struct {
	Rate float64 `json:"rate"`
}
