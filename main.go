package main

import (
	"fmt"
	"game-backend-example/flatbuffer_model/mygame/sample"
	flatbuffers "github.com/google/flatbuffers/go"
)

func main() {
	// construct a Builder with 1024 byte backing array
	// builder will automatically resize the backing buffer when necessary
	builder := flatbuffers.NewBuilder(0)

	// string is a reference type, we first need to serialize it before
	// returns numerical offset to where that data resides in the buffer
	name := builder.CreateString("MonsterName")

	// backing array is the low-level array behind the scenes that holds the elements.

	sample.MonsterStart(builder)
	sample.MonsterAddName(builder, name)
	sample.MonsterAddPos(builder, sample.CreateVec3(
		builder, 1.0, 2.0, 3.0))
	sample.MonsterAddColor(builder, sample.ColorRed)

	// serialized monster to the flatbuffer and have its offset.
	monster := sample.MonsterEnd(builder)

	builder.Finish(monster)

	// buffer access - bytes can be sent over the network
	buf := builder.FinishedBytes()

	// deserialization
	// flatbuffers doesn't deserialize the whole buffer when accessed.
	// just decodes the data that is requested, leaving all the other data untouched.
	// deserialize -> accessing data from a binary flatbuffer

	de_monster := sample.GetRootAsMonster(buf, 0)

	fmt.Println(string(de_monster.Name()))
}
