package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
)

const (
	SbAuthServiceEndpointBaseUrl = "http://127.0.0.1:5000/"
)

var upgrader = websocket.Upgrader{}

type Env struct {
	db *gorm.DB
}

type FFA struct {
	// NOTE: not sure if we need a db struct in here
	Env
	Players        Players
	RoundIsOngoing bool
}

func (f *FFA) ffa(w http.ResponseWriter, r *http.Request) {
	// Handles the FFA Room
	// The First Message needs to contain a valid auth token usually obtained from the SB-Auth-Service
	// further msges should be a json object representing a player struct representing the client
	// will return a Error Response if authentication fails
	// will return players every iteration after auth was valid

	var conn, _ = upgrader.Upgrade(w, r, nil)

	go func(conn *websocket.Conn) {
		auth := false
		playerJoined := false

		for {
			// check if client is authenticated
			// the idea is to have guest tokens for guest players i dont think we should allow clients to connect without any kind of auth
			if !auth {
				// check if token is valid
				msg := Token{}
				err := conn.ReadJSON(&msg)
				if err != nil {
					fmt.Println("this panics 1")
				}
				token := msg.Token

				if !isAuth(token) {
					// token was invalid
					err = conn.WriteJSON(&ErrorResponse{Msg: "Authentication failed!"})
					if err != nil {
						fmt.Println("this panics 4")
						fmt.Println(err)
					}
					continue
				}

				// token is valid setting auth to true
				auth = true
			}

			// sending client the connected players
			conn.WriteJSON(f.Players)

			// adding the player to the ffa players if the round is not ongoing and the client has not already joined
			if !f.RoundIsOngoing || !playerJoined {

				player := Player{Name: "name of the player we get from sb auth service which is not yet implemented", Active: true, Mutex: &sync.Mutex{}}
				*f.Players = append(*f.Players, player)
				playerJoined = true
			}

			// reading client msgs
			updatedPlayer := Player{}
			conn.ReadJSON(updatedPlayer)

		}
	}(conn)

}

func main() {

	// enviorment structs
	ffa := FFA{}

	// Router /Routes
	r := mux.NewRouter()

	r.HandleFunc("/ffa", ffa.ffa)

	// TESTING
	// Index page for testing and stuff
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	r.HandleFunc("/testing", func(w http.ResponseWriter, r *http.Request) {
		var conn, _ = upgrader.Upgrade(w, r, nil)
		go func(conn *websocket.Conn) {
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					conn.Close()
				}
				fmt.Println(string(msg))
			}
		}(conn)

		go func(conn *websocket.Conn) {
			ch := time.Tick(5 * time.Second)

			for range ch {
				conn.WriteJSON(myStruct{
					Username:  "mvansickle",
					FirstName: "Michael",
					LastName:  "Van Sickle",
				})
			}
		}(conn)
	})

	// Server
	log.Fatal(http.ListenAndServe(":3000", r))
}

type myStruct struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
