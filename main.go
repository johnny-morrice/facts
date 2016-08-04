package main

import (
        "strings"

        "github.com/pkg/errors"
)

func main() {
        frame := &RealityFrame{}

        frame.realm = &EloUnderground{}

        frame.structure = &Room{
                Walls: MakeWalls(MakeMaterials("brick:white", 1), 4)
                Floor: MakeMaterials("lino:blue", 1),
                Roof: MakeMaterials("brick:white", 1),
                Size: SMALL_ROOM,
        }

        frame.structure.Walls[0].Door = &Door{
                Panes: MakeMaterials("steel:unpainted")
        }

        frame.playerentity = MakeHuman("adult:male:prime")

        psych := PrisonCell()


}

type RealityFrame struct {

}

type Checker interface {
        Is(p Attribute) (bool, error)
}

type Describer interface {
        Describe(ps []Attribute) (string, error)
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

func (r *Room) Is(p Attribute) (bool, error) {
        switch (p)
        {
        case Tiny:
                return r.Size == CUPBOARD
        case Small:
                return r.Size == SMALL_ROOM
        case Inside:
                return true
        case Sealed:
                return true
        }
}

func (r *Room) Describe(attrs AttrSet) (string, error) {
        wallSet := map[Wall]int

        for _, w := r.Walls {
                if prev, ok := wallSet[w]; ok {
                        wallSet[w] = prev + 1
                } else {
                        wallSet[w] = 1
                }
        }

        out := ""

        if attrs.Has(Grim) {
                out = "This awful room is "
        } else {
                out = "This room is"
        }

        switch (r.Size) {
        case CUPBOARD:
                out += "tiny. "
        case SMALL_ROOM:
                out += "small. "
        case LARGE_ROOM:
                out += "large. "
        }

        out += fmt.Sprintf("It has %v walls", len(r.Walls))

        if len(wallSet) == 1 {
                out += ", all the same. "
        } else {
                out += ". "
        }

        for w, c := range wallSet {
                wd, err := w.Describe(attrs)

                if err != nil {
                        errors.Wrap(err, "error: describing room")
                }

                out += fmt.Sprintf("%v walls: %v", c, wd)
        }
}

type AttrSet map[Attribute]bool

func (a AttrSet) Has(attr Attribute) bool {
        return a[attr]
}

func Attributes(attrs... []Attribute) AttrSet {
        attrSet := map[Attribute]bool

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

type Wall struct {
        Panes []Material
        Door Door
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
        __clsmap[Inside] = PositionClass

        __clsmap[Tiny] = SizeClass
        __clsmap[Small] = SizeClass

        __clsmap[Confined] = ArchitectureClass

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
                all = all || pred
        }

        return all
}

type PerceptionFrame struct {
        Conditional
        Perception Perception
        Name string
}

func (pf *PerceptionFrame) Run(e Entity) (string, error) {
        if !pf.Valid() {
                return "", errors.New("Failed condition")
        }

        return pf.Perception.Describe(e)
}

type Perception interface {
        Describe(e Entity) (string, error)
}

func PrisonCell() PerceptionFrame {
        cell := PerceptionFrame{}

        cell.PlaceName = "Cell"

        cell.Format = &CellFormatter{}

        cell.AndConditions = []Attribute{
                Location,
                Sealed,
        }

        cell.OrConditions = []Attribute{
                Tiny,
                Small,
        }
}

type CellPerception struct {}

func (cf *CellPerception) Describe(e Entity) (string, error) {
        // TODO compute odd walls/floor
        // TODO use golang template library

        var out string

        if e.Is(Tiny) {
                out = "You are in a cramped prison cell.  The walls seem to press in around you."
        } else {
                out = "You are in a spacious prison cell."
        }

        out += e.Describe(Attributes(Grim))

        return out, nil
}
