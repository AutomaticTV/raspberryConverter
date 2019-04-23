package main

import (
       "log"
       "net/http"
       "fmt"
       "app/frontend"
			 "app/storage"
 )


 func DashboardHandler(w http.ResponseWriter, r *http.Request) {
        conditionsMap := map[string]interface{}{}
         if !IsLoggedIn(r) {
           http.Redirect(w, r, "/login", http.StatusFound)
           return
         }
         if r.FormValue("UpdateStatus") != "" {
           fmt.Println("Updating status:")
         } else if r.FormValue("UpdatePlayer") != "" {
           fmt.Println("Updating player:")
           fmt.Println("Video", r.FormValue("Video"))
           fmt.Println("AudioDecoding", r.FormValue("AudioDecoding"))
           fmt.Println("URL", r.FormValue("URL"))
           fmt.Println("Transport", r.FormValue("Transport"))
           fmt.Println("Buffer", r.FormValue("Buffer"))
           fmt.Println("Username", r.FormValue("Username"))
           fmt.Println("Password", r.FormValue("Password"))
           fmt.Println("Volume", r.FormValue("Volume"))
           fmt.Println("Autoplay", r.FormValue("Autoplay"))
           conditionsMap["PlayerError"] = true
           conditionsMap["PlayerErrorMessage"] = "CUSTOM VALIDATION MEASSAGE"
         } else if r.FormValue("UpdateNetwork") != "" {
           fmt.Println("Updating network:")
           fmt.Println("Mode", r.FormValue("Mode"))
           fmt.Println("IP", r.FormValue("IP"))
           fmt.Println("Gateway", r.FormValue("Gateway"))
           fmt.Println("Netmask", r.FormValue("Netmask"))
           fmt.Println("DNS1", r.FormValue("DNS1"))
           fmt.Println("DNS2", r.FormValue("DNS2"))
           conditionsMap["NetworkError"] = true
           conditionsMap["NetworkErrorMessage"] = "CUSTOM VALIDATION MEASSAGE"
         } else if r.FormValue("UpdatePassword") != "" {
           fmt.Println("Updating password:")
         }
         // Set up values to fill template
         fmt.Println(r.URL.Path)
         switch r.URL.Path {
         case "/dashboard/status":
           conditionsMap["Page"] = "Status"
           conditionsMap["CPU"] = 70
           conditionsMap["RAM"] = 24
           conditionsMap["URL"] = "https://some.example"
           conditionsMap["Video"] = "1080p59.94"
           conditionsMap["Status"] = "Running"
         case "/dashboard/":
           conditionsMap["Page"] = "Status"
           conditionsMap["CPU"] = 70
           conditionsMap["RAM"] = 24
           conditionsMap["URL"] = "https://some.example"
           conditionsMap["Video"] = "1080p59.94"
           conditionsMap["Status"] = "Running"
         case "/dashboard/player":
           conditionsMap["Page"] = "Player"
           conditionsMap["Video"] = "1080p59.94"
           conditionsMap["AudioDecoding"] = "Software"
           conditionsMap["URL"] = "https://some.example"
           conditionsMap["Transport"] = "HTTP"
           conditionsMap["Buffer"] = 300
           conditionsMap["Username"] = "admin"
           conditionsMap["Password"] = "admin"
           conditionsMap["Volume"] = 0
           conditionsMap["Autoplay"] = "No"
         case "/dashboard/network":
           conditionsMap["Page"] = "Network"
           conditionsMap["Mode"] = "Static"
           conditionsMap["IP"] = "192.168.1.57"
           conditionsMap["Netmask"] = "255.255.255.0"
           conditionsMap["Gateway"] = "192.168.1.1"
           conditionsMap["DNS1"] = "8.8.8.8"
           conditionsMap["DNS2"] = "8.8.8.8"
         case "/dashboard/password": conditionsMap["Page"] = "Password"
         default: conditionsMap["Page"] = "Status"
         }

         // Create html from template
         if err := template.Dashboard.Execute(w, conditionsMap); err != nil {
                 log.Println(err)
         }
 }

 func LoginHandler(w http.ResponseWriter, r *http.Request) {
         conditionsMap := map[string]interface{}{}

         // ALREADY LOGGED IN? Go to dashboard
         if IsLoggedIn(r) {
           http.Redirect(w, r, "/dashboard", http.StatusFound)
           return
         }

         // CHECK PASSWORD
         if r.FormValue("Login") != "" && r.FormValue("Username") != "" {
                 username := r.FormValue("Username")
                 password := r.FormValue("Password")
                 // PASSWORD CORRECT, GO TO DASHBOARD
								 if PasswordIsCorrect(username, password) {
                         conditionsMap["LoginError"] = false
                         err := Login(w, r, username)
                         if err != nil {
                           conditionsMap["LoginError"] = true
                         }

                         http.Redirect(w, r, "/dashboard", http.StatusFound)
                         return
                 } else { // PASSWORD ERROR, ADD ERROR MESSAGE TO TEMPLATE
                         conditionsMap["LoginError"] = true
                 }
         }

         // SERVE LOGIN PAGE
         if err := template.Login.Execute(w, conditionsMap); err != nil {
                 fmt.Println(err)
         }
 }

 func LogoutHandler(w http.ResponseWriter, r *http.Request) {
         err := Logout(w, r)
         if err != nil {
                 log.Println(err)
         }
         // redirec to login
         http.Redirect(w, r, "/login", http.StatusFound)
 }

 func main() {
         storage.InitStorage()
         fmt.Println("Server starting, point your browser to localhost:80/login to start")
				 // ENDPOINTS
         http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))
         http.HandleFunc("/", LoginHandler)
         http.HandleFunc("/login", LoginHandler)
				 http.HandleFunc("/logout", LogoutHandler)
         http.HandleFunc("/dashboard/", DashboardHandler)
         http.ListenAndServe(":80", nil)
 }
