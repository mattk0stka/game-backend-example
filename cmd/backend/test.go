package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net"
	"path/filepath"
	"strconv"
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

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating listener")
	}

	log.Info().Msgf("listening on port 8000")

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal().Err(err).Msg("Error accepting connection")
			}

			defer conn.Close()
		}
	}()

	log.Debug().Msg("Listening on port 8000")
}
