package player

import (
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/gobuffalo/packr"
)

// player is the struct that controls the playback on the omxplayer
// using an exec.Cmd command
type player struct {
	state     string
	nextState string
	command   *exec.Cmd
	pipeIn    io.WriteCloser
}

type killing struct {
	inProgress bool
	mu         sync.Mutex
}

// static image to be used in displayingIP mode.
var box = packr.NewBox("./assets")

// channel used to controll the player and possible values for the messages
var channel = make(chan string)
var statusChannel = make(chan string)

const startMsg = "Start"
const stopMsg = "Stop"
const restartMsg = "Reset"
const statusMsg = "Status"
const doneMsg = "Done"
const errorMsg = "Error"

// Possible values for p.State and p.nextState
const playing = "Playing"
const displayingIP = "Displaying IP"
const runningNothing = "Nothing running"
const defaultState = displayingIP

const autoPlayPeriod = 30 // seconds between each autoplay check

// Init is a function that initializes the player, and the storage
func Init() {
	// Create DB table if it doesn't exist. Run migrations if there is a new version
	initStorage()
	// Endless loop that gets requests concurrently and controll the player (start, stop, ...).
	go playerController(displayingIP)
	// Sends start message to player controller every autoPlayPeriod seconds if Autoplay is enabled
	go autoPlay()
}

// Start is a function that makes player start playing according to stored config
// if the player is already streaming this function has no effect.
func Start() {
	channel <- startMsg
}

// Restart is a function that makes player start playing according to stored config
// if the player is already streaming it will be stopped before playing as described.
func Restart() {
	channel <- restartMsg
}

// Stop is a function that terminate the streaming video process, and switch to displaying IP
func Stop() {
	channel <- stopMsg
}

// GetStatus is a function that terminate the streaming video process, and switch to displaying IP
func GetStatus() string {
	channel <- statusMsg
	return <-statusChannel
}

// playerController is a function that acts as a concurrency controller fot the player,
// it gets messages from the channel, triggered either by the importer of the package (startMsg || restartMsg || stopMsg) or by the player loop (errorMsg || doneMsg).
// according to the received message and the current state of the player it will decide to (stop playing | start playing | stop playing and then start again).
// Note that when the player is not playing the system displays a static image that shows the IP of the device.
func playerController(initialState string) {
	// Initialize the process with default action
	p := player{state: runningNothing, nextState: initialState}
	k := killing{}
	go playerLoop(&p, &k)
	// Endless channel reader loop:
	var err error
	for {
		// wait for new mwssage
		fmt.Println(p)
		msg := <-channel
		// decide the next state of the player based on the message
		switch msg {
		case startMsg:
			if p.state != playing {
				p.nextState = playing
				err = killRuningProcess(&p, &k)
			}
		case restartMsg:
			p.nextState = playing
			err = killRuningProcess(&p, &k)
		case stopMsg:
			p.nextState = displayingIP
			err = killRuningProcess(&p, &k)
		case statusMsg:
			statusChannel <- p.state
		case doneMsg:
			p.nextState = displayingIP
		case errorMsg:
			fmt.Println("The player have experimented an error. Initializing again.")
			err = killRuningProcess(&p, &k)
			p = player{state: runningNothing, nextState: displayingIP}
			go playerLoop(&p, &k)
		default:
			// IN THEORY THIS CASE SHOULD NEVER HAPPEN!
			fmt.Println("Unrecognized message received from the channel: ", msg)
			if p.state == runningNothing {
				p.nextState = defaultState
			}
		}
		// ERROR KILLING PROCESS, IN THEORY THIS CASE SHOULD NEVER HAPPEN!
		// if err != nil {
		// 	go func(err error) {
		// 		for err != nil {
		fmt.Println("Error killing process!", err)
		// 			// TODO: HERE A BIT OF PROPPER ERROR HANDLING SHOULD HAPPEN (EVALUATE IF IT'S REALLY NEEDED).
		// 			err = killRuningProcess(&p, &k)
		// 		}
		// 	}(err)
		// }
	}
}

// killRuningProcess is a function that terminates the process that is being runned by p
func killRuningProcess(p *player, k *killing) error {
	fmt.Println("killer mode")
	fmt.Println(p)
	if p.state == playing {
		fmt.Println("sending q")
		p.pipeIn.Write([]byte("q"))
		return nil
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	if p.state != runningNothing && !k.inProgress {
		k.inProgress = true
		// kill command
		fmt.Println("Killing current process")
		if err := p.command.Process.Kill(); err != nil {
			k.inProgress = false
			fmt.Println("failed to kill process: " + err.Error())
			return err
		}
		fmt.Println("Process murdered")
	}
	return nil
}

// autoPlay simulates a "Start" received from the web UI every autoPlayPeriod seconds
// if Autoplay is enabled in PLAYER config.
func autoPlay() {
	for {
		time.Sleep(autoPlayPeriod * time.Second)
		config, _ := GetConfig()
		if config.Autoplay == "Yes" {
			channel <- startMsg
		}
	}
}
