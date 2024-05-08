package l

import (
	"github.com/natefinch/lumberjack"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// var logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
var Z = zerolog.New(&lumberjack.Logger{
	Filename:   "logs/zero.log", // File name
	MaxSize:    20,              // Size in MB before file gets rotated
	MaxBackups: 50,              // Max number of files kept before being overwritten
	MaxAge:     7,               // Max number of days to keep the files
	Compress:   true,            // Whether to compress log files using gzip
	LocalTime:  true,            // Local time zone
}).With().Timestamp().Logger()

// Same as 'Error' but with stack trace which makes it ~7x slower
func ErrorTrace(err error) {

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	err = errors.WithStack(err)
	Z.Error().Stack().Err(err).Send()
}

func Error(err error) {
	Z.Error().Err(err).Send()
}

func Warning(warn string) {
	Z.Warn().Msg(warn)
}

func Info(inf string) {
	Z.Info().Msg(inf)
}
