package otp

import (
  "fmt"
  "time"
  "errors"
)

var debugUtil string = "~~~~~~~~~~~~(util)~~~~~~~~~~~~~~~>"

// 
func ReceiveInfinity(in <-chan interface{}) (interface{}, error) {
  fmt.Println(debugUtil, "ReceiveIninity ")

  select {
  case res := <-in:
    fmt.Println(debugUtil, "Receive message")
    return res, nil
  }
}

func ReceiveTimeout(in <-chan interface{}, timeout int) (interface{}, error) {
  fmt.Println(debugUtil, "ReceiveTimeout ")

  select {
  case res := <-in:
    fmt.Println(debugUtil, "Receive message")
    return res, nil
  case <-time.After(time.Millisecond * time.Duration(timeout)):
    fmt.Println(debugUtil, "Receive message timeout")
    return nil, errors.New("timeout")
  }
}


