package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
)

type PerceptionFrame struct {
	Conditional
	Perception Perception
	Name       string
}

func (pf *PerceptionFrame) Describe(w io.Writer, rf *RealityFrame) error {
	if err := pf.Validate(rf); err != nil {
		panic(errors.Wrap(err, "PerceptionFrame description"))
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

type WardPerception struct{}

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

type CellPerception struct{}

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
