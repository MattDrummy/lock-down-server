package main

import (
  "log"

  "github.com/googollee/go-socket.io"
  
)

func socketConnectionHandler(s socketio.Socket)  {
  s.On("open", func (user [2]string)  {
    s.Join(user[1])
    log.Println(user[0] + " connected to " + user[1])
    s.BroadcastTo(user[1], "message", user[0] + " connected to " + user[1])
  })
  s.On("message", func (msg [2]string)  {
    log.Printf("%+v\n", msg)
    s.Emit("message", msg[1])
    s.BroadcastTo(msg[0], "message", msg[1])
  })
  s.On("disconnection", func(){
    log.Println("a user has disconnected")
  })
}

func socketErorrHandler(s socketio.Socket, err error)  {
  log.Println("error: ", err)
}
