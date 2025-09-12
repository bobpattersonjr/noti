package command

import (
	"fmt"

	"github.com/bobpattersonjr/noti/service/blink1"
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

func getBlink1(title, message string, v *viper.Viper) notification {
	return &blink1.Notification{
		Brightness: v.GetInt("blink1.brightness"),
		Color:      v.GetString("blink1.color"),
		Delay:      v.GetInt("blink1.delay"),
		Fade:       v.GetInt("blink1.fade"),
		Glimmer:    v.GetInt("blink1.glimmer"),
		Path:       v.GetString("blink1.path"),
		Random:     v.GetInt("blink1.random"),
		Repeats:    v.GetInt("blink1.repeats"),
	}
}
