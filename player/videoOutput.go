package player

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
)

// setVideoOutput will change the output mode of the HDMI port based on the given mode
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

// getCurrentOutputMode return information about the video output that is currently used in the HDMI port
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
