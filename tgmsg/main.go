package main

import (
	"fmt"
	"goapi/lib/telegram/bot"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <msg>...\n", os.Args[0])
		os.Exit(0)
	}
	content := strings.Join(os.Args[1:], " ")
	err := bot.SendPlain(content)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("OK")
}
