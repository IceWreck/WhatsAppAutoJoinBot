package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Rhymen/go-whatsapp"
)

var wacon *whatsapp.Conn

func main() {

	setup := func() {
		wac, err := whatsapp.NewConnWithOptions(&whatsapp.Options{
			// timeout
			Timeout: 20 * time.Second,
			// set custom client name
			ShortClientName: "WhatsApp AutoJoin",
			LongClientName:  "WhatsApp AutoJoin",
		})
		if err != nil {
			newLog("Cant create connection", err)
			panic(err)
		}

		wac.SetClientVersion(2, 2123, 7)
		wac.AddHandler(&waHandler{wac, uint64(time.Now().Unix())})

		// run this in background
		go func() {
			if err = login(wac); err != nil {
				fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
			}
		}()

		// set global wacon = connnection
		// shoddy work, i knowbut whatever
		wacon = wac
	}

	setup()

	// ticker to recheck and relogin in case it breaks connection
	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				fmt.Println("Login Check tick at", t)
				if !wacon.GetLoggedIn() {
					newLog("Logged Out, retry login")
					setup()
					// check if relogin worked
					time.Sleep(2 * time.Minute)
					if !wacon.GetLoggedIn() {
						// send a notification about logout
						// comment out if you dont want notifications

						// _, _ = http.PostForm("https://hookmsg.example.com/hooks/matrix/script-reports", url.Values{
						// 	"secret":  {"HOOKMSG-SECRET"},
						// 	"content": {"**WhatsApp AutoJoin Bot:** Can't Login"},
						// })
					}
				}
			}
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ts, _ := template.ParseFiles("./templates/index.tmpl")
		ts.Execute(w, struct {
			LoggedIn bool
		}{
			wacon.GetLoggedIn(),
			// false,
		})
	})

	// start a new login
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		setup()
		// tell user to wait
		ts, _ := template.ParseFiles("./templates/login.tmpl")
		ts.Execute(w, nil)
	})

	http.HandleFunc("/qr", func(w http.ResponseWriter, r *http.Request) {
		ts, _ := template.ParseFiles("./templates/qr.tmpl")
		ts.Execute(w, struct {
			QRString string
			NeedQR   bool
		}{
			getState("qr"),
			func() bool {
				needed := false
				if time.Since(lastUpdated) < time.Second*20 {
					needed = true
				}
				if len(getState("qr")) == 0 {
					needed = false
				}
				return needed
			}(),
		})
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		err := wacon.Logout()
		if err != nil {
			w.Write([]byte("Done"))
		} else {
			newLog("error during logout", err)
			w.Write([]byte("Error"))
		}
	})

	// view logs
	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		ts, _ := template.ParseFiles("./templates/logs.tmpl")
		ts.Execute(w, struct {
			Logs []string
		}{
			getLogs(),
		})

	})

	// delete saved files and reset everything
	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			removeSession()
			// exit program, systemd should handle restart
			os.Exit(0)
		}()

		w.Write([]byte("Resetting and Restarting"))
	})

	log.Fatal(http.ListenAndServe(":8755", nil))
}
