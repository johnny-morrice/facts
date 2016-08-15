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
        Text() *Update
        Choose() *Update
        Ctrl() *Update
        React(ev *Event)
}

func Loop(game Game, outch chan<- *Update, evch <-chan *Event) {
        for reply := range evch {
                game.React(reply)

                ctrl := game.Ctrl()

                if ctrl != nil {
                        outch<- ctrl
                }

                txt := game.Text()

                outch<- txt

                choose := game.Choose()

                outch<- choose
        }
}
