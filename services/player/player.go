// INSPIRATION:  https://github.com/dplesca/go-omxremote

package player

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"os/exec"
)

// status is a struct used for string constants that represent the state of the player
type status struct {
	NotInit, Playing, Paused, IP string
}

// player is the struct that controls the playback on the omxplayer
// using an exec.Cmd command
type player struct {
	Command *exec.Cmd
	PipeIn  io.WriteCloser
	Playing bool
	State   string
}

var s = status{
	NotInit: "N",
	Playing: "P",
	Paused:  "p",
	IP:      "I",
}

var p = player{}

// Init is a function that initializes the player, and the storage
func Init() {
	initStorage()
	initPlayer()
}

// initPlayer is a function that intializes the player acccording to user setings
func initPlayer() error {
	setings, err := GetConfig()
	if err != nil {
		return playerError(err)
	}
	if setings.Autoplay == "Yes" {
		err = Play()
		if err != nil {
			return playerError(err)
		}
	} else {
		err = displayIP()
		if err != nil {
			return playerError(err)
		}
	}
	return nil
}

// VOLUME = -o hdmi --vol [-6000:0]
type playerTemplateData struct {
	Volume int // [-6000:0]
	URL    string
}

var omxOptionsTemplate = template.Must(template.New("omxCommand").Parse(
	"omxplayer -o hdmi --vol {{.Volume}} {{.URL}} &",
))

//Play is a function that starts the stream with the stored setings
func Play() error {
	// IF ALREADY PLAYING RETURN
	// if p.Playing {
	// 	return errors.New("Already playing")
	// }
	// setings, err := GetConfig()
	// if err != nil {
	// 	return playerError(err)
	// }
	// var omxOptions bytes.Buffer
	// if err := omxOptionsTemplate.Execute(&omxOptions, playerTemplateData{
	// 	-3000,
	// 	"/home/pi/tmp",
	// }); err != nil {
	// 	return playerError(err)
	// }
	//
	// p.Command = exec.Command("/bin/sh", "-c", omxOptions.String())
	// p.PipeIn, err = p.Command.StdinPipe()
	// if err != nil {
	// 	return playerError(err)
	// }
	//
	// p.Playing = true
	// err = p.Command.Start()
	//
	// if err != nil {
	// 	p.Playing = false
	// 	return playerError(err)
	// }

	return nil
}

func Stop() error {
	return nil
}

func Pause() error {
	return nil
}

func displayIP() error {
	return nil
}

func playerError(err error) error {
	fmt.Println(err)
	return errors.New("Error geting network configuration.")
}
