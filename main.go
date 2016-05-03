package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

const URL = "https://api.telegram.org/bot" + TOKEN + "/"

type Incoming struct {
	UpdateID int64 `json:"update_id"`
	Message  struct {
		MessageID int64 `json:"message_id"`
		From      struct {
			ID        int64
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string
			Type      string
		}
		Chat struct {
			ID        int64
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string
			Type      string
		}
		Date     int64
		Text     string
		Entities []struct {
			Type   string
			Offset int64
			Length int64
		}
	}
}

type Informoj struct {
	Celvorto string
	Diveno   []rune
	Vicoj    int
}

func iksigi(in string) string {
	conversion := []struct {
		from string
		to   string
	}{
		{"cx", "ĉ"},
		{"gx", "ĝ"},
		{"hx", "ĥ"},
		{"jx", "ĵ"},
		{"sx", "ŝ"},
		{"ux", "ŭ"},
		{"Cx", "Ĉ"},
		{"Gx", "Ĝ"},
		{"Hx", "Ĥ"},
		{"Jx", "Ĵ"},
		{"Sx", "Ŝ"},
		{"Ux", "Ŭ"},
		{"CX", "Ĉ"},
		{"GX", "Ĝ"},
		{"HX", "Ĥ"},
		{"JX", "Ĵ"},
		{"SX", "Ŝ"},
		{"UX", "Ŭ"},
	}
	for _, c := range conversion {
		in = strings.Replace(in, c.from, c.to, -1)
	}
	return in
}

func telegram(w http.ResponseWriter, r *http.Request) {
	var mymessage string

	c := appengine.NewContext(r)
	client := &http.Client{
		Transport: &urlfetch.Transport{
			Context: c,
		},
	}
	request, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf(c, "%v", err)
	}
	r.Body.Close()
	var Output Incoming
	err = json.Unmarshal(request, &Output)
	if err != nil {
		log.Errorf(c, "%v", err)
	}
	log.Debugf(c, "%v", Output)
	command := regexp.MustCompile("/[^ @]*").FindString(Output.Message.Text)
	text := regexp.MustCompile("^/[^ ]* ").ReplaceAllString(Output.Message.Text, "")
	switch command {
	case "/start":
		mymessage = "Saluton! Ĉi tiu roboto ankoraŭ ne pretas. Ĉi tie vi baldaŭ povos lerni Esperanton! Sendu komentojn pri la roboto al @lapingvino"
	case "/echo":
		mymessage = regexp.MustCompile(`(["\\])`).ReplaceAllString(text, `\$1`)
	}
	client.Post(URL+"sendMessage", "application/json", strings.NewReader(fmt.Sprintf("{\"chat_id\": %v, \"text\": \"%v\"}", Output.Message.Chat.ID, mymessage)))
}

func init() {
	http.HandleFunc("/"+SECRETLINK, telegram)
}
