package main

import (
        "errors"
        "log"
        "testing"
        "time"
)

type gameStub struct {
        Counter int
        Events []EventType
        Updates []Update
}

func (stub *gameStub) Ctrl(ev Event) Update {
        current := stub.Counter
        stub.Counter++

        if stub.Events[current] != ev.Type {
                err := Update{Type: ErrorUpdate,}
                err.Data = errors.New("Unexpected Event")

                return err
        }

        return stub.Updates[current]
}

func TestLoop(t *testing.T) {
        text1 := Update{}
        text1.Type = TextUpdate
        text1.Data = "First text"

        text2 := Update{}
        text2.Type = TextUpdate
        text2.Data = "Second text"

        choice := Update{}
        choice.Type = ChooseUpdate
        choice.Data = []string{"Choice A", "Choice B"}

        stop := Update{}
        stop.Type = CloseUpdate

        stub := &gameStub{}
        stub.Events = []EventType{
                BeginEvent,
                NextEvent,
                ChoiceEvent,
                CloseEvent,
        }
        stub.Updates = []Update{
                text1,
                choice,
                text2,
                stop,
        }

        evch := make(chan Event)
        outch := make(chan Update)
        done := make(chan bool, 1)
        timeout := time.After(500 * time.Millisecond)

        go func() {
                select {
                case <-done:
                case <-timeout:
                        t.Error("timeout")
                }
        }()
        go eventGenerator(t, evch, outch)

        err := Loop(stub, evch, outch)
        done<- true

        if err != nil {
                t.Error("Game error:", err)
        }

        if stub.Counter != len(stub.Updates) {
                t.Error("Unexpected stub completion state:", stub)
        }
}

func eventGenerator(t *testing.T, evch chan<- Event, upch <-chan Update) {
        evch<- Event{Type: BeginEvent,}

        log.Println("client sent begin")

        txt := <-upch

        if txt.Type != TextUpdate {
                t.Error("Expected TextUpdate but got", txt)
        }

        evch<- Event{Type: NextEvent,}

        log.Println("client sent next")

        choose := <-upch

        if choose.Type != ChooseUpdate {
                t.Error("Expected ChooseUpdate but got", choose)
        }

        choice := Event{}
        choice.Type = ChoiceEvent
        choice.Data = 0

        evch<- choice

        log.Println("client sent choice")

        txt2 := <-upch

        if txt2.Type != TextUpdate {
                t.Error("Expected TextUpdate but got", txt2)
        }

        evch<- Event{Type: CloseEvent,}

        last := <-upch

        if last.Type != CloseUpdate {
                t.Error("Expected CloseUpdate but got", last)
        }

        close(evch)

        log.Println("client close")
}
