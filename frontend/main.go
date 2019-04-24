package template

import (
  "html/template"
  "github.com/gobuffalo/packr"
)

var box = packr.NewBox("./templates")

var Dashboard = template.Must(template.New("Dashboard").Parse(
  getTemplate("header.html") +
  getTemplate("dashboard.html") +
  getTemplate("status.html") +
  getTemplate("password.html") +
  getTemplate("network.html") +
  getTemplate("player.html"),
))

var Login = template.Must(template.New("Login").Parse(
  getTemplate("header.html") +
  getTemplate("login.html") ,
))

func getTemplate(name string) string {
  temp, _ := box.FindString(name)
  return temp
}
