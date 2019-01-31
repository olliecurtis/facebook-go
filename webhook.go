package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	PAGE_TOKEN = "EAALD2ojScyEBAHYKb9myrpnkd4dtINBjelRHATU4Tvu2QjC2Ag92HsOtZBOwRo6P9bJZBZBdZCz5oWdQacA6qfQElQtBqaFbTsrGqyHpcyVlZAjJ4EVBHnO3JQe5sQQVn5e5EUZBw7Vg0aE3FSheM37qjVqalVwcQ7plhZAMM2lLwZDZD"
	AUTH_TOKEN = ""
)

func main() {
	http.HandleFunc("/webhook", Verify)
	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

type MessengerInput struct {
	Entry []struct {
		Time      uint64 `json:"time,omitempty"`
		Messaging []struct {
			Sender struct {
				Id string `json:"id"`
			} `json:"sender,omitempty"`
			Recipient struct {
				Id string `json:"id"`
			} `json:"recipient,omitempty"`
			Timestamp uint64 `json:"timestamp,omitempty"`
			Message   *struct {
				Mid  string `json:"mid,omitempty"`
				Seq  uint64 `json:"seq,omitempty"`
				Text string `json:"text"`
			} `json:"message,omitempty"`
		} `json:"messaging"`
	}
}

func Verify(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		challenge := r.URL.Query().Get("hub.challenge")
		verify_token := r.URL.Query().Get("hub.verify_token")
		if len(verify_token) > 0 && len(challenge) > 0 && verify_token == "developers-are-great" {
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, challenge)
			return
		}
	} else if r.Method == "POST" {
		defer r.Body.Close()

		input := new(MessengerInput)
		if err := json.NewDecoder(r.Body).Decode(input); err == nil {

			//lets swap sender and recipient
			reply := input.Entry[0].Messaging[0]
			reply.Sender, reply.Recipient = reply.Recipient, reply.Sender
			fmt.Println(input.Entry[0].Messaging[0].Message.Text)

			reply.Message.Text = input.Entry[0].Messaging[0].Message.Text
			reply.Message.Seq = 0 //these fields are not used so remove them with omit empty
			reply.Message.Mid = ""

			b, _ := json.Marshal(reply)
			url := fmt.Sprintf("https://graph.facebook.com/v2.6/me/messages?access_token=%s", PAGE_TOKEN)
			http.Post(url, "application/json", bytes.NewReader(b))
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(400)
	fmt.Fprintf(w, "Bad Request")
}
