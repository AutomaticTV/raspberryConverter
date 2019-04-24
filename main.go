package main

import (
       "log"
       "net/http"
       "fmt"
       "reflect"
       "raspberryConverter/services"
       "raspberryConverter/frontend"
			 "raspberryConverter/storage"
       "github.com/gobuffalo/packr"
 )


 func DashboardHandler(w http.ResponseWriter, r *http.Request) {
        conditions := map[string]interface{}{}
        // IF NOT LOGGED IN, GO TO LOGIN PAGE
         if !services.IsLoggedIn(r) {
           http.Redirect(w, r, "/login", http.StatusFound)
           return
         }
         // CHECK IF THE REQUEST INCLUDES A FORM
         if r.FormValue("UpdateStatus") != "" {

        // UPDATE PLAYER CONFIG
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
           conditions["PlayerError"] = true
           conditions["PlayerErrorMessage"] = "CUSTOM VALIDATION MEASSAGE"
        // UPDATE NETWORK
         } else if r.FormValue("UpdateNetwork") != "" {
           fmt.Println("Updating network:")
           fmt.Println("Mode", r.FormValue("Mode"))
           fmt.Println("IP", r.FormValue("IP"))
           fmt.Println("Gateway", r.FormValue("Gateway"))
           fmt.Println("Netmask", r.FormValue("Netmask"))
           fmt.Println("DNS1", r.FormValue("DNS1"))
           fmt.Println("DNS2", r.FormValue("DNS2"))
           conditions["NetworkError"] = true
           conditions["NetworkErrorMessage"] = "CUSTOM VALIDATION MEASSAGE"
        // UPDATE PASSWORD
         } else if r.FormValue("UpdatePassword") != "" {
           err := services.UpdatePassword(
             r.FormValue("Username"),
             r.FormValue("OldPassword"),
             r.FormValue("NewPassword"),
             r.FormValue("RePassword"),
           )
           if (err != nil){
             fmt.Println(err)
             AddError(err.Error(), conditions)
           } else {
             AddSuccess("Password updated successfuly", conditions)
           }
         }

         // SET CONDITIONS TO FILL TEMPLATE IF REQUIRED AFTER FORM
         switch r.URL.Path {
         case "/dashboard/status":
           SetConditions(storage.GetStatus, "Status", conditions, "Unable to retrieve last status.")
         case "/dashboard/player":
           SetConditions(storage.GetPlayer, "Player", conditions, "Unable to retrieve last player setings.")
         case "/dashboard/network":
           SetConditions(services.GetNetworkConfig, "Network", conditions, "Unable to retrieve last network setings.")
         case "/dashboard/password":
           conditions["Page"] = "Password"
         default:
           http.Redirect(w, r, "/dashboard/status", http.StatusFound)
           return
         }

         // CREATE HTML FROM TEMPLATE USING CONDITIONS
         if err := template.Dashboard.Execute(w, conditions); err != nil {
                 log.Println(err)
         }
 }

 func SetConditions(get func() (interface{}, error), page string, conditions map[string]interface{}, errorMessage string) {
   values, err := get()
   if err != nil {
     fmt.Println(err)
     AddError(errorMessage, conditions)
   } else {
     AddToConditionsMap(&values, conditions)
   }
   conditions["Page"] = page
 }

 func AddError(message string, conditions map[string]interface{}) {
   // test!!!!
   conditions["Error"] = true
   if _, ok := conditions["ErrorMessage"]; ok {
        conditions["ErrorMessage"] = conditions["ErrorMessage"].(string) + "; " + message
    } else {
      conditions["ErrorMessage"] = message
    }
 }

  func AddSuccess(message string, conditions map[string]interface{}) {
    conditions["Success"] = true
    conditions["SuccessMessage"] = message
  }

 func AddToConditionsMap(record *interface{}, conditions map[string]interface{}) {
    values := reflect.ValueOf(*record).Elem()
    names := values.Type()
    for i := 0; i < values.NumField(); i++ {
      name := names.Field(i).Name
      if name != "Model" {
        conditions[name] = values.Field(i).Interface()
      }
    }
 }

 func LoginHandler(w http.ResponseWriter, r *http.Request) {
         conditions := map[string]interface{}{}

         // ALREADY LOGGED IN? Go to dashboard
         if services.IsLoggedIn(r) {
           http.Redirect(w, r, "/dashboard", http.StatusFound)
           return
         }

         // CHECK PASSWORD
         if r.FormValue("Login") != "" && r.FormValue("Username") != "" {
                 username := r.FormValue("Username")
                 password := r.FormValue("Password")
                 // PASSWORD CORRECT, GO TO DASHBOARD
								 if services.PasswordIsCorrect(username, password) {
                         conditions["LoginError"] = false
                         err := services.Login(w, r, username)
                         if err != nil {
                           conditions["LoginError"] = true
                         }

                         http.Redirect(w, r, "/dashboard", http.StatusFound)
                         return
                 } else { // PASSWORD ERROR, ADD ERROR MESSAGE TO TEMPLATE
                         conditions["LoginError"] = true
                 }
         }

         // SERVE LOGIN PAGE
         if err := template.Login.Execute(w, conditions); err != nil {
                 fmt.Println(err)
         }
 }

 func LogoutHandler(w http.ResponseWriter, r *http.Request) {
         err := services.Logout(w, r)
         if err != nil {
                 log.Println(err)
         }
         // redirec to login
         http.Redirect(w, r, "/login", http.StatusFound)
 }

 func main() {
         staticFiles := packr.NewBox("frontend/static")
         storage.InitStorage()
         fmt.Println("Server starting, point your browser to localhost:80/login to start")
				 // ENDPOINTS
         http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticFiles)))
         http.HandleFunc("/", LoginHandler)
         http.HandleFunc("/login", LoginHandler)
				 http.HandleFunc("/logout", LogoutHandler)
         http.HandleFunc("/dashboard/", DashboardHandler)
         http.ListenAndServe(":5555", nil)
 }
