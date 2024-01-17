package main

import (
	"fmt"
	// api "yourtube/api"
	// app "yourtube/internal/controllers"
	"yourtube/internal/server"
)

func main() {
	fmt.Println("checking main")
	server := server.NewServer()
	// app.Transcoder()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
