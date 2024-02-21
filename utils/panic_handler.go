package utils

func PanicHandler(logger Logger) {
	err := recover()
	if err != nil {
		logger.Panic("%s\n", err)
	}
}
