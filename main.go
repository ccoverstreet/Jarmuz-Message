package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"github.com/ccoverstreet/jablkodev"
	"github.com/gorilla/mux"
)

type Config struct {
	BotID string `json:"botID"`
}

func main() {
	log.Println("Jarmuz-Message starting...")

	// Load environment variables
	err := jablkodev.LoadJablkoEnv()
	if err != nil {
		panic("Jablko environment variables aren't set. Make sure to run this as a JMOD or set up a fake environment")
	}

	var conf Config

	confStr := jablkodev.GetJablkoModConfig()

	if len(confStr) < 4 {
		conf = Config{"PLACEHOLDER"}
	} else {
		err := json.Unmarshal([]byte(confStr), &conf)
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

	http.ListenAndServe("127.0.0.1:"+strconv.Itoa(jablkodev.GetJablkoModPort()), router)
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

func httpErrorHandler(err error, msg string, w http.ResponseWriter) {
	_, filename, line, _ := runtime.Caller(1)
	log.Printf("ERROR: %s: %d\n\t%v\n\t%s", filename, line, err, msg)
	fmt.Fprintf(w, `{"err": "%s"}`, msg)
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
	if err != nil {
		httpErrorHandler(err, "Unable to parse JSON body", w)
		return
	}

	err = ctx.SendMessage(data.Message)
	if err != nil {
		httpErrorHandler(err, "Unable to send GroupMe message", w)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "successful"}`)
}
