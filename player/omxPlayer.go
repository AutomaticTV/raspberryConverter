package player

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

const connectionTimeout = "3" // time in seconds for the player to connect

var pipeIn io.WriteCloser

// play send omxPlayer a command to start playing an uri.
// If the player is already playing this has no effect.
func playCommand(uri string) error {
	// Check if the player is playing
	_, err := execCommand("sh /var/lib/raspberryConverter/omxController.sh status")
	// If not playing, send play command
	if err != nil {
		fmt.Println("Error geting omxPlayer status (assuming its not playing), LETS PLAY!!!!! \n", uri)
		// Send play command
		cmd := exec.Command("/bin/sh", "-c", uri)
		var err error
		pipeIn, err = cmd.StdinPipe()
		if err != nil {
			return err
		}
		err = cmd.Start()
		if err != nil {
			return err
		}
		fmt.Println("Done sending play command")
		// Force (hidevideo/unhidevideo)???
		execCommand("sleep " + connectionTimeout + " && /var/lib/raspberryConverter/omxController.sh hidevideo")
		execCommand("/var/lib/raspberryConverter/omxController.sh unhidevideo")
	}
	return nil
}

// stopCommand send omxPlayer a command to stop playing a video.
// If the player is already stoppped this has no effect.
func stopCommand() {
	execCommand("/var/lib/raspberryConverter/omxController.sh stop")
}

// statusCommand returns readable information on what is displaying the omxPlayer. It also return boolean indicating if something is being played
func statusCommand() (string, bool) {
	// Check player status
	output, err := execCommand("/var/lib/raspberryConverter/omxController.sh status")
	if err != nil {
		fmt.Println("Error geting omxPlayer status (assuming its not playing)")
		return "Not playing", false
	}
	return string(output), true
}

func execCommand(command string) ([]byte, error) {
	cmd := exec.Command("/bin/sh", "-c", command)
	out, err := cmd.CombinedOutput()
	fmt.Println(command, " OUTPUT =========== \n", string(out))
	if err != nil {
		return nil, errors.New(command + " ERROR: " + err.Error())
	}
	return out, nil
}
