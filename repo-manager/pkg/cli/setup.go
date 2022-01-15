package cli

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup(component string) {
	if IsEnvSet("DEBUG") {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Logger = log.Logger.With().Str("component", component).Logger()
}
