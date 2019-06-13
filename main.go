// Package raspberryConverter is a JEOS (Just Enought Operating System) to play RTMP on Raspberry Pi 3.
// The project consist of two main elements: a web server that provides a web interface to configure the service and a player to show the streaming content through HDMI.
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"raspberryConverter/auth"
	"raspberryConverter/frontend"
	"raspberryConverter/models"
	"raspberryConverter/monitor"
	"raspberryConverter/network"
	"raspberryConverter/player"
	"strconv"
	"time"

	"github.com/gobuffalo/packr"
)

// port describe the port listened by the server, example: http://localhost:5555
const port = ":1234"

// handler processes the requests received by the server and respon to them
// according to the content of the request.
func handler(w http.ResponseWriter, r *http.Request) {
	// IF NOT LOGGED IN, GO TO LOGIN PAGE
	if !auth.IsLoggedIn(r) && r.URL.Path != "/login" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// CHECK IF REQUEST INCLUDES ANY FORM (AND TRIGGER NECESSARY ACTIONS ACCORDING TO)
	isUpdate, requestError := requestHandler(w, r)

	// RETURN THE APROPIATE FRONTEND PAGE ACCORDING TO THE PATH
	switch r.URL.Path {
	case "/login":
		if requestError == nil && isUpdate {
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		}
		err := frontend.Login(isUpdate, requestError, w)
		if err != nil {
			http.Redirect(w, r, "/static/error.html", http.StatusFound)
		}
		return
	case "/dashboard/status":
		status, err := monitor.GetStatus()
		err = frontend.Status(isUpdate, status, combineErrors(requestError, err), w)
		if err != nil {
			http.Redirect(w, r, "/static/error.html", http.StatusFound)
		}
		return
	case "/dashboard/player":
		config, err := player.GetConfig()
		err = frontend.Player(isUpdate, config, combineErrors(requestError, err), w)
		if err != nil {
			http.Redirect(w, r, "/static/error.html", http.StatusFound)
		}
		return
	case "/dashboard/network":
		config, err := network.GetConfig()
		err = frontend.Network(isUpdate, config, combineErrors(requestError, err), w)
		if err != nil {
			http.Redirect(w, r, "/static/error.html", http.StatusFound)
		}
		return
	case "/dashboard/password":
		err := frontend.Password(isUpdate, requestError, w)
		if err != nil {
			http.Redirect(w, r, "/static/error.html", http.StatusFound)
		}
		return
	default:
		http.Redirect(w, r, "/dashboard/status", http.StatusFound)
		return
	}
}

// requestHandler checks if the request includes an expected form and if so it gets handled.
// Returns true and a possible error if r includes a form otherwisr (false, nil)
func requestHandler(w http.ResponseWriter, r *http.Request) (bool, error) {
	// LOGIN
	if r.FormValue("Login") != "" && r.FormValue("Username") != "" {
		username := r.FormValue("Username")
		password := r.FormValue("Password")
		if auth.PasswordIsCorrect(username, password) {
			return true, auth.Login(w, r, username)
		}
		return true, errors.New("Wrong username or password")
	}
	// PLAY
	if r.FormValue("Start") != "" {
		player.Start()
		return true, nil
	}
	// RESTART
	if r.FormValue("Restart") != "" {
		player.Restart()
		return true, nil
	}
	// STOP
	if r.FormValue("Stop") != "" {
		player.Stop()
		return true, nil
	}
	// UPDATE PLAYER CONFIG
	if r.FormValue("UpdatePlayer") != "" {
		buf, err := strconv.Atoi(r.FormValue("Buffer"))
		var vol int
		if err == nil {
			vol, err = strconv.Atoi(r.FormValue("Volume"))
		}
		if err == nil {
			err = player.SetConfig(models.PlayerConfig{
				Video:         r.FormValue("Video"),
				AudioDecoding: r.FormValue("AudioDecoding"),
				URL:           r.FormValue("URL"),
				Buffer:        buf,
				Username:      r.FormValue("Username"),
				Password:      r.FormValue("Password"),
				Volume:        vol,
				Autoplay:      r.FormValue("Autoplay"),
			})
		}
		return true, err
	}
	// UPDATE NETWORK CONFIG
	if r.FormValue("UpdateNetwork") != "" {
		return true, network.SetConfig(models.NetworkConfig{
			Mode:    r.FormValue("Mode"),
			IP:      r.FormValue("IP"),
			Gateway: r.FormValue("Gateway"),
			Netmask: r.FormValue("Netmask"),
			DNS1:    r.FormValue("DNS1"),
			DNS2:    r.FormValue("DNS2"),
		})
	}
	// CHANGE PASSWORD
	if r.FormValue("UpdatePassword") != "" {
		return true, auth.UpdatePassword(
			r.FormValue("Username"),
			r.FormValue("OldPassword"),
			r.FormValue("NewPassword"),
			r.FormValue("RePassword"),
		)
	}
	// NO REQUEST FOUND
	return false, nil
}

// combineErrors return a single error messages given two errors (that may be nil)
func combineErrors(err1 error, err2 error) error {
	var errorString string
	if err1 != nil {
		errorString = err1.Error() + ". "
	}
	if err2 != nil {
		errorString += err2.Error() + "."
	}
	if errorString != "" {
		return errors.New(errorString)
	}
	return nil
}

// logoutHandler clean the session associated to the user of the request and redirrects to login page
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	err := auth.Logout(w, r)
	if err != nil {
		log.Println(err)
	}
	// redirec to login
	http.Redirect(w, r, "/login", http.StatusFound)
}

// main get executed at the beginning, and initializes the web server and the rest of the components that require initialization.
func main() {
	staticFiles := packr.NewBox("frontend/static")
	frontend.Init()
	auth.Init()
	fmt.Println("Server starting, point your browser to localhost" + port + " to start")
	// ENDPOINTS
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticFiles)))
	http.HandleFunc("/", handler)
	http.HandleFunc("/login", handler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/dashboard/", handler)
	err := http.ListenAndServe(port, nil)
	const cmd = "sudo fbi --noverbose -a -T 7 -d /dev/fb0 /var/lib/raspberryConverter/IPImage.png"
	for err != nil {
		player.LastIP = "fail"
		player.MakeImage("NETWORK ERROR / NO INTERNET")
		exec.Command("/bin/sh", "-c", cmd).Run()
		fmt.Println("Error starting http server: ", err)
		fmt.Println("trying again in 5 seconds")
		time.Sleep(5000000000) // 5 second expressed in nanoseconds
		err = http.ListenAndServe(port, nil)
		player.LastIP = "try"
	}
	player.Init()
}
