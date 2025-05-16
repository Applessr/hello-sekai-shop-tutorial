package main

import (
	"context"
	"log"
	"os"

	"github.com/Applessr/hello-sekai-shop-tutorial/config"
	database "github.com/Applessr/hello-sekai-shop-tutorial/pkg"
	"github.com/Applessr/hello-sekai-shop-tutorial/server"
)

func main() {
	ctx := context.Background()

	//*Initialize config
	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path is required")
		}
		return os.Args[1]
	}())

	//Connect to database
	db := database.DbConnect(ctx, &cfg)
	defer db.Disconnect(ctx)

	//Start server
	server.Start(ctx, &cfg, db)

}
