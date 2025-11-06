package main

import "calendar-server/internal/app"

func main() {
	application := app.New()
	application.Run()
}
