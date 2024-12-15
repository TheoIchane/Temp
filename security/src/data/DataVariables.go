package data

import (
	"database/sql"
	"time"
)

type User struct {
	ID       []byte //UUID
	Email    string
	Username string
	Password []byte //encrypted password
	Avatar string
	Liked_Posts string
	Liked_Comments string
	Admin bool
}

type Session struct {
	SessionID  []byte
	UserUUID   []byte
	Expires_at time.Time
}

type Post struct {
	ID       int
	Title    string
	Content  string
	Images   *Images
	User     *User
	Likes    []byte
	Dislikes []byte
	Topic    []*Topic
	Date     time.Time
}

type Images struct {
	ID       int
	ImageUrl string
}

type Topic struct {
	ID      int
	Title   string
	Content string
}

type Comment struct {
	ID       int
	Content  string
	Likes    []byte
	Dislikes []byte
	User     *User
	Post     *Post
	Date     time.Time
}


type Notif struct {
	ID int
	User *User
	Post *Post
	Content string
	Date time.Time
}

type Credentials struct {
	ID []byte
	User *User
	Credential string
}

var DB *sql.DB // Global variable to store the database connection
