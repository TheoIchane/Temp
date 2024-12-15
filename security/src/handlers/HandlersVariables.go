package handlers

import "forum/src/data"

type Page struct {
	Title         string
	Error         string
	User          *data.User
	Notifications []string
	Data          map[string]interface{}
}

type Error struct {
	Error   error
	Code    int
	Message string
}
