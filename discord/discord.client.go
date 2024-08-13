package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/kfukue/lyle-labs-libraries/v2/utils"
)

const (
	DISCORD_PATH = "DISCORD_LYLE_PATH"
)

var DiscordClientSession *discordgo.Session

func init() {
	initializeClient()
}

func initializeClient() {
	var err error
	botToken, err := GetDiscordBotToken()
	DiscordClientSession, err = discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func GetDiscordBotToken() (string, error) {
	if utils.GetEnv() == "production" {
		log.Println("In production!")
		discordPath := utils.MustGetenv(DISCORD_PATH)
		// Set each value dynamically w/ Sprintf
		token, err := utils.AccessSecretVersion(discordPath)
		if err != nil {
			log.Fatal(err)
		}
		return token, nil
	} else {
		discordPath := utils.GoDotEnvVariable(DISCORD_PATH)
		token, err := utils.AccessSecretVersion(discordPath)
		if err != nil {
			log.Fatal(err)
		}
		return token, nil
	}
}

func SendMessage(channelID string, message string) error {
	if DiscordClientSession == nil {
		initializeClient()
	}
	DiscordClientSession.Identify.Intents = discordgo.IntentGuildMembers
	_, err := DiscordClientSession.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
