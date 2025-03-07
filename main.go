package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/kindiregg/gator/internal/config"
	"github.com/kindiregg/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("failed to read config")
	}

	st := &state{cfg: &cfg}

	// fmt.Printf("DB URL: %s\n", st.cfg.DBUrl)
	// fmt.Printf("Current User: %s\n", st.cfg.CurrentUsername)

	db, err := sql.Open("postgres", st.cfg.DBUrl)
	if err != nil {
		log.Fatalf("failed to connect to database at: %s", st.cfg.DBUrl)
	}

	dbQueries := database.New(db)

	st.db = dbQueries
	cmds := &commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

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
