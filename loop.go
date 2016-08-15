package main

type Packet struct {
        Data interface{}
        Error error
}

type UpdateType string

const (
        SliceUpdate = "SliceUpdate"
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
        DoEvent = "DoEvent"
        MatchEvent = "MatchEvent"
        BeginEvent = "BeginEvent"
        QuitEvent = "QuitEvent"
)

type Event struct {
        Packet
        Type EventType
}

type Game interface {
        Ctrl(ev *Event) (*Update, error)
}

func Loop(game Game, outch chan<- *Update, evch <-chan *Event) {
        for input := range evch {
                update, err := game.Ctrl(input)

                if err != nil {
                        outch<- fail(err)
                }

                if update == nil {
                        outch<- end()
                        return
                } else {
                        outch<- update
                }
        }
}

func end() *Update {
        return &Update{
                Type: CloseUpdate,
        }
}

func fail(err error) *Update {
        bad := &Update{}
        bad.Type = ErrorUpdate
        bad.Data = err

        return bad
}
