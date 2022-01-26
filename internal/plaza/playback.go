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

	opts := dca.StdEncodeOptions
	opts.Bitrate = 96
	opts.RawOutput = true
	opts.CompressionLevel = 5

	for {
		encoder, err := dca.EncodeFile("http://radio.plaza.one/opus", opts)
		if err != nil {
			logContext.Error(err)
			voice.Disconnect()
			return
		}

		voice.Speaking(true)

		done := make(chan error)
		dca.NewStream(encoder, voice, done)

		select {
		case <-ch:
			encoder.Stop()
			encoder.Cleanup()
			voice.Speaking(false)
			voice.Disconnect()
			logContext.Info("playback stopped")
			return
		case err := <-done:
			encoder.Stop()
			encoder.Cleanup()
			voice.Speaking(false)
			if err == nil {
				voice.Disconnect()
				logContext.Info("playback stopped")
				return
			}
		}
	}
}
