package plaza

import (
	"net/http"

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

	resp, err := http.Get("http://radio.plaza.one/opus")
	if err != nil {
		logContext.Error(err)
		return
	}
	defer resp.Body.Close()

	encoder, err := dca.EncodeMem(resp.Body, dca.StdEncodeOptions)
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
