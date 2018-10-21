package main

import (
	"fmt"
	"log"
	"net/http"

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

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/sms", incomingHandler)
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

	incountry := r.Form["FromCountry"][0]
	country := fmt.Sprintf("Hi! It looks like your phone number was born in %s", incountry)
	w.Header().Set("content-type", "application/xml")
	w.Write([]byte(fmt.Sprintf("<Response><Message>%s</Message></Response>", country)))
}

// loghttp just logs the path and headers
func loghttp(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path, r.Header)
		h.ServeHTTP(w, r)
	})
}
