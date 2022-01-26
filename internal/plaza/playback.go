package plaza

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/sirupsen/logrus"
)

func playback(ch <-chan struct{}, voice *discordgo.VoiceConnection) {
	logContext := logrus.WithFields(logrus.Fields{
		"guild":   voice.GuildID,
		"channel": voice.ChannelID,
	})
	logContext.Info("starting playback")

	encoder, err := dca.EncodeFile("http://radio.plaza.one/opus", dca.StdEncodeOptions)
	if err != nil {
		logContext.Error(err)
		voice.Disconnect()
		return
	}
	defer encoder.Cleanup()

	voice.Speaking(true)

	done := make(chan error)
	dca.NewStream(encoder, voice, done)

	<-ch
	encoder.Stop()
	<-done
	logContext.Info("playback stopped")

	voice.Speaking(false)
	voice.Disconnect()
}
