package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kindiregg/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("failed to read config")
	}

	st := &state{Config: &cfg}

	fmt.Printf("DB URL: %s\n", st.DBUrl)
	fmt.Printf("Current User: %s\n", st.CurrentUsername)

	cmds := &commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Error: not enough arguments provided")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}

	if err = cmds.run(st, cmd); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

}
