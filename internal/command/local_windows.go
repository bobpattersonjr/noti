package command

import (
	"fmt"

	"github.com/bobpattersonjr/noti/service/notifyicon"
	"github.com/bobpattersonjr/noti/service/speechsynthesizer"
	"github.com/spf13/viper"
)

func getBanner(title, message string, v *viper.Viper) notification {
	return &notifyicon.Notification{
		BalloonTipTitle: title,
		BalloonTipText:  message,
		BalloonTipIcon:  notifyicon.BalloonTipIconInfo,
	}
}

func getSpeech(title, message string, v *viper.Viper) notification {
	return &speechsynthesizer.Notification{
		Text:  fmt.Sprintf("%s %s", title, message),
		Rate:  3,
		Voice: v.GetString("speechsynthesizer.voice"),
	}
}
