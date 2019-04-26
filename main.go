package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"raspberryConverter/frontend"
	"raspberryConverter/models"
	"raspberryConverter/services/auth"
	"raspberryConverter/services/monitor"
	"raspberryConverter/services/network"
	"raspberryConverter/services/player"
	"strconv"

	"github.com/gobuffalo/packr"
)

// port is a constant that describe the port listened by the server, example: http://localhost:5555
const port = ":5555"

func handler(w http.ResponseWriter, r *http.Request) {
	// IF NOT LOGGED IN, GO TO LOGIN PAGE
	if !auth.IsLoggedIn(r) && r.URL.Path != "/login" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	isUpdate, requestError := requestHandler(w, r)

	// RETURN THE APROPIATE FRONTEND PAGE ACCORDING TO THE PATH
	fmt.Println(r.URL.Path)
	switch r.URL.Path {
	case "/login":
		fmt.Println("login1")
		if requestError == nil && isUpdate {
			fmt.Println("login2")
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		}
		fmt.Println("login3")
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
			fmt.Println(err)
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

// requestHandler is a function that checks if the request includes an expected form and if so it gets handled.
// Returns true and a possible error if r includes a form otherwisr (false, nil)
func requestHandler(w http.ResponseWriter, r *http.Request) (bool, error) {
	// LOGIN
	if r.FormValue("Login") != "" && r.FormValue("Username") != "" {
		username := r.FormValue("Username")
		password := r.FormValue("Password")
		if auth.PasswordIsCorrect(username, password) {
			return true, auth.Login(w, r, username)
		}
		// PASSWORD ERROR, ADD ERROR MESSAGE TO TEMPLATE
		return true, errors.New("Wrong username or password")
	}
	// PLAY
	if r.FormValue("Play") != "" {
		return true, player.Play()
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
				Transport:     r.FormValue("Transport"),
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

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	err := auth.Logout(w, r)
	if err != nil {
		log.Println(err)
	}
	// redirec to login
	http.Redirect(w, r, "/login", http.StatusFound)
}

func main() {
	staticFiles := packr.NewBox("frontend/static")
	player.Init()
	auth.Init()
	fmt.Println("Server starting, point your browser to localhost" + port + "/login to start")
	// ENDPOINTS
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticFiles)))
	http.HandleFunc("/", handler)
	http.HandleFunc("/login", handler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/dashboard/", handler)
	http.ListenAndServe(port, nil)
}
