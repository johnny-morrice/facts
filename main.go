package main

import (
	"bufio"
	"log"
	"os"

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

func init() {
	__clsmap[Tiny] = SizeClass
	__clsmap[Small] = SizeClass

	__clsmap[Inside] = ArchitectureClass
	__clsmap[Sealed] = ArchitectureClass

	__clsmap[Location] = RootClass

	__clsmap[Grim] = EmotionClass
}
