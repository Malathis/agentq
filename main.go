package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
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

type Spl struct {
	Name        string      `json:"spl"`
	ScheduleTag ScheduleTag `json:"scheduleTag"`
}

// let's declare a global PlaybackLogs array
// that we can populate to read from xp servers
var PlaybackLogs []PlaybackLog

var Queue []Spl

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "AgentQ homepage endpoint hit")
}

func receiveSpls(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: receiveSpls")

	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(responseData))

	// pushes the received spl into queue
	Ack := make([]Spl, 10)

	json.Unmarshal(responseData, &Ack)

	Queue = append(Queue, Ack...)
}

func receiveAckFromAgentQL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: receiveAckFromAgentQL")

	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(responseData))
}

func handleRequests() {
	http.HandleFunc("/", homePage)

	// add receiveSpls route and map it to our receiveSpls function like so
	http.HandleFunc("/receiveSpls", receiveSpls)

	// ack route and map it to our receiveAckFromAgentQ function like so
	http.HandleFunc("/ackFromAgentQL", receiveAckFromAgentQL)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func readAndPushLogs() {
	fmt.Println("Start of readAndPushLogs")

	// reads the logs from xp server
	PlaybackLogs = []PlaybackLog{
		{Id: 1, ScheduleTag: ScheduleTag{Id: 12500, Name: "Andhra GAP 112206603 AP"}, Screen: "Screen1", Cpl: "Cpl1"},
		{Id: 2, ScheduleTag: ScheduleTag{Id: 12501, Name: "Andhra GAP 112206603 AP AS"}, Screen: "Screen1", Cpl: "Cpl2"},
	}

	// push logs to agentql via receiveLogs api
	jsonValue, _ := json.Marshal(PlaybackLogs)
	response, err := http.Post("http://localhost:8083/receiveLogs", "application/json", bytes.NewBuffer(jsonValue))

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

func main() {
	// initialise agentq apis
	go handleRequests()

	// end-less process
	for {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Second * time.Duration(rand.Intn(10)))

		// sends acknowledgement of spl delivery
		if len(Queue) != 0 {
			jsonValue, _ := json.Marshal(Queue[0].ScheduleTag.Id)

			_, err := http.Post("http://localhost:8081/ackFromAgentQ", "application/json", bytes.NewBuffer(jsonValue))

			if err != nil {
				fmt.Print(err.Error())
			}
		}

		if len(Queue) > 1 {
			Queue = Queue[1:]
		} else {
			Queue = make([]Spl, 0)
		}

		// function to read logs from device server and send logs after a random wait to agentql
		readAndPushLogs()
	}
}
