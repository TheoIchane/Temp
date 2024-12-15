package main

// If we add a TCP server, use a go routine, we already are running the server in a go routine for graceful shutdown.

import (
	"fmt"
	"forum/src/data"
	s "forum/src/server"
	"os"
	"strings"
)

func main() {
	InitEnv()
	db := data.InitDB()
	defer db.Close()
	s.Server() // HTTP server runs inside a goroutine for graceful shutdown

	//go TcpServer() --> This is how we would add a TCP server if needed later
}

func InitEnv() {
	file, err := os.ReadFile(".env")
	if err != nil{
		fmt.Println(err)
		return
	}
	Envs := strings.Split(string(file),"\n")
	for _, v := range Envs {
		os.Setenv(strings.Split(v,"=")[0],strings.Split(v,"=")[1])
	}
}