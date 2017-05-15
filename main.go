package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"net/http"
	"flag"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 1 * time.Second

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

var (
	addr      = flag.String("addr", ":8080", "http service address")
	world       World
)


func pumpStdout(ws *websocket.Conn, done chan struct{}) {
	defer func() {
	}()

	for {
		p, _, err := world.RunGeneration()
		if err != nil {
			ws.Close()
			break
		}
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.TextMessage, p); err != nil {
			ws.Close()
			break
		}
		time.Sleep(filePeriod)
	}

	close(done)

	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
}

func ping(ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
			}
		case <-done:
			return
		}
	}
}

func internalError(ws *websocket.Conn, msg string, err error) {
	log.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

var upgrader = websocket.Upgrader{}

func serveWs(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer ws.Close()

	stdoutDone := make(chan struct{})
	go pumpStdout(ws, stdoutDone)
	go ping(ws, stdoutDone)

	select {
	case <-stdoutDone:
	case <-time.After(time.Second):
		<-stdoutDone
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	log.Println("serveHome")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "home.html")
}

func main() {
	f, err := os.OpenFile("gameoflife.log", os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	filename := os.Args[1]
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error case ", err)
		os.Exit(1)
	}

	log.Println(string(data))

	world.readFromJson(data)

	log.Println("Width = ", world.Width)
	log.Println("Height = ", world.Height)

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

