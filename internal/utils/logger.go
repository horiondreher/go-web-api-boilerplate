package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func StartLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339, // Customize the time format or use an empty string to hide the time
		NoColor:    false,        // Set to true if you do not want colored output
	})

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}
