package frontend

import (
	"fmt"
	"html/template"
	"net/http"
	"raspberryConverter/models"

	"github.com/gobuffalo/packr"
)

type page struct {
	Page           string
	Error          bool
	ErrorMessage   string
	Success        bool
	SuccessMessage string
	Data           interface{}
}

var box = packr.NewBox("./templates")

//dashboardTemplate is a template that produces the HTML code for the main page
var dashboardTemplate *template.Template

//login is a template that produces the HTML code for the login page
var loginTemplate *template.Template

// Init loads the static files to be rendered
func Init() {
	loginTemplate = template.Must(template.New("Login").Parse(getTemplate("header.html") +
		getTemplate("login.html"),
	))
	dashboardTemplate = template.Must(template.New("Dashboard").Parse(getTemplate("header.html") +
		getTemplate("dashboard.html") +
		getTemplate("status.html") +
		getTemplate("password.html") +
		getTemplate("network.html") +
		getTemplate("player.html"),
	))
}

// Login is a function that render the status.html template and send it through w
func Login(isUpdate bool, err error, w http.ResponseWriter) error {
	p := page{}
	setSuccessOrError(
		&p, err, isUpdate,
		"Login succesful.")
	return loginTemplate.Execute(w, p)
}

// Status is a function that render the status.html template and send it through w
func Status(isUpdate bool, s models.Status, err error, w http.ResponseWriter) error {
	p := page{Page: "Status", Data: s}
	setSuccessOrError(
		&p, err, isUpdate,
		"Command sended succesfuly.")
	return dashboardTemplate.Execute(w, p)
}

// Network is a function that render the network.html template and send it through w
func Network(isUpdate bool, nc models.NetworkConfig, err error, w http.ResponseWriter) error {
	p := page{Page: "Network", Data: nc}
	setSuccessOrError(
		&p, err, isUpdate,
		"Network setings updated successfuly. You may have to reboot the device to make the changes efective and see it reflected on the page information.")
	return dashboardTemplate.Execute(w, p)
}

// Player is a function that render the player.html template and send it through w
func Player(isUpdate bool, pc models.PlayerConfig, err error, w http.ResponseWriter) error {
	p := page{Page: "Player", Data: pc}
	setSuccessOrError(
		&p, err, isUpdate,
		"Player configuration updated! Restart the player at the page STATUS for the changes to take effect.")
	return dashboardTemplate.Execute(w, p)
}

// Password is a function that render the password.html template and send it through w
func Password(isUpdate bool, err error, w http.ResponseWriter) error {
	p := page{Page: "Password"}
	setSuccessOrError(
		&p, err, isUpdate,
		"Password updated succesfuly.")
	return dashboardTemplate.Execute(w, p)
}

func setSuccessOrError(p *page, err error, isUpdate bool, successMessage string) {
	if err != nil {
		p.Error = true
		p.ErrorMessage = err.Error()
	} else if isUpdate {
		p.Success = true
		p.SuccessMessage = successMessage
	}
}

func getTemplate(name string) string {
	temp, err := box.Find(name)
	if err != nil {
		fmt.Println("error geting from tha box", err)
	}
	return string(temp)
}
