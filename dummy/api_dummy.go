package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	phocus_api "github.com/wolffshots/phocus/v2/api"
	phocus_messages "github.com/wolffshots/phocus/v2/messages"
)

func randomLast() {
	for {
		input := fmt.Sprintf("(1 %d B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006\x36%x\r", rand.Intn(2000)+10000, rand.Intn(2)+10)
		actual, err := phocus_messages.InterpretQPGSn(input, rand.Intn(2)+1)
		if err == nil {
			phocus_api.SetLast(actual)
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {
	go randomLast()
	err := phocus_api.SetupRouter(gin.DebugMode).Run("0.0.0.0:8080")
	if err != nil {
		log.Fatalf("fatal err in router: %v", err)
	}
}
