package main

import (
	"flag"
	"fmt"
	"game-backend-example/flatbuffer_model/mygame/network"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {

	// setting global log level
	debug := flag.Bool("debug", false, "enable debug mode")

	flag.Parse()

	// default level is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	log.Logger = log.With().Caller().Logger()

	if *debug {
		fmt.Println("debug mode enabled")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// short termination solution
	// todo: something like this - https://victoriametrics.com/blog/go-graceful-shutdown/
	done := make(chan bool)

	go func() {
		log.Debug().Msg("listening signals...")
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		close(done)
	}()

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating listener")
	}

	log.Info().Msgf("listening on port 8000")
	log.Debug().Msg("listening for connections")

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal().Err(err).Msg("Error accepting connection")
			}

			log.Debug().Msgf("accepted new connection from %s", conn.RemoteAddr())

			go handleTCPConnection(conn)
			defer conn.Close()
		}
	}()

	<-done
	log.Info().Msg("backend server shutting down")
}

func handleTCPConnection(connection net.Conn) {
	defer connection.Close()

	packet := make([]byte, 0)
	tmp := make([]byte, 1024)

	for {
		n, err := connection.Read(tmp)
		if err != nil {
			if err != io.EOF {
				log.Err(err).Msg("Error reading from connection")
			}
			break
		}
		packet = append(packet, tmp[0:n]...)
	}

	log.Debug().Msgf("packet received: %d bytes", len(packet))

	env := network.GetRootAsEnvelope(packet, 0)

	table := new(flatbuffers.Table)
	if env.Msg(table) {
		switch env.MsgType() {
		case network.NetworkUnionRequest:
			log.Debug().Msg("union request found")
			request := network.GetRootAsRequest(table.Bytes, 0)
			fmt.Println(request.Status().String())
		}
	}

}
