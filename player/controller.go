// Package player is in charge of controlling everything is displayed over the HDMI port
package player

import (
	"errors"
	"fmt"
	"raspberryConverter/network"
	"regexp"
	"strconv"
	"time"
)

const whatchdogPeriod = 5                               // seconds between each autoplay check
const destinationPath = "/var/lib/raspberryConverter/"  // path to store images (used to display IP)
const destinationFile = destinationPath + "IPImage.png" // filepath to store images (used to display IP)

// Init is a function that initializes the player, and the storage
func Init() {
	// Create DB table if it doesn't exist. Run migrations if there is a new version
	initStorage()
	// Start with the IP image before trying to stream video
	DisplayImageCommand()
	// Check the status of the player every whatchdogPeriod seconds and decide if show IP image, keep playing or try to play.
	go whatcdog()
}

// Start is a function that makes player start playing according to stored config
// if the player is already streaming this function has no effect.
func Start() {
	fmt.Println("======== Player Controller: START")
	uri, err := getPlayUri()
	if err != nil {
		fmt.Println("Error geting play command: ", err)
		return
	}
	playCommand(uri)
	if err != nil {
		fmt.Println("ERROR executing play command: ", err)
		return
	}
}

// getPlayUri is a function that starts the stream with the stored setings
func getPlayUri() (string, error) {
	config, err := GetConfig()
	if err != nil {
		return "", errors.New("Error geting config from storage: " + err.Error())
	}
	err = setVideoOutput(config.Video)
	if err != nil {
		return "", err
	}
	// BUILD THE PLAY COMMAND ACCORDING TO CONFIG
	// transform volume [0:100] => [-6000:0]
	volume := "--vol " + strconv.Itoa(-6000+60*config.Volume) + " "
	var decode string
	if config.AudioDecoding == "Hardware" {
		decode = "--hw "
	}
	// transform URL http(s)://(www.)... => rtmp://(username:paszsword@)...
	re := regexp.MustCompile(`(https:\/\/www\.|http:\/\/www\.|rtmp:\/\/www\.|https:\/\/|http:\/\/|rtmp:\/\/|www\.)`)
	var auth string
	if config.Username != "" || config.Password != "" {
		auth = config.Username + ":" + config.Password + "@"
	}
	url := re.ReplaceAllString(config.URL, "rtmp://"+auth)
	// transform the buffer ms => s
	threshold := "--threshold " + strconv.FormatFloat(float64(config.Buffer)/1000.0, 'f', 3, 64) + " "
	return "omxplayer -o hdmi " + volume + threshold + decode + url, nil
}

// Restart is a function that makes player start playing according to stored config
// if the player is already streaming it will be stopped before playing as described.
func Restart() {
	fmt.Println("======== Player Controller: RESTART")
	Stop()
	Start()
}

// Stop is a function that terminate the streaming video process, and switch to displaying IP
func Stop() {
	fmt.Println("======== Player Controller: STOP")
	stopCommand()
}

// GetStatus is a function that terminate the streaming video process, and switch to displaying IP
func GetStatus() string {
	statusString, _ := statusCommand()
	return statusString
}

// whatcdog check the status of the player every whatchdogPeriod seconds
// if Autoplay is enabled in PLAYER config, it will try to play. Otherwise will display IP image.
func whatcdog() {
	for {
		_, isPlaying := statusCommand()
		config, _ := GetConfig()
		// Show IP
		if !isPlaying {
			// GET IP
			config, err := network.GetConfig()
			if err != nil {
				fmt.Println(config)
				err = MakeImage("NO INTERNET")
				if err != nil {
					fmt.Println("Error saving the new image: ", err)
				}
			}
			// IF IP HAS CHANGED SINCE LAST IMAGE WAS GENERATED
			if config.IP != LastLabel {
				// MAKE A NEW IMAGE
				err = MakeImage("http://" + config.IP)
				if err != nil {
					fmt.Println("Error saving the new image: ", err)
				}
			}
			err = DisplayImageCommand()
			if err != nil {
				fmt.Println("Error displaying image:", err)
			}
		}
		// if autoplay: Play
		if config.Autoplay == "Yes" && !isPlaying {
			Start()
		}
		// if is playing: Keep playing
		if isPlaying {
			fmt.Println("======== Player Controller: KEEP PLAYING")
		}
		time.Sleep(whatchdogPeriod * time.Second)
	}
}
