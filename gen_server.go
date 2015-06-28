package otp

import (
  "fmt"
  "time"
  "errors"
//  "reflect"
)

var debugPrefixGenServer = "++++++++++(gen_server)++++++++++++> GenServer debug:"

// Genserver interface
type gsModIf interface {
  Init() (string, interface{}, error)

  // handle call message
  HandleCall(args interface{}) (string, interface{}, error)

  // handle cast info
  HandleInfo(args interface{}) (string, int, error)

  // stop
  Stop(string) error

  // terminate
  Terminate(string)

  // Debug interface
  HandleMessage(interface{}) (string, interface{}, error)
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
  state, timeout, err := gs.gsMod.Init()
  if err != nil {
    fmt.Println(debugPrefixGenServer, "gen_server error:", err)
    return err
  }

  switch state {
  case STOP:
    return nil
  }

  // receive message
  go gs.loop(timeout)

  return nil
}

func (gs *GenServer) Stop(reason string) error {
  fmt.Println(debugPrefixGenServer, "gen_server stop")
  err := gs.gsMod.Stop(reason)
  if err != nil {
    fmt.Println(debugPrefixGenServer, "gen_server stop error:", err)
    return err
  }
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

func (gs *GenServer) loop(to interface{}) {
  fmt.Println(debugPrefixGenServer, "gen_server loop: ", time.Now())
  switch to.(type) {
  case int:
    fmt.Println(debugPrefixGenServer, "loop timeout:", to.(int))
    select {
    case callMsg := <-gs.callMsg:
      fmt.Println(debugPrefixGenServer, "Recevie call message :", time.Now())
      state, res, err := gs.gsMod.HandleCall(callMsg)
      if err != nil {
        fmt.Println(debugPrefixGenServer, "Recevie call error:", err, time.Now())
        return
      }

      // switch state
      fmt.Println(debugPrefixGenServer, "gen_server state:", state)
      switch state {
      case REPLY:
        fmt.Println(debugPrefixGenServer, "gen_server reply")
        // TODO
        gs.result<-res
        gs.loop(to)
      case NOREPLY:
        fmt.Println(debugPrefixGenServer, "gen_server noreply")
        // TODO
        gs.loop(to)
      case STOP:
        fmt.Println(debugPrefixGenServer, "gen_server stop")
        // TODO
        gs.terminate(NORMAL)
      }
    case castMsg := <-gs.castMsg:
      fmt.Println(debugPrefixGenServer, "Recevie cast message :", time.Now())

      state, timeout, err := gs.gsMod.HandleInfo(castMsg)
      if err != nil {
        fmt.Println(debugPrefixGenServer, "Recevie cast error:", err, time.Now())
        return
      }

      // switch state
      fmt.Println(debugPrefixGenServer, "gen_server state:", state, timeout)
      switch state {
      case STOP:
        fmt.Println(debugPrefixGenServer, "gen_server stop")
        // TODO
        gs.terminate(NORMAL)
      }
      gs.loop(timeout)
    case <-time.After(time.Millisecond * time.Duration(to.(int))):
      fmt.Println(debugPrefixGenServer, "Recevie timeout message ", time.Now())

      state, timeout, err := gs.gsMod.HandleInfo(TIMEOUT)
      if err != nil {
        fmt.Println(debugPrefixGenServer, "Recevie cast error:", err, time.Now())
        return
      }

      // switch state
      fmt.Println(debugPrefixGenServer, "gen_server timeout state:", state, timeout)
      switch state {
      case STOP:
        fmt.Println(debugPrefixGenServer, "gen_server timeout stop")
        // TODO
        gs.terminate(NORMAL)
      }
      gs.loop(timeout)
    }
  case string, nil:
    fmt.Println(debugPrefixGenServer, "loop not timeout")
    select {
    case callMsg := <-gs.callMsg:
      fmt.Println(debugPrefixGenServer, "Recevie call message :", time.Now())
      state, res, err := gs.gsMod.HandleCall(callMsg)
      if err != nil {
        fmt.Println(debugPrefixGenServer, "Recevie call error:", err, time.Now())
        return
      }

      // switch state
      fmt.Println(debugPrefixGenServer, "gen_server state:", state)
      switch state {
      case REPLY:
        fmt.Println(debugPrefixGenServer, "gen_server reply")
        // TODO
        gs.result<-res
        gs.loop(to)
      case NOREPLY:
        fmt.Println(debugPrefixGenServer, "gen_server noreply")
        // TODO
        gs.loop(to)
      case STOP:
        fmt.Println(debugPrefixGenServer, "gen_server stop")
        // TODO
        gs.terminate(NORMAL)
      }
    case castMsg := <-gs.castMsg:
      fmt.Println(debugPrefixGenServer, "Recevie cast message :", time.Now())

      state, timeout, err := gs.gsMod.HandleInfo(castMsg)
      if err != nil {
        fmt.Println(debugPrefixGenServer, "Recevie cast error:", err, time.Now())
        return
      }

      // switch state
      fmt.Println(debugPrefixGenServer, "gen_server state:", state, timeout)
      switch state {
      case STOP:
        fmt.Println(debugPrefixGenServer, "gen_server stop")
        // TODO
        gs.terminate(NORMAL)
      }
      gs.loop(timeout)
    }
  }
}

func (gs *GenServer) terminate(reason string) {
  fmt.Println(debugPrefixGenServer, "gen_server stoop:")
  gs.gsMod.Terminate(reason)
}
