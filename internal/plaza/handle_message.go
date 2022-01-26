package plaza

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func handleMessage(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == s.State.User.ID {
		return
	}

	logContext := logrus.WithFields(logrus.Fields{
		"guild":   msg.GuildID,
		"content": msg.Content,
	})
	logContext.Info("got message")

	switch msg.Content {
	case "!play":
		logContext.Info("playback requested")

		guild, err := s.State.Guild(msg.GuildID)
		if err != nil {
			logContext.Error(err)
			return
		}

		for _, vs := range guild.VoiceStates {
			if vs.UserID == msg.Author.ID {
				voice, err := s.ChannelVoiceJoin(msg.GuildID, vs.ChannelID, false, true)
				if err != nil {
					logContext.Error(err)
					return
				}

				ch := make(chan struct{}, 1)
				bot.stopChans[msg.GuildID] = ch
				go playback(ch, voice)

				vChan, err := s.Channel(vs.ChannelID)
				if err != nil {
					logContext.Error(err)
					return
				}
				s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Now playing in %s. Enjoy!", vChan.Name))
				return
			}
		}

		s.ChannelMessageSend(msg.ChannelID, "You are not in voice channel. Where to join??")
	case "!stop":
		logContext.Info("stop requested")

		if ch, ok := bot.stopChans[msg.GuildID]; ok {
			ch <- struct{}{}
			close(ch)
			delete(bot.stopChans, msg.GuildID)
		}

		s.ChannelMessageSend(msg.ChannelID, "Stopped playback")
	}
}
