package utils

import "github.com/kaiquegarcia/gostudy/v2/logging"

func PanicHandler(logger logging.Logger) {
	err := recover()
	if err != nil {
		logger.Panic("%s\n", err)
	}
}
