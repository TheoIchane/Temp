package data

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	// Open the database connection
	database, err := sql.Open("sqlite3", "src/data/forum.db")
	if err != nil {
		log.Fatal(err)
	}

	// Store the connection globally for further use
	DB = database

	// Create the users table
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		"uuid" BLOB NOT NULL PRIMARY KEY UNIQUE,
		"username" TEXT NOT NULL UNIQUE,
		"password" BLOB NOT NULL,
		"email" TEXT UNIQUE,
		"avatar" TEXT DEFAULT "/upload/default.png" NOT NULL,
		"liked_posts" TEXT DEFAULT "" NOT NULL,
		"liked_comments" TEXT DEFAULT "" NOT NULL,
		"admin" BINARY NOT NULL
	);`
	_, err = database.Exec(createUserTable)
	if err != nil {
		log.Fatal(err)
	}

	createSessionTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		"session_id" BLOB NOT NULL PRIMARY KEY UNIQUE,
		"user_UUID" BLOB NOT NULL,
		"expires_at" TIMESTAMP,
		FOREIGN KEY("user_uuid") REFERENCES users("uuid") ON DELETE CASCADE
	);`

	_, err = database.Exec(createSessionTable)
	if err != nil {
		log.Fatal(err)
	}

	createImageTable := `
	CREATE TABLE IF NOT EXISTS images (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	imageURL TEXT UNIQUE
	)`
	_, err = database.Exec(createImageTable)
	if err != nil {
		log.Fatal(err)
	}

	// Create the comments table
	createTopicTable := `
    CREATE TABLE IF NOT EXISTS topics (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
       	title TEXT UNIQUE,
		content TEXT
    );`
	_, err = database.Exec(createTopicTable)
	if err != nil {
		log.Fatal(err)
	}

	// Create the posts table
	createPostTable := `
    CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
		likes BLOB,
		dislikes BLOB,
        user_id BLOB,
		topic_id TEXT NOT NULL,
		image_id INTEGER DEFAULT 0, 
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(image_id) REFERENCES images(id),
        FOREIGN KEY(user_id) REFERENCES users(uuid)
    );`
	_, err = database.Exec(createPostTable)
	if err != nil {
		log.Fatal(err)
	}

	createCredentialsTable := `
	CREATE TABLE IF NOT EXISTS credentials (
		"id" BLOB NOT NULL UNIQUE,
		"credential" TEXT NOT NULL,
		"user_id" BLOB NOT NULL UNIQUE,
		FOREIGN KEY(user_id) REFERENCES user(uuid)
		);`

	_, err = database.Exec(createCredentialsTable)
	if err != nil {
		log.Fatal(err)
	}

	createNotifTable := `
	CREATE TABLE IF NOT EXISTS notifs (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"user_id" BLOB NOT NULL,
		"post_id" INT NOT NULL,
		"content" TEXT NOT NULL,
		"date" DATETIME,
		FOREIGN KEY("user_id") REFERENCES users("uuid"),
		FOREIGN KEY("post_id") REFERENCES posts("id")
	);`

	_, err = database.Exec(createNotifTable)
	if err != nil {
		log.Fatal(err)
	}
	// Create the comments table
	createCommentTable := `
    CREATE TABLE IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        content TEXT NOT NULL,
        post_id INTEGER,
        user_id BLOB,
		likes BLOB,
		dislikes BLOB,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY(post_id) REFERENCES posts(id),
        FOREIGN KEY(user_id) REFERENCES users(uuid)
    );`
	_, err = database.Exec(createCommentTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.Exec(`
	INSERT OR IGNORE INTO topics(title,content)
	VALUES ('Autres Catégories','All posts non related to an existing topic');
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.Exec(`
	INSERT OR IGNORE INTO topics(title,content)
	VALUES ('Jeux Vidéos','All posts related to video games');
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.Exec(`
	INSERT OR IGNORE INTO users(uuid,username,password,email,avatar,admin)
	VALUES (?,'Texio',?,'texio974@gmail.com','./uploads/avatars/default.png',1);
	`, []byte{62, 166, 203, 33, 145, 11, 65, 225, 156, 20, 110, 101, 63, 63, 139, 31}, []byte("$2a$15$eX/qrGiqOg6VVshuvAAfiOe.XhOazdre/x54/tvAmYf9t2qRTJtgG"))
	if err != nil {
		log.Fatal(err)
	}

	return database
}
