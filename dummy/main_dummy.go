package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	phocus_api "github.com/wolffshots/phocus/v2/api"
)

func main() {
	if len(os.Args) <= 1 || os.Args[1] == "router" {
		go randomLast()
		if err := phocus_api.SetupRouter(gin.DebugMode, false).Run("0.0.0.0:8080"); err != nil {
			log.Fatalf("fatal error in router: %v", err)
		}
	} else if os.Args[1] == "serial" {
		log.Println("Starting dummy serial")
		serialDummy() // doesn't need goroutine, run on current thread
	} else {
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}
