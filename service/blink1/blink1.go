package blink1

// use blink1control-tool to control blink(1)
// https://github.com/todbot/blink1-tool

import (
	"fmt"
	"os/exec"
)

// Flags for command line building
type Flags struct {
	// Brightness
	Brightness string
	// color in multiple formats (hex, rgb, etc.) with pattern built in see blink1-tool --list
	Color string
	// Delay between blinks
	Delay string
	// fade time
	Fade string
	// Glimmer number of times
	Glimmer string
	// Path to blink1control-tool
	Path string
	// Flash Random colors
	Random string
	// blink count
	Repeats string
}

const (
	BRIGHTNESS = "brightness"
	COLOR      = "rgb"
	DELAY      = "delay"
	FADE       = "millis"
	GLIMMER    = "glimmer"
	PATH       = "path"
	RANDOM     = "random"
	REPEATS    = "blink"
)

// map of possible colors
var colors = map[string]bool{
	"red":     true,
	"green":   true,
	"blue":    true,
	"cyan":    true,
	"magenta": true,
	"yellow":  true,
	"white":   true,
}

type Notification struct {
	// Brightness
	Brightness int
	// color in multiple formats (hex, rgb, etc.) with pattern built in see blink1control-tool --list
	Color string
	// Delay between blinks
	Delay int
	// fade time
	Fade int
	// Glimmer number of times
	Glimmer int
	// Path to blink1control-tool
	Path string
	// Flash Random colors
	Random int
	// blink count
	Repeats int
}

// Send triggers a blink1control notification.
func (n *Notification) Send() error {
	var args []string
	if n.Brightness > 0 {
		args = append(args, fmt.Sprintf("--%s", BRIGHTNESS))
		args = append(args, fmt.Sprintf("%d", n.Brightness))
	}
	if n.Color != "" {
		// if color is in the list of colors
		if colors[n.Color] {
			args = append(args, fmt.Sprintf("--%s", n.Color))
		} else {
			args = append(args, fmt.Sprintf("--%s", COLOR))
			args = append(args, n.Color)
		}
	}
	if n.Delay > 0 {
		args = append(args, fmt.Sprintf("--%s", DELAY))
		args = append(args, fmt.Sprintf("%d", n.Delay))
	}
	if n.Fade > 0 {
		args = append(args, fmt.Sprintf("--%s", FADE))
		args = append(args, fmt.Sprintf("%d", n.Fade))
	}
	if n.Glimmer > 0 {
		args = append(args, fmt.Sprintf("--%s", GLIMMER))
		args = append(args, fmt.Sprintf("%d", n.Glimmer))
	}
	if n.Random > 0 {
		args = append(args, fmt.Sprintf("--%s", RANDOM))
		args = append(args, fmt.Sprintf("%d", n.Random))
	}
	if n.Repeats != 0 {
		args = append(args, fmt.Sprintf("--%s", REPEATS))
		args = append(args, fmt.Sprintf("%d", n.Repeats))
	}

	args = append(args, "-q")

	path := n.Path
	if path == "" {
		path = "blink1-tool"
	}
	cmd := exec.Command(path, args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("blink1-tool: %w", err)
	}
	return nil
}
