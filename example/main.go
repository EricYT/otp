package main

import (
  "fmt"
  "otp"
  "time"
)

var debugTest string = "********(test)*********> "

type test struct {
  name string
  *otp.GenServer
}

func (t *test) Init() error {
  fmt.Println(debugTest, "test init")
  return nil
}

func (t *test) HandleMessage(message interface{}) (interface{}, error) {
  fmt.Println(debugTest, "test handle message:", message.(string))
  t.name = message.(string)
  return message, nil
}

func (t *test) HandleCall(msg interface{}) (interface{}, error) {
  fmt.Println(debugTest, "test handle call")
  time.Sleep(time.Second*2)
  return msg, nil
}

func (t *test) HandleInfo(msg interface{}) error {
  fmt.Println(debugTest, "test handle cast")
  time.Sleep(time.Second*2)
  return nil
}

func main() {
  fmt.Println(debugTest, "start gen_server")

  t := &test{}
  genServer := otp.NewGenServer(t)
  err := genServer.Start()
  if err != nil {
    fmt.Println(debugTest, "test start gen_server error:", err)
    return
  }

  genServer.Cast("hello world")
  fmt.Println(debugTest, "test send cast info end:", time.Now())

  //res, err := genServer.Call("gen_server call", 4000)
  //res, err := genServer.Call("gen_server call", otp.INFINITY)
  res, err := genServer.Call("gen_server call", nil)
  if err != nil {
    fmt.Println(debugTest, "test call error:", err)
    return
  }
  fmt.Println(debugTest, "test call result:", res.(string), time.Now())

  fmt.Println(debugTest, "test name", t.name)

  stop := make(chan string, 1)
  go func() {
    time.Sleep(time.Second * 10)
    stop <-"stop"
  }()

  select {
  case <-stop:
    fmt.Println(debugTest, "test stop")
  }

  fmt.Println(debugTest, "test name", t.name)
}

