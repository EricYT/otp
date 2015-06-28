package otp

import (
  "fmt"
  "time"
  "errors"
  "reflect"
)

var debugPrefixGenServer = "++++++++++(gen_server)++++++++++++> GenServer debug:"

// Genserver interface
type gsModIf interface {
  Init() error

  // handle call message
  HandleCall(args interface{}) (interface{}, error)

  // handle cast info
  HandleInfo(args interface{}) error

  // Debug interface
  HandleMessage(interface{}) (interface{}, error)
}


type GenServer struct {
  Name string            // gen_server name

  // message channel call
  callMsg chan interface{}

  // message channel cast
  castMsg chan interface{}

  // user gen_server interface
  gsMod gsModIf

  // call result
  result chan interface{}

  // gen model
  *gen
}

func NewGenServer(gIf gsModIf) *GenServer {
  return &GenServer{
    callMsg : make(chan interface{}),
    castMsg : make(chan interface{}),
    result : make(chan interface{}),
    gsMod : gIf,
    gen : NewGen(),
  }
}

func (gs *GenServer) Start() error {
  fmt.Println(debugPrefixGenServer, "gen_server start ")
  err := gs.StartGen(gs)
  if err != nil {
    fmt.Println(debugPrefixGenServer, "gen_server start error:", err)
    return err
  }
  return nil
}

func (gs *GenServer) InitIt() error {
  fmt.Println(debugPrefixGenServer, "gen_server init ")
  err := gs.gsMod.Init()
  if err != nil {
    fmt.Println(debugPrefixGenServer, "gen_server error:", err)
    return err
  }

  // receive message
  go gs.loop()

  return nil
}

func (gs *GenServer) Call(msg, timeout interface{}) (interface{}, error) {
  // send message to gen_server
  gs.callMsg<-msg

  // Block until receive the response
  switch timeout.(type) {
  case string:
    if timeout.(string) == INFINITY {
      fmt.Println(debugPrefixGenServer, "block forever until receive the response")
      return ReceiveInfinity(gs.result)
    }
    return nil, errors.New("timeout error")
  case int:
    if timeout.(int) >= 0 {
      fmt.Println(debugPrefixGenServer, "block until receive the response or timeout")
      return ReceiveTimeout(gs.result, timeout.(int))
    }
    return nil, errors.New("timeout must bigger than 0")
  case nil:
    fmt.Println(debugPrefixGenServer, "block until receive the response or default timeout")
    return ReceiveTimeout(gs.result, DEFAULT_TIMEOUT)
  }
  return nil, errors.New("timeout type error")
}

func (gs *GenServer) Cast(msg interface{}) {
  gs.castMsg<-msg
}

func (gs *GenServer) Debug(debugMsg interface{}) {
  gs.castMsg<-debugMsg
}

func (gs *GenServer) loop() {
  fmt.Println(debugPrefixGenServer, "gen_server loop: ", time.Now())
  select {
  case callMsg := <-gs.callMsg:
    fmt.Println(debugPrefixGenServer, "Recevie call message :", time.Now())
    res, err := gs.gsMod.HandleCall(callMsg)
    if err != nil {
      fmt.Println(debugPrefixGenServer, "Recevie call error:", err, time.Now())
      return
    }
    gs.result<-res
    gs.loop()
  case castMsg := <-gs.castMsg:
    fmt.Println(debugPrefixGenServer, "Recevie cast message :", time.Now())
    gs.gsMod.HandleInfo(castMsg)
    gs.loop()
  case <-time.After(time.Millisecond * 1):
    fmt.Println(debugPrefixGenServer, "Recevie timeout message ", time.Now())
    gs.gsMod.HandleInfo(TIMEOUT)
    gs.loop()
  }
}

