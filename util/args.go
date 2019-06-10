package util

import (
	"log"
	"os"
	"strings"
)

func BodyFrom(args []string) string {
	log.Println(args)
	if (len(args) < 2) || os.Args[1] == "" {
		return "hello"
	}
	return strings.Join(args[1:], " ")
}

func SeverityFrom(args []string) string {
	if (len(args) < 2) || os.Args[1] == "" {
		return "info"
	}
	return os.Args[1]
}
