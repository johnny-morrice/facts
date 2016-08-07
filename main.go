package main

import (
        "bufio"
        "fmt"
        "io"
        "log"
        "os"
        "strings"

        "github.com/pkg/errors"
)

func main() {
        frame := &RealityFrame{}

        //frame.realm = &EloUnderground{}

        room := &Room{
                Walls: MakeWalls(MakeMaterials("brick:white", 1), 4),
                //Floor: MakeMaterials("lino:blue", 1),
                //Roof: MakeMaterials("brick:white", 1),
                Size: SMALL_ROOM,
        }
/*
        room.Walls[0].Door = Door{
                Panes: MakeMaterials("steel:unpainted", 1),
        }
        */

        frame.Structure = room

        //frame.PlayerEntity = MakeHuman("adult:male:prime")

        /*
        psych := PrisonCell()

        */

        out := bufio.NewWriter(os.Stdout)

        err := frame.Describe(out, nil)

        if err != nil {
                log.Fatal(errors.Wrap(err, "error: "))
        }

        err = out.Flush()

        if err != nil {
                log.Fatal(errors.Wrap(err, "error: flushing buffer"))
        }
}

type RealityFrame struct {
        Structure Entity
        //PlayerEntity Entity
}

func (rf *RealityFrame) Describe(w io.Writer, attrs AttrSet) error {
        return rf.Structure.Describe(w, attrs)
}

type Checker interface {
        Is(p Attribute) bool
}

type Describer interface {
        Describe(w io.Writer, attrs AttrSet) error
}

type Entity interface {
        Checker
        Describer
}

type Structure Entity

type Room struct {
        Walls []Wall
        Floor []Material
        Size RoomSize
}

func (r *Room) Is(attr Attribute) bool {
        switch (attr) {
        case Tiny:
                return r.Size == CUPBOARD
        case Small:
                return r.Size == SMALL_ROOM
        case Inside:
                return true
        case Sealed:
                return true
        default:
                return false
        }
}

func (r *Room) Describe(w io.Writer, attrs AttrSet) error {
        if attrs.Has(Grim) {
                fmt.Fprint(w, "This awful room is ")
        } else {
                fmt.Fprint(w, "This room is ")
        }

        switch (r.Size) {
        case CUPBOARD:
                fmt.Fprint(w, "tiny. ")
        case SMALL_ROOM:
                fmt.Fprint(w, "small. ")
        case LARGE_ROOM:
                fmt.Fprint(w, "large. ")
        }

        fmt.Fprintf(w, "It has %v walls.", len(r.Walls))

        fmt.Fprintf(w, "\n")

        for _, wall := range r.Walls {
                err := wall.Describe(w, attrs)

                if err != nil {
                        return errors.Wrap(err, "room describing walls")
                }

                fmt.Fprintf(w, "\n")
        }

        return nil
}

type AttrSet map[Attribute]bool

func (a AttrSet) Has(attr Attribute) bool {
        return a[attr]
}

func Attributes(attrs... Attribute) AttrSet {
        attrSet := map[Attribute]bool{}

        for _, a := range attrs {
                attrSet[a] = true
        }

        return attrSet
}

type RoomSize uint8

const (
        CUPBOARD = RoomSize(iota)
        SMALL_ROOM = RoomSize(iota)
        LARGE_ROOM = RoomSize(iota)
)

type Material struct {
        Name string
        Color string
}

func (m *Material) Describe(w io.Writer, attrs AttrSet) error {
        fmt.Fprintf(w, "%v of %v", m.Name, m.Color)

        return nil
}

type Wall struct {
        Panes []Material
        Door Door
}

func (wall *Wall) Describe(w io.Writer, attrs AttrSet) error {
        if len(wall.Panes) == 0 {
                return errors.New("invisible wall")
        }

        fmt.Fprintf(w, "A wall made from ")

        between := ""
        for _, p := range wall.Panes {
                fmt.Fprintf(w, "%v", between)

                err := p.Describe(w, attrs)

                if err != nil {
                        return errors.Wrap(err, "wall pane error")
                }

                fmt.Fprintf(w, ".")

                between = " and "
        }

        return nil
}

type Door struct {
        Panes []Material
        Handle DoorHandle
}

type DoorHandle struct {
        Shape string
        Material
}

func MakeMaterial(desc string) Material {
        parts := strings.Split(desc, ":")

        return Material{
                Name: parts[0],
                Color: parts[1],
        }
}

func MakeMaterials(desc string, count int) []Material {
        out := make([]Material, count, count)
        copy := MakeMaterial(desc)

        for i, _ := range out {
                out[i] = copy
        }

        return out
}

func MakeWalls(materials []Material, count int) []Wall {
        out := make([]Wall, count, count)
        copy := MakeWall(materials)

        for i, _ := range out {
                out[i] = copy
        }

        return out
}

func MakeWall(materials []Material) Wall {
        return Wall {
                Panes: materials,
        }
}

type Action uint

const (
        Look = Action(iota)
        Walk
        Touch
        Lick
        Listen
        Smell
)

type Attribute uint

const (
        Location = Attribute(iota)

        Inside
        Sealed

        Tiny
        Small

        Grim
)

type PredicateClass uint

const (
        RootClass = PredicateClass(iota)

        ArchitectureClass

        SizeClass

        EmotionClass
)

var __clsmap = map[Attribute]PredicateClass{}

func init() {
        __clsmap[Tiny] = SizeClass
        __clsmap[Small] = SizeClass

        __clsmap[Inside] = ArchitectureClass
        __clsmap[Sealed] = ArchitectureClass

        __clsmap[Location] = RootClass

        __clsmap[Grim] = EmotionClass
}

func GetClass(p Attribute) PredicateClass {
        cls, ok := __clsmap[p]

        if !ok {
                panic("bug: predicate without class")
        }

        return cls
}

type Conditional struct {
        AndConditions []Attribute
        OrConditions []Attribute
}

func (c *Conditional) Valid(e Entity) bool {
        for _, pred := range c.AndConditions {
                if !e.Is(pred) {
                        return false
                }
        }

        all := false
        for _, pred := range c.OrConditions {
                all = all || e.Is(pred)
        }

        return all
}

type PerceptionFrame struct {
        Conditional
        Perception Perception
        Name string
}

func (pf *PerceptionFrame) Describe(w io.Writer, e Entity) error {
        if !pf.Valid(e) {
                return errors.New("Failed condition")
        }

        return pf.Perception.Describe(w, e)
}

type Perception interface {
        Describe(w io.Writer, e Entity) error
}

func PrisonCell() *PerceptionFrame {
        cell := &PerceptionFrame{}

        cell.Name = "Cell"

        cell.Perception = &CellPerception{}

        cell.AndConditions = []Attribute{
                Location,
                Sealed,
        }

        cell.OrConditions = []Attribute{
                Tiny,
                Small,
        }

        return cell
}

type CellPerception struct {}

func (cf *CellPerception) Describe(w io.Writer, e Entity) error {
        if e.Is(Tiny) {
                fmt.Fprint(w, "You are in a cramped prison cell. ")
                fmt.Fprint(w, "The walls seem to press in around you.")
        } else {
                fmt.Fprint(w, "You are in a spacious prison cell.")
        }

        err := e.Describe(w, Attributes(Grim))

        if err != nil {
                return errors.Wrap(err, "CellPerception describe failed")
        }

        return nil
}
