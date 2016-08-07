package main

import (
	"io"
	"math/rand"

	"github.com/pkg/errors"
)

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
			scored[score] = []*PerceptionFrame{perc}
		}
	}

	top := scored[highest]
	rogue := rand.Intn(len(top))

	return top[rogue], nil
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
