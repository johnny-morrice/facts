package main

import (
	"fmt"
	"io"
	"strings"
)

type Checker interface {
	Is(p Attribute) bool
}

type Describer interface {
	Describe(w io.Writer, attrs AttrSet) error
}

type Attribute string

type AttrSet map[Attribute]bool

func (a AttrSet) Has(attr Attribute) bool {
	return a[attr]
}

func Attributes(attrs ...Attribute) AttrSet {
	attrSet := map[Attribute]bool{}

	for _, a := range attrs {
		attrSet[a] = true
	}

	return attrSet
}

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

func GetClass(p Attribute) AttrClass {
	cls, ok := __clsmap[p]

	if !ok {
		panic("bug: predicate without class")
	}

	return cls
}

type Conditional struct {
	AndConditions []Attribute
	OrConditions  []Attribute
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

	if len(c.OrConditions) == 0 {
		return nil
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

type Action uint

const (
	Look = Action(iota)
	Walk
	Touch
	Lick
	Listen
	Smell
	PickUp
	PutDown
)
