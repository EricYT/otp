package otp

import (
  "fmt"
)

var debugPrefix string = "--------(gen)---------> "

type GenModIf interface {
  // Init function
  InitIt() error

  //TODO:
}

type gen struct {
  // call wait channel
  result chan string
}

func NewGen() *gen {
  return &gen{
    result : make(chan string),
  }
}

func (g *gen) StartGen(genMod GenModIf) error {
  // start message routine to handle messages
  fmt.Println(debugPrefix, "start gen model")
  err := genMod.InitIt()
  if err != nil {
    fmt.Println(debugPrefix, "gen start error:", err)
    return err
  }
  return nil
}

