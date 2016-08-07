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

        door := &Door{
                Panes: MakeMaterials("steel:unpainted", 1),
        }

        room.Walls[0].Door = door

        if len(os.Args) > 1 {
                door.Handle = &DoorHandle{}
        } 

        frame.Structure = room

        //frame.PlayerEntity = MakeHuman("adult:male:prime")

        psych := Personality{}

        psych.Pframes = []*PerceptionFrame{
                PrisonCell(),
                SimpleWard(),
        }

        out := bufio.NewWriter(os.Stdout)

        err := psych.Describe(out, frame)

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
                return r.Size <= SMALL_ROOM
        case Austere:
                return r.areWallsAustere()
        case CantLeave:
                return r.doorsWillEntrap()
        }

        return false;
}

func (r *Room) doorsWillEntrap() bool {
        for _, w := range r.Walls {
                if w.Door != nil && w.Door.Handle != nil {
                        return false
                }
        }

        return true
}

func (r *Room) areWallsAustere() bool {
        for _, w := range r.Walls {
                if !w.isAustere() {
                        return false
                }
        }

        return true
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

        newline(w)

        for _, wall := range r.Walls {
                err := wall.Describe(w, attrs)

                if err != nil {
                        return errors.Wrap(err, "room describing walls")
                }

                newline(w)
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

type Matlist []Material

func (m *Material) isAustere() bool {
        austere := map[string]bool{
                "brick": true,
                "steel": true,
                "tiles": true,
        }

        return austere[m.Name]
}

func (m *Material) Describe(w io.Writer, attrs AttrSet) error {
        fmt.Fprintf(w, "%v of %v", m.Name, m.Color)

        return nil
}

func (mats Matlist) Describe(w io.Writer, attrs AttrSet) error {
        between := ""
        for _, m := range mats {
                fmt.Fprint(w, between)

                err := m.Describe(w, attrs)

                if err != nil {
                        return err
                }

                between = " and "
        }

        return nil
}

type Wall struct {
        Panes []Material
        Door *Door
}

func (wall *Wall) isAustere() bool {
        for _, m := range wall.Panes {
                if !m.isAustere() {
                        return false
                }
        }

        return true
}

func (wall *Wall) Describe(w io.Writer, attrs AttrSet) error {
        if len(wall.Panes) == 0 {
                return errors.New("invisible wall")
        }

        fmt.Fprintf(w, "A wall made from ")

        Matlist(wall.Panes).Describe(w, attrs)

        fmt.Fprint(w, ".")

        if wall.Door != nil {
                newline(w)
                fmt.Fprintf(w, "There is a door here.")
                newline(w)
                wall.Door.Describe(w, attrs)
        }

        return nil
}

type Door struct {
        Panes []Material
        Handle *DoorHandle
}

func (d *Door) Describe(w io.Writer, attrs AttrSet) error {
        fmt.Fprint(w, "This door is made from ")

        Matlist(d.Panes).Describe(w, attrs)

        fmt.Fprint(w, ".")

        newline(w)

        if d.Handle != nil {
                if d.Handle.Material.Name != "" {
                        fmt.Fprint(w, "It has an unusual doorhandle made from ")

                        err := d.Handle.Material.Describe(w, attrs)

                        if err != nil {
                                return errors.Wrap(err, "door describing handle")
                        }

                        fmt.Fprint(w, ".")
                        newline(w)
                }
        } else {
                fmt.Fprint(w, "It has no door handle.")
        }

        return nil
}

type DoorHandle struct {
        Material Material
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
const Austere Attribute = "austere"
const CantLeave Attribute = "CantLeave"

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

func (c *Conditional) MatchingConditions(rf *RealityFrame) []Attribute {
        hold := []Attribute{}

        allattrs := append(c.AndConditions, c.OrConditions...)

        for _, attr := range allattrs {
                if rf.Is(attr) {
                        hold = append(hold, attr)
                }
        }

        return hold
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
        Pframes []*PerceptionFrame
}

func (p *Personality) ChooseFrame(rf *RealityFrame) (*PerceptionFrame, error) {
        matching := []*PerceptionFrame{}

        for _, p := range p.Pframes {
                if err := p.Validate(rf); err == nil {
                        matching = append(matching, p)
                }
        }

        if len(matching) == 0 {
                return nil, errors.New("no matching perception frame")
        }

        scored := map[int][]*PerceptionFrame{}
        highest := 0

        for _, perc := range matching {
                score := len(perc.MatchingConditions(rf))

                if score > highest {
                        highest = score
                }

                if peers, ok := scored[score]; ok {
                        scored[score] = append(peers, perc)
                } else {
                        scored[score] = []*PerceptionFrame{perc,}
                }
        }

        top := scored[highest]

        rouge := rand.Intn(len(top))

        return top[rouge], nil
}

func (psych *Personality) Describe(w io.Writer, rf *RealityFrame) error {
        perc, err := psych.ChooseFrame(rf)

        if err != nil {
                return errors.Wrap(err, "personality choosing pframe")
        }

        err = perc.Describe(w, rf)

        if err != nil {
                return errors.Wrap(err, "personalty describing perception")
        }

        return nil
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
                Austere,
                CantLeave,
        }

        cell.OrConditions = []Attribute{
                Tiny,
                Small,
        }

        return cell
}

func SimpleWard() *PerceptionFrame {
        ward := &PerceptionFrame{}

        ward.Name = "SimpleWard"

        ward.Perception = &WardPerception{}

        ward.AndConditions = []Attribute{
                Location,
                Sealed,
                Austere,
        }

        return ward
}

type WardPerception struct {}

func (WardPerception) Describe(w io.Writer, rf *RealityFrame) error {
        if rf.Is(Tiny) {
                fmt.Fprintf(w, "You are in a simple hospital room.")
        } else if rf.Is(Small) {
                fmt.Fprintf(w, "You are in a small hospital ward.")
        } else {
                fmt.Fprintf(w, "You are in a hospital ward.")
        }

        newline(w)

        err := rf.Describe(w, Attributes(nil...))

        if err != nil {
                return errors.Wrap(err, "WardPerception describe failed")
        }

        return nil
}

type CellPerception struct {}

func (CellPerception) Describe(w io.Writer, rf *RealityFrame) error {
        if rf.Is(Tiny) {
                fmt.Fprint(w, "You are in a cramped prison cell.")
                fmt.Fprint(w, "The walls seem to press in around you.")
        } else {
                fmt.Fprint(w, "You are in a prison cell.")
        }

        newline(w)

        err := rf.Describe(w, Attributes(Grim))

        if err != nil {
                return errors.Wrap(err, "CellPerception describe failed")
        }

        return nil
}

func newline(w io.Writer) {
        fmt.Fprint(w, "\n")
}
