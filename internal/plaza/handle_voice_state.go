package plaza

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func handleVoiceState(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	if vsu.UserID == bot.session.State.User.ID {
		return
	}

	logContext := logrus.WithFields(logrus.Fields{
		"guild": vsu.GuildID,
	})
	logContext.Info("voice state changed")

	if vsu.ChannelID == "" {
		logContext.Info("user left, checking if they were the last")

		guild, err := s.State.Guild(vsu.GuildID)
		if err != nil {
			logContext.Error(err)
			return
		}

		if len(guild.VoiceStates) == 1 && guild.VoiceStates[0].UserID == s.State.User.ID {
			logContext.Info("all users left, disconnecting")

			if ch, ok := bot.stopChans[vsu.GuildID]; ok {
				ch <- struct{}{}
				close(ch)
				delete(bot.stopChans, vsu.GuildID)
			}
		}
	}
}
