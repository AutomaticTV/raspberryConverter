package player

import (
	"errors"
	"fmt"
	"os/exec"
	"raspberryConverter/services/network"
	"regexp"
	"strconv"
)

// playerLoop is a function that executes commands syncronously in a infinit loop.
// The commands to be executed are set by playerController
func playerLoop(p *player, k *killing) {
	initImageMaker()
	var failCounter int
	const failLimit = 10
	for {
		p.state = p.nextState
		cmd := getNextCommand(p.nextState)
		p.command = exec.Command("/bin/sh", "-c", cmd)
		var err error
		p.pipeIn, err = p.command.StdinPipe()
		if err != nil {
			fmt.Println(err)
		}
		// Set new state, run the command then set runing nothing state
		fmt.Println("Player: executing command START: ", cmd)
		msg, err := p.command.CombinedOutput()
		p.state = runningNothing
		k.mu.Lock()
		k.inProgress = false
		k.mu.Unlock()
		// Logs
		if err != nil {
			failCounter++
			fmt.Println("The command: '", cmd, "', has finished with error: ", string(msg), err)
			if failCounter > failLimit {
				fmt.Println("Too many consecutive failures, trying to restart the player")
				channel <- errorMsg
				return
			}
		} else {
			failCounter = 0
			channel <- doneMsg
			fmt.Println("Player: executing command END: ", cmd)
		}
	}
}

func getNextCommand(nextState string) string {
	const survivalCommand = "sleep 10"
	var cmd string
	var err error
	switch nextState {
	case playing:
		cmd, err = getPlayCommand()
	default:
		cmd, err = getDisplayCommand()
	}
	if err != nil {
		fmt.Println("!!Error while trying to get next command in order to achieve the next state: ", nextState, ":")
		fmt.Println(err)
	}
	if cmd == "" {
		fmt.Printf("next command not found!! Seting survival command!!!")
		return survivalCommand
	}
	return cmd
}

// getPlayCommand is a function that starts the stream with the stored setings
func getPlayCommand() (string, error) {
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
	// transform URL http(s)://(www.)... => rtmp://(username:password@)...
	re := regexp.MustCompile(`(https:\/\/www\.|http:\/\/www\.|rtmp:\/\/www\.|https:\/\/|http:\/\/|rtmp:\/\/|www\.)`)
	var auth string
	var protocol string
	if config.Username != "" || config.Password != "" {
		auth = config.Username + ":" + config.Password + "@"
	}
	if config.Transport == "UDP" {
		protocol = "udp://"
	} else if config.Transport == "TCP" {
		protocol = "tcp://"
	} else {
		protocol = "rtmp://"
	}
	url := re.ReplaceAllString(config.URL, protocol+auth)
	// transform the buffer ms => s
	threshold := "--threshold " + strconv.FormatFloat(float64(config.Buffer)/1000.0, 'f', 3, 64) + " "
	return "omxplayer -o hdmi " + volume + threshold + decode + url + " && sudo killall fbi", nil
}

var lastIP string

// displayIP is a function that generates an image containing the IP of the system, and display it through the player
func getDisplayCommand() (string, error) {
	const cmd = "sudo fbi --noverbose -a -T 7 -d /dev/fb0 " + destinationFile + " && read x < /dev/fd/1"
	// GET IP
	config, err := network.GetConfig()
	if err != nil {
		fmt.Println(config)
		return cmd, errors.New("Error geting config from network: " + err.Error())
	}
	// IF IP HAS CHANGED SINCE LAST IMAGE WAS GENERATED
	if config.IP != lastIP {
		// MAKE A NEW IMAGE
		err = makeImage("http://" + config.IP)
		if err != nil {
			return cmd, errors.New("Error saving the new image: " + err.Error())
		}
		lastIP = config.IP
	}
	// RUN THE COMMAND
	return cmd, nil
}

func setVideoOutput(mode string) error {
	// CHECK IF THE VIDEO OUTPUT IS ALREADY IN THE MODE
	m, err := stringToOutputMode(mode)
	if err != nil {
		return err
	}
	currentM, err := getCurrentOutputMode()
	if err != nil {
		return err
	}
	if m == currentM {
		// VIDEO OUTPUT IS ALREADY IN THE MODE
		return nil
	}
	// CHANGE VIDEO OUTPUT MODE
	msg, err := exec.Command("/bin/sh", "-c", "tvservice -e "+m+" && sleep 1 && fbset -depth 8 && fbset -depth 16").CombinedOutput()
	if err != nil {
		fmt.Println(err, msg)
		return errors.New("Error while changing video outpot mode")
	}
	return nil
}

func getCurrentOutputMode() (string, error) {
	msg, err := exec.Command("/bin/sh", "-c", "tvservice -s").CombinedOutput()
	if err != nil {
		fmt.Println(err, msg)
		return "", errors.New("Error while geting the current video output mode")
	}
	regFind := regexp.MustCompile(`CEA \([0-9]{1,2}\)|DMT \([0-9]{1,2}\)`)
	regClean := regexp.MustCompile(`\(|\)`)
	return `"` + regClean.ReplaceAllString(regFind.FindString(string(msg)), "") + `"`, nil
}

// GetCurrentOutputModeString returns a string that descrives the video output, in a human friendly format.
func GetCurrentOutputModeString() (string, error) {
	msg, err := exec.Command("/bin/sh", "-c", "tvservice -s").CombinedOutput()
	if err != nil {
		fmt.Println(err, msg)
		return "", errors.New("Error while geting the current video output mode")
	}
	reg := regexp.MustCompile(`.*\],`)
	return reg.ReplaceAllString(string(msg), ""), nil
}

// stringToOutputMode transforms mode string as specified by PlayerConfig model to the modes of "tvservice" in raspberry
func stringToOutputMode(mode string) (string, error) {
	var m string
	switch mode {
	case "480i59.94":
		m = `"CEA 6"`
	case "480+i59.94":
		m = `"CEA 7"`
	case "576i50":
		m = `"CEA 21"`
	case "576+i50":
		m = `"CEA 22"`
	case "720p50":
		m = `"CEA 19"`
	case "720p59.94":
		m = `"CEA 4"`
	case "1080i50":
		m = `"CEA 20"`
	case "1080i59.94":
		m = `"CEA 5"`
	case "1080p50":
		m = `"CEA 31"`
	case "1080p59.94":
		m = `"CEA 16"`
	default:
		return "", errors.New("Invalid video output mode")
	}
	return m, nil
}
