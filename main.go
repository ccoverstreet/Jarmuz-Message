package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Config struct {
	BotID string `json:"botID"`
}

func main() {
	log.Println("Jarmuz-Message starting...")

	JablkoCorePort := os.Getenv("JABLKO_CORE_PORT")
	JMODPort := os.Getenv("JABLKO_MOD_PORT")
	JMODKey := os.Getenv("JABLKO_MOD_KEY")
	JMODConfig := os.Getenv("JABLKO_MOD_CONFIG")

	log.Println(JablkoCorePort, JMODPort, JMODKey)

	var conf Config

	if len(JMODConfig) < 4 {
		conf = Config{"PLACEHOLDER"}
	} else {
		err := json.Unmarshal([]byte(JMODConfig), &conf)
		if err != nil {
			log.Println("ERROR: Unable to unmarshal config", err)
			panic(err)
		}
	}

	context := CreateAppContext(Config{conf.BotID})

	router := mux.NewRouter()
	router.HandleFunc("/webComponent", handleWebComponent)
	router.HandleFunc("/instanceData", handleInstanceData)
	router.HandleFunc("/jmod/sendMessage", wrapHandle(handleSendMessage, context))

	http.ListenAndServe("127.0.0.1:"+JMODPort, router)
}

func ParseJSONBody(body io.ReadCloser, dest interface{}) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil
	}
	return json.Unmarshal(b, dest)
}

func wrapHandle(handle func(*AppContext, http.ResponseWriter, *http.Request), ctx *AppContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handle(ctx, w, r)
	}
}

func handleWebComponent(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "class{}")
}

func handleInstanceData(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "[]")
}

func handleSendMessage(ctx *AppContext, w http.ResponseWriter, r *http.Request) {
	var data struct {
		Message string `json:"message"`
	}

	err := ParseJSONBody(r.Body, &data)
	log.Println(data)
	log.Println(err)

	ctx.SendMessage(data.Message)
}
