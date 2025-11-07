package main

import (
	"calendar-server/internal/app"
	"calendar-server/pkg/logger/zappretty"
)

func main() {
	logger := zappretty.SetupLogger()

	application := app.New(logger)
	if err := application.Run(); err != nil {
		logger.Fatal("Failed to run application",
			zappretty.Field("error", err),
		)
	}
}
