package main

import (
	"avito-backend/pkg/handlers"
	"fmt"
	"log"
	"os"
)

func main() {
	cmdArgs := os.Args[1:]
	if len(cmdArgs) != 1 {
		log.Fatalf("incorrect arguments quantity: %d\n", len(cmdArgs))
	}

	role := cmdArgs[0]
	token := handlers.NewJWT(role)
	if token == nil {
		log.Fatalln("failed to generate token")
	}
	fmt.Println(*token)
}
