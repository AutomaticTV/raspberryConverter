package template

import (
  "html/template"
)

var basepath = "frontend/templates/"
var Dashboard = template.Must(template.ParseFiles(
  basepath + "header.html",
  basepath + "dashboard.html",
  basepath + "status.html",
  basepath + "password.html",
  basepath + "network.html",
  basepath + "player.html",
))

var Login = template.Must(template.ParseFiles(
  basepath + "header.html",
  basepath + "login.html",
))
