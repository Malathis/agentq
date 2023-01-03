package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ScheduleTag struct {
	Id   int    `json:"id"`
	Name string `json:"schedule"`
}

type PlaybackLog struct {
	Id          int         `json:"id"`
	ScheduleTag ScheduleTag `json:"scheduleTag"`
	Screen      string      `json:"schedule"`
	Cpl         string      `json:"orderId"`
}

// let's declare a global Logs array
// that we can populate to read from xp servers
var PlaybackLogs []PlaybackLog

func getSpls(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getSpls")

	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(responseData))
}

func readAndPushLogs() {
	fmt.Println("Start of readAndPushLogs")

	// reads the logs from xp server
	PlaybackLogs = []PlaybackLog{
		PlaybackLog{Id: 1, ScheduleTag: ScheduleTag{Id: 12500, Name: "Andhra GAP 112206603 AP"}, Screen: "Screen1", Cpl: "Cpl1"},
		PlaybackLog{Id: 2, ScheduleTag: ScheduleTag{Id: 12501, Name: "Andhra GAP 112206603 AP AS"}, Screen: "Screen2", Cpl: "Cpl2"},
	}

	jsonValue, _ := json.Marshal(PlaybackLogs)
	response, err := http.Post("http://localhost:8083/saveLogs", "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(responseData))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "AgentQ homepage endpoint hit")
}

func handleRequests() {
	http.HandleFunc("/", homePage)

	// add our spls route and map it to our
	// getSpls function like so
	http.HandleFunc("/getSpls", getSpls)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func main() {
	go readAndPushLogs()
	handleRequests()
}
