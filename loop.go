package main

type Packet struct {
        Data interface{}
        Error error
}

type UpdateType string

const (
        TextUpdate = "TextUpdate"
        ChooseUpdate = "ChooseUpdate"
        ErrorUpdate = "ErrorUpdate"
        CloseUpdate = "CloseUpdate"
)

type Update struct {
        Packet
        Type UpdateType
}

type EventType string

const (
        NextEvent = "NextEvent"
        ChoiceEvent = "ChoiceEvent"
        BeginEvent = "BeginEvent"
        CloseEvent = "CloseEvent"
)

type Event struct {
        Packet
        Type EventType
}

type Game interface {
        Ctrl(ev Event) Update
}

func Loop(game Game, evch <-chan Event, outch chan<- Update) error {
        var err error

        for input := range evch {
                update := game.Ctrl(input)

                outch<- update

                if update.Type == ErrorUpdate {
                        err = update.Data.(error)
                        goto CLOSE
                }

                if update.Type == CloseUpdate {
                        goto CLOSE
                }
        }

CLOSE:
        close(outch)
        return err
}
