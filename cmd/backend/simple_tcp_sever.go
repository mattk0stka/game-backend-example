package main

import (
	"flag"
	"fmt"
	"game-backend-example/flatbuffer_model/mygame/network"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {

	debug := flag.Bool("debug", false, "enable debug mode")

	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if *debug {
		fmt.Println("debug mode enabled")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	done := make(chan struct{})

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

		<-c
		close(done)
	}()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal().Err(err).Msg("creating listener failed")
	}

	log.Info().Msg("listening on port 8080")

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal().Err(err).Msg("error accepting connection")
			}

			log.Debug().Msg("accepted new connection")

			go handelTCPConnection(conn)
		}
	}()

	<-done
	log.Info().Msg("tcp server shutting down")
}

func handelTCPConnection(connection net.Conn) {
	defer connection.Close()

	n, err := connection.Write(createFlatBufferPkt())
	//n, err := connection.Write(createSimpleFlatBufferPkt())
	if err != nil {
		log.Debug().Err(err).Msg("error sending message")
	}

	log.Info().Msgf("sent %d bytes", n)
}

func createFlatBufferPkt() []byte {
	builder := flatbuffers.NewBuilder(0)

	// 1. - Request
	network.RequestStart(builder)
	network.RequestAddStatus(builder, 1)
	request := network.RequestEnd(builder)

	// 2. - wrap tp envelope
	network.EnvelopeStart(builder)
	network.EnvelopeAddMsgType(builder, network.NetworkUnionRequest)
	network.EnvelopeAddMsg(builder, request)
	envelope := network.EnvelopeEnd(builder)
	builder.Finish(envelope)

	return builder.FinishedBytes()
}

func createSimpleFlatBufferPkt() []byte {
	builder := flatbuffers.NewBuilder(0)

	network.RequestStart(builder)
	network.RequestAddStatus(builder, 1)
	request := network.RequestEnd(builder)

	builder.Finish(request)
	return builder.FinishedBytes()
}
