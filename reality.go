package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
)

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

type Entity interface {
	Checker
	Describer
}

type Structure Entity

type Room struct {
	Walls []Wall
	Floor []Material
	Size  RoomSize
}

func (r *Room) Is(attr Attribute) bool {
	trivial := map[Attribute]bool{
		Inside:   true,
		Sealed:   true,
		Location: true,
	}

	if trivial[attr] {
		return true
	}

	switch attr {
	case Tiny:
		return r.Size == CUPBOARD
	case Small:
		return r.Size <= SMALL_ROOM
	case Austere:
		return r.areWallsAustere()
	case CantLeave:
		return r.doorsWillEntrap()
	}

	return false
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

	switch r.Size {
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

type RoomSize uint8

const (
	CUPBOARD   = RoomSize(iota)
	SMALL_ROOM = RoomSize(iota)
	LARGE_ROOM = RoomSize(iota)
)

type Material struct {
	Name  string
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
	fmt.Fprintf(w, "%v %v", m.Color, m.Name)

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
	Door  *Door
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
	Panes  []Material
	Handle *DoorHandle
}

func (d *Door) Describe(w io.Writer, attrs AttrSet) error {
	fmt.Fprint(w, "This door is made from ")

	Matlist(d.Panes).Describe(w, attrs)

	fmt.Fprint(w, ".")

	if d.Handle != nil {
		if d.Handle.Material.Name != "" {
			newline(w)
			fmt.Fprint(w, "It has an unusual doorhandle made from ")

			err := d.Handle.Material.Describe(w, attrs)

			if err != nil {
				return errors.Wrap(err, "door describing handle")
			}

			fmt.Fprint(w, ".")
		}
	} else {
		newline(w)
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
		Name:  parts[0],
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
	return Wall{
		Panes: materials,
	}
}
