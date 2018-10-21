package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

// TWIMLResponse is an TWIML XML response
type TWIMLResponse struct {
	Message string `xml:"Message"`
}

// TWIMLRequest represents a form posting from Twilio
type TWIMLRequest struct {
	ToCountry     string
	ToState       string
	SmsMessageSid string
	NumMedia      string
	ToCity        string
	FromZip       string
	SmsSid        string
	FromState     string
	SmsStatus     string
	FromCity      string
	Body          string
	FromCountry   string
	To            string
	ToZip         string
	MessageSid    string
	AccountSid    string
	From          string
	ApiVersion    string
}

var tasklist []string

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/todo", incomingHandler).Methods("POST")
	http.Handle("/", loghttp(r))
	if err := http.ListenAndServe(":12001", nil); err != nil {
		log.Fatal(err)
	}
}

func incomingHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Println("Can't parse form")
	}
	decoder := schema.NewDecoder()
	var in TWIMLRequest
	decoder.Decode(in, r.Form)
	log.Println(r.Form)

	echo := "Sorry, please use: add [task], remove [id], or list"
	body := r.Form["Body"][0]
	// clean up this if-then, to maybe a switch
	if strings.Contains(body, "add ") {
		task := strings.TrimSpace(strings.TrimLeft(body, `add `))
		echo = fmt.Sprintf("adding '%s'", task)
		tasklist = append(tasklist, task)
	} else if strings.Contains(body, "list") {
		if len(tasklist) == 0 {
			echo = "No tasks, please use: add [task]"
		} else {
			echo = ""
			for k, v := range tasklist {
				echo = fmt.Sprintf("%s%v. %s\n", echo, k+1, v)
			}
		}
	} else if strings.Contains(body, "remove ") {
		// lots of missing guards here, out of bound checks
		taskid, err := strconv.Atoi(strings.TrimSpace(strings.TrimLeft(body, "remove ")))
		if err != nil {
			echo = "Task list item to remove must be a number"
		} else {
			taskid = taskid - 1
			tasklist = append(tasklist[:taskid], tasklist[taskid+1:]...)
			echo = fmt.Sprintf("i'll remove item %v", taskid+1)
		}
	}

	message := fmt.Sprintf("<Response><Message>%s</Message></Response>", echo)

	w.Header().Set("content-type", "application/xml")
	w.Write([]byte(message))
}

// loghttp just logs the path and headers
func loghttp(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path, r.Header)
		h.ServeHTTP(w, r)
	})
}
