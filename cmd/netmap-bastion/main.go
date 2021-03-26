package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"netmap-bastion/internal/domain"
	"netmap-bastion/internal/repositories/graphrepo"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]string) // connected clients

var inboundChannel = make(chan domain.TrafficGraph)

// Configure the upgrader
var upgrader = websocket.Upgrader{}

// database
var graphRepository *graphrepo.PostgresRepo

// database connection details
var (
	dbhost     = "localhost"
	dbport     = 5432
	dbuser     = "postgres"
	dbpassword = "password"
	dbname     = "netmap"
)

func init() {

	flag.StringVar(&dbhost, "dbhost", "localhost", "DB host; default localhost")
	flag.IntVar(&dbport, "dbport", 5432, "DB port; default 5432")
	flag.StringVar(&dbuser, "dbuser", "postgres", "DB user; default")
	flag.StringVar(&dbpassword, "dbpassword", "password", "DB password; default")

	flag.Parse()
}

func main() {
	var err error
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname)
	graphRepository, err = graphrepo.NewPostgresRepo(psqlconn)
	if err != nil {
		panic(err)
	}

	// Create a simple file server
	fs := http.FileServer(http.Dir("../../public"))
	http.Handle("/", fs)

	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// Start listening for incoming chat messages
	go handleInbound()

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8000")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()
	defer fmt.Println("<<< disconnected: " + ws.LocalAddr().String())

	// Register our new client
	clients[ws] = ws.LocalAddr().String()
	fmt.Println(">>> new connection: " + clients[ws])

	for {
		var graph domain.TrafficGraph
		err := ws.ReadJSON(&graph)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// fmt.Println(graph)
		fmt.Println("sending payload for processing...")
		// Send payload for processing
		inboundChannel <- graph
	}
}

func handleInbound() {
	for {
		// Grab the next message from the inbound channel
		graph := <-inboundChannel
		fmt.Println("received payload, inserting in db...")
		err := graphRepository.Insert("testing", graph)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
