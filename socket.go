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
	s.On("updateRecord", func(){
		s.Emit("updateRecord", "")
		s.BroadcastTo("lobby", "updateRecord", "")
	})
}

func socketErorrHandler(s socketio.Socket, err error) {
	log.Println("error: ", err)
}
