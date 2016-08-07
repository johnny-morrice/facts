package main

import (
        "bufio"
        "fmt"
        "io"
        "log"
        "os"
        "math/rand"
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

        perc := PrisonCell()

        out := bufio.NewWriter(os.Stdout)

        err := perc.Describe(out, frame)

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

func (rf *RealityFrame) Is(attr Attribute) bool {
        return rf.Structure.Is(attr)
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
        trivial := map[Attribute]bool {
                Inside: true,
                Sealed: true,
                Location: true,
        }

        if trivial[attr] {
                return true
        }

        switch (attr) {
        case Tiny:
                return r.Size == CUPBOARD
        case Small:
                return r.Size == SMALL_ROOM
        }

        return false;
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

type Attribute string

const Location Attribute = "location"
const Inside Attribute = "inside"
const Sealed Attribute = "sealed"
const Tiny Attribute = "tiny"
const Small Attribute = "small"
const Grim Attribute = "grim"

type AttrClass Attribute

const RootClass AttrClass = "root"
const ArchitectureClass AttrClass = "architecture"
const SizeClass AttrClass = "size"
const EmotionClass AttrClass = "emotion"

var __clsmap = map[Attribute]AttrClass{}

func init() {
        __clsmap[Tiny] = SizeClass
        __clsmap[Small] = SizeClass

        __clsmap[Inside] = ArchitectureClass
        __clsmap[Sealed] = ArchitectureClass

        __clsmap[Location] = RootClass

        __clsmap[Grim] = EmotionClass
}

func GetClass(p Attribute) AttrClass {
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

func (c *Conditional) Validate(rf *RealityFrame) error {
        for _, pred := range c.AndConditions {
                if !rf.Is(pred) {
                        return fmt.Errorf("invalid AND condition: %v", pred)
                }
        }

        all := false
        for _, pred := range c.OrConditions {
                all = all || rf.Is(pred)
        }

        if !all {
                // Join all Attributes for error.
                ors := make([]string, len(c.OrConditions))
                for i, or := range c.OrConditions {
                        ors[i] = string(or)
                }

                together := strings.Join(ors, ", ")

                return fmt.Errorf("invalid OR condition: %v", together)

        }

        return nil
}

type Personality struct {
        Pframes []PerceptionFrame
}

func (p *Personality) ChooseFrame(rf *RealityFrame) (*PerceptionFrame, error) {
        matching := []PerceptionFrame{}

        for _, p := range p.Pframes {
                if err := p.Validate(rf); err == nil {
                        matching = append(matching, p)
                }
        }

        if len(matching) == 0 {
                return nil, errors.New("no matching perception frame")
        }

        index := rand.Intn(len(matching))

        return &matching[index], nil
}

type PerceptionFrame struct {
        Conditional
        Perception Perception
        Name string
}

func (pf *PerceptionFrame) Describe(w io.Writer, rf *RealityFrame) error {
        if err := pf.Validate(rf); err != nil {
                return errors.Wrap(err, "Failed condition")
        }

        return pf.Perception.Describe(w, rf)
}

type Perception interface {
        Describe(w io.Writer, rf *RealityFrame) error
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

func (cf *CellPerception) Describe(w io.Writer, rf *RealityFrame) error {
        if rf.Is(Tiny) {
                fmt.Fprint(w, "You are in a cramped prison cell. ")
                fmt.Fprint(w, "The walls seem to press in around you.")
        } else {
                fmt.Fprint(w, "You are in a prison cell.")
        }

        fmt.Fprint(w, "\n")

        err := rf.Describe(w, Attributes(Grim))

        if err != nil {
                return errors.Wrap(err, "CellPerception describe failed")
        }

        return nil
}
