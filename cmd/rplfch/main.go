package main

import (
	"fmt"
	"goapi/bilibili/comment"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <oid>", os.Args[0])
		os.Exit(0)
	}
	oid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	of, err := os.Create(os.Args[1] + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	replies, err := comment.GetReplies(comment.TYPE_TREND, oid, true)
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range replies {
		fmt.Fprintf(of, "%d,%d,%s\n", r.Uid, r.Rpid, r.Content)
	}
}
