package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"fmt"
	"log"
	"net/http"
)

var wrongMsg string
var tokens map[string]string
var authorizationKey string

func main() {
	authorizationKey = os.Args[1]
	tokens = make(map[string]string)
	http.HandleFunc("/addToken/", handleAddToken)
	http.HandleFunc("/sendMessage/", handleSendMessage)
	log.Fatal(http.ListenAndServe(":65000", nil))
}

func sendMessageToFirebase(token string, msg string) {
	var jsonStr = []byte(`{
		"data": {
			"msg": "`+msg+`"
		},
		"to" : "`+token+`"
	}`)
	req, err := http.NewRequest(
		"POST",
		"https://fcm.googleapis.com/fcm/send",
		bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key="+authorizationKey)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func handleAddToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST requests only", http.StatusMethodNotAllowed)
		return
	} else {
		var token string;
		token = r.FormValue("id")
		_,ok := tokens[token]
		if (ok == false) {
			tokens[token] = token
			log.Println("addToken: " + token)
		} else {
			log.Println("Token still exists")
		}
	}
	fmt.Fprintf(w, "Hallo Client")
}

func handleSendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST requests only", http.StatusMethodNotAllowed)
		return
	} else {
		var toToken string;
		var msg string;
		toToken = r.FormValue("to")
		msg = r.FormValue("msg")
		log.Println("SendMessage: " + toToken + " -> " + msg)
		_,ok := tokens[toToken]
		if (ok == false) {
			log.Println("Token not exists")
		} else {
			sendMessageToFirebase(toToken, msg)
		}
	}
	fmt.Fprintf(w, "Hallo Client")
}
