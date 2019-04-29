// INSPIRATION:  https://github.com/dplesca/go-omxremote

package player

import (
	"bytes"
	"fmt"
	"html/template"
	"os/exec"
	"sync"
)

// status is a struct used for string constants that represent the state of the player
type status struct {
	NotInit, Playing, Paused, IP string
}

// player is the struct that controls the playback on the omxplayer
// using an exec.Cmd command
type player struct {
	State   string
	Command *exec.Cmd
	mu      sync.Mutex
}

// killBlocker is a struct used to prevent killing twice the same process
type killBlocker struct {
	waitingToDie bool
	mu           sync.Mutex
}

// playPriority is a struct used to prevent displaying IP to happen when there is a request for play pending
type playPriority struct {
	waitingToPlay bool
	mu            sync.Mutex
}

// Possible values for p.State
const playing = "Playing"
const paused = "Paused"
const runningNothing = "Nothing running"
const displayingIP = "Displaying IP"

// Channel used to controll the player asyncronously and possible values for the actionChannel messages
var actionChannel = make(chan string)

const startAction = "Start"
const stopAction = "Stop"
const restartAction = "Reset"
const errorAction = "Error"
const doneAction = "Finish"

// Init is a function that initializes the player, and the storage
func Init() {
	initStorage()
	go playerController()
}

// Start is a function that plays RTMP streaming according to stored config
// if the player is already streaming this function has no effect.
func Start() {
	actionChannel <- startAction
}

// Restart is a function that plays RTMP streaming according to stored config
// if the player is already streaming it will be stopped before playing as described.
func Restart() {
	actionChannel <- startAction
}

// Stop is a function that terminate the streaming video process, and switch to displaying IP
func Stop() {
	actionChannel <- stopAction
}

// playerController is a function that acts as a concurrency controller fot the player,
// it gets messages from the actionChannel, triggered either by the importer of the package (startAction || restartAction || stopAction) or by other actions (doneAction || errorAction).
func playerController() {
	p := player{State: runningNothing}
	pp := playPriority{waitingToPlay: false}
	k := killBlocker{waitingToDie: false}
	// Chose next action to be done
	doAction := func(action string) {
		switch action {
		case startAction:
			start(&p, &k, &pp)
		case restartAction:
			restart(&p, &k, &pp)
		case stopAction:
			killRuningProcess(&p, &k)
			displayIP(&p, &pp)
		case doneAction:
			k.waitingToDie = false
			displayIP(&p, &pp)
		case errorAction:
			k.waitingToDie = false
			displayIP(&p, &pp)
		default:
			displayIP(&p, &pp)
		}
	}

	// Initialize the process with default action
	go doAction("")
	// Endless loop!
	for {
		// wait for new actions
		action := <-actionChannel
		// execute the action asyncronously
		go doAction(action)
	}

}

// VOLUME = -o hdmi --vol [-6000:0]
type playerTemplateData struct {
	Volume int // [-6000:0]
	URL    string
}

// play is a function that starts the stream with the stored setings
func start(p *player, k *killBlocker, pp *playPriority) {
	// IF ALREADY PLAYING RETURN
	if p.State == playing || pp.waitingToPlay {
		fmt.Println("Already playing")
		return
	}
	// KILL ANY RUNNING PROCESS, AND INDICATE THAT WE'VE A PLAY REQUEST PENDING
	pp.mu.Lock()
	defer pp.mu.Unlock()
	pp.waitingToPlay = true
	killRuningProcess(p, k)
	// GET PLAYER CONFIG
	setings, err := GetConfig()
	if err != nil {
		fmt.Println("Error geting config from storage: ", err)
		actionChannel <- errorAction
		return
	}
	// BUILD THE PLAY COMMAND ACCORDING TO CONFIG
	var omxOptionsTemplate = template.Must(template.New("omxCommand").Parse(
		"omxplayer -o hdmi --vol {{.Volume}} {{.URL}}",
	))
	var omxOptions bytes.Buffer
	if err = omxOptionsTemplate.Execute(&omxOptions, playerTemplateData{
		setings.Volume,
		setings.URL,
	}); err != nil {
		fmt.Println("Error generating omxplayer command: ", err)
		actionChannel <- errorAction
		return
	}
	// RUN THE COMMAND
	err = runAction(p, omxOptions.String(), playing, true)
	// NOTIFY THAT THE PROCESS IS DONE VIA actionChannel
	pp.waitingToPlay = false
	if err != nil {
		actionChannel <- errorAction
		return
	}
	actionChannel <- doneAction
	return
}

func restart(p *player, k *killBlocker, pp *playPriority) {
	// KILL ANY RUNNING PROCESS
	killRuningProcess(p, k)
	// PLAY
	start(p, k, pp)
}

// killRuningProcess is a function that terminates the process that is being runned by p
func killRuningProcess(p *player, k *killBlocker) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if p.State != runningNothing && !k.waitingToDie {
		k.waitingToDie = true
		// kill command
		fmt.Println("Killing current process")
		if err := p.Command.Process.Kill(); err != nil {
			fmt.Println("failed to kill process: ", err)
			k.waitingToDie = false
			killRuningProcess(p, k)
		}
		fmt.Println("Process murdered")
	}
}

// displayIP is a function that generates an image containing the IP of the system, and display it through the player
func displayIP(p *player, pp *playPriority) {
	// IF THE PLAYER IS ALREADY DOING SOMETHING, return
	if p.State != runningNothing || pp.waitingToPlay {
		fmt.Println("The player is playing, waiting to play or already displaying IP.")
		return
	}
	// Gen IP IMAGE
	// magic command tells omxplayer to display the IP IMAGE
	// RUN THE COMMAND
	err := runAction(p, `sleep 10`, displayingIP, false)
	// NOTIFY THAT THE PROCESS IS DONE VIA actionChannel
	if err != nil {
		actionChannel <- errorAction
		return
	}
	actionChannel <- doneAction
}

// runAction is a function that execute a command synchronously and once at a time (it blocks)
func runAction(p *player, cmd, newState string, isPlay bool) error {
	p.Command = exec.Command("/bin/sh", "-c", cmd)
	// Block writes on p until the function is done
	fmt.Println("BLOCKING: ", cmd)
	p.mu.Lock()
	defer p.mu.Unlock()
	// Set new state, run the command then set runing nothing state
	p.State = newState
	fmt.Println("Player: executing command START: ", cmd)
	msg, err := p.Command.CombinedOutput()
	p.State = runningNothing
	// Logs
	if err != nil {
		fmt.Println("The command: '", cmd, "', has finished with error: ", msg, err)
	} else {
		fmt.Println("Player: executing command END: ", cmd)
	}
	return err
}
