package plaza

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type plaza struct {
	session   *discordgo.Session
	stopChans map[string]chan<- struct{}
}

var bot plaza

func Init() error {
	if token := os.Getenv("DISCORD_TOKEN"); token != "" {
		session, err := discordgo.New("Bot " + token)
		if err != nil {
			return err
		}
		bot.session = session
	} else {
		return errors.New("DISCORD_TOKEN cannot be empty")
	}

	bot.session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	bot.session.AddHandler(handleMessage)
	bot.session.AddHandler(handleVoiceState)

	bot.stopChans = make(map[string]chan<- struct{})

	return nil
}

func Run() error {
	if err := bot.session.Open(); err != nil {
		return err
	}

	bot.session.UpdateListeningStatus("https://plaza.one")

	logrus.Info("started bot")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	logrus.Info("closing connection")

	for _, ch := range bot.stopChans {
		ch <- struct{}{}
	}
	bot.session.Close()

	return nil
}
