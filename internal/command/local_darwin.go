package command

import (
	"fmt"

	"github.com/bobpattersonjr/noti/service/blink1"
	"github.com/bobpattersonjr/noti/service/nsuser"
	"github.com/bobpattersonjr/noti/service/say"
	"github.com/spf13/viper"
)

func getBanner(title, message string, v *viper.Viper) notification {
	return &nsuser.Notification{
		Title:           title,
		InformativeText: message,
		ContentImage:    v.GetString("banner.icon"),
		SoundName:       v.GetString("nsuser.soundName"),
	}
}

func getSpeech(title, message string, v *viper.Viper) notification {
	return &say.Notification{
		Voice: v.GetString("say.voice"),
		Text:  fmt.Sprintf("%s %s", title, message),
		Rate:  200,
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
