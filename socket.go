package main

import (
	"log"

	"github.com/googollee/go-socket.io"
)

func socketConnectionHandler(s socketio.Socket) {
	s.On("open", func(con [2]string) {
		user := con[0]
		room := con[1]
		s.Join(con[1])
		log.Println(user + " connected to " + room)
		s.BroadcastTo(room, "message", user+" connected to "+room)
		s.Emit("message", "You have connected to "+room)
	})
	s.On("message", func(msg [2]string) {
		room := msg[0]
		message := msg[1]
		log.Printf("%+v\n", msg)
		s.Emit("message", message)
		s.BroadcastTo(room, "message", message)
	})
	s.On("close", func(con [2]string) {
		user := con[0]
		room := con[1]
		log.Println(user + " has left " + room)
		s.BroadcastTo(room, "close", user+" has left "+room)
	})
	s.On("disconnection", func() {
		log.Println("a user has disconnected")
	})
	s.On("gameDeleted", func(){
		log.Println("a user has deleted a game")
		s.BroadcastTo("lobby", "gameDeleted", "")
	})
	s.On("updateGameList", func(){
		log.Println("a user has added a game")
		s.BroadcastTo("lobby", "updateGameList", "")
	})
	s.On("gameConsole", func(msg [2]string){
		room := msg[0]
		message := msg[1]
		log.Println(room + " : " + message)
		s.BroadcastTo(room, "gameConsole", message)
	})
}

func socketErorrHandler(s socketio.Socket, err error) {
	log.Println("error: ", err)
}
