//go:build !darwin && !windows
// +build !darwin,!windows

package command

import (
	"fmt"

	"github.com/bobpattersonjr/noti/service/blink1"
	"github.com/bobpattersonjr/noti/service/espeak"
	"github.com/bobpattersonjr/noti/service/freedesktop"
	"github.com/spf13/viper"
)

func getBanner(title, message string, v *viper.Viper) notification {
	icon := v.GetString("banner.icon")
	if icon == "" {
		icon = "utilities-terminal"
	}
	return &freedesktop.Notification{
		Summary:       title,
		Body:          message,
		ExpireTimeout: 5000,
		AppIcon:       icon,
	}
}

func getSpeech(title, message string, v *viper.Viper) notification {
	return &espeak.Notification{
		Text:      fmt.Sprintf("%s %s", title, message),
		VoiceName: v.GetString("espeak.voiceName"),
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
