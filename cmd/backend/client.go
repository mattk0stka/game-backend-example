package main

import (
	"fmt"
	"game-backend-example/flatbuffer_model/mygame/network"
	flatbuffers "github.com/google/flatbuffers/go"
	"net"
)

func main() {

	builder := flatbuffers.NewBuilder(0)

	network.RequestStart(builder)
	network.RequestAddStatus(builder, 1)
	request := network.RequestEnd(builder)

	network.EnvelopeStart(builder)
	network.EnvelopeAddMsgType(builder, network.NetworkUnionRequest)
	network.EnvelopeAddMsg(builder, request)
	envelope := network.EnvelopeEnd(builder)
	builder.Finish(envelope)

	buf := builder.FinishedBytes()

	{ // client: send binary to backend server
		conn, err := net.Dial("tcp", ":8000")
		if err != nil {
			panic(err)
		}

		_, err = conn.Write(buf)
		if err != nil {
			fmt.Println("write to server failed", err.Error())
		}

		conn.Close()
	}

	/* - deserialize
	env := network.GetRootAsEnvelope(buf, 0)
	table := new(flatbuffers.Table)

	if env.Msg(table) {
		switch env.MsgType() {
		case network.NetworkUnionRequest:
			fmt.Println("something found")
				req := new(network.Request)
				req.Init(table.Bytes, table.Pos)
				fmt.Println(req.Status().String())
			req := network.GetRootAsRequest(table.Bytes, 0)
			fmt.Println(req.Status().String())
		}
	}
	*/

}
