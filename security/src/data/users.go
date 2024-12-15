package data

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Search for a user in the database using his UUID
func GetUserByUUID(uuid []byte) (*User, error) {
	user := User{}
	user_rows, err := DB.Query("SELECT * FROM users WHERE uuid = ?", uuid)
	if err != nil {
		return &User{}, errors.New("Unvalid Identifier")
	}
	for user_rows.Next() {
		err := user_rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email,&user.Avatar,&user.Liked_Posts,&user.Liked_Comments, &user.Admin)
		if err != nil {
			fmt.Print(err)
			return &User{}, errors.New("Unvalid Identifier")
		}
	}
	return &user, nil
}

// Get the connected user UUID
func GetCurrentUserID(r *http.Request) ([]byte, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, errors.New("User not logged in")
	}
	decodedSession, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, errors.New("Invalid session ID")
	}

	var userUUID []byte

	err = DB.QueryRow(`
	SELECT user_UUID FROM sessions WHERE session_id = ?`,
		decodedSession,
	).Scan(&userUUID)

	return userUUID, nil
}

func GetCurrentUser(r *http.Request) (*User, error) {
	id, err := GetCurrentUserID(r)
	if err != nil {
		return nil, err
	}
	return GetUserByUUID(id)
}

// Search for a user in the database using his username
func GetUserByUsername(s string) (*User, error) {
	user := User{}
	user_rows, err := DB.Query("SELECT * FROM users WHERE LOWER(username) = ?", strings.ToLower(s))
	if err != nil {
		return &User{}, errors.New("Unvalid Identifier")
	}
	for user_rows.Next() {
		err = user_rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Avatar,&user.Liked_Posts,&user.Liked_Comments, &user.Admin)
		if err != nil {
			fmt.Print(err)
			return &User{}, errors.New("Unvalid Identifier")
		}
	}
	return &user, nil
}

// Search for a user in the database using his email
func GetUserByEmail(s string) (*User, error) {
	user := User{}
	user_rows, err := DB.Query(`SELECT * FROM users WHERE "email" = ?;`, strings.ToLower(s))
	if err != nil {
		return &User{}, errors.New("Unvalid Identifier")
	}

	for user_rows.Next() {
		err = user_rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email,&user.Avatar,&user.Liked_Posts,&user.Liked_Comments, &user.Admin)
		if err != nil {
			fmt.Print(err)
			return &User{}, errors.New("Unvalid Identifier")
		}
	}
	return &user, nil
}

// Insert the user data in the data and check if username and email are not already used.
func InsertUser(user *User) error {
	if user.Avatar == "" {
		user.Avatar = "./uploads/avatars/default.png"
	}
	_, err := DB.Exec(`
	INSERT INTO users (
		"uuid",
		"username",
		"password",
		"email",
		"avatar",
		"admin"
	)
		VALUES (?,?,?,?,?,?);`, user.ID, user.Username, user.Password, strings.ToLower(user.Email), user.Avatar, user.Admin)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			return errors.New("Email already used, try to connect")
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
			return errors.New("Username already used")
		}
	}
	return nil
}

func IsSessionValid(sessionID string) (bool, error) {

	decodedSession, err := base64.StdEncoding.DecodeString(sessionID)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	var count int
	err = DB.QueryRow(`
        SELECT COUNT(*) FROM sessions
        WHERE session_id = ? AND expires_at > datetime('now')`,
		decodedSession,
	).Scan(&count)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return count > 0, nil
}

func InsertSession(session *Session) error {
	existingSession := &Session{}
	err := DB.QueryRow(`
		SELECT session_id FROM sessions WHERE user_UUID = ?`,
		session.UserUUID,
	).Scan(&existingSession.SessionID)

	if err != nil && err != sql.ErrNoRows {
		return err // Return any error except "no rows found"
	}

	if err == sql.ErrNoRows {
		// No existing session, insert a new one
		_, err = DB.Exec(`
			INSERT INTO sessions (session_id, user_UUID, expires_at)
			VALUES (?, ?, ?);`,
			session.SessionID, session.UserUUID, session.Expires_at.Format("2006-01-02 15:04:05"),
		)
	} else {
		// Existing session found, replace it
		_, err = DB.Exec(`
			UPDATE sessions
			SET session_id = ?, expires_at = ?
			WHERE user_UUID = ?;`,
			session.SessionID, session.Expires_at.Format("2006-01-02 15:04:05"), session.UserUUID,
		)
	}

	return err
}

func UpdateUserAvatar(userID []byte, avatarURL string) error {
	_, err := DB.Exec(`UPDATE users SET avatar = ? WHERE uuid = ?`, avatarURL, userID)
	return err
}

func ModifyUserLikesPost(user *User) {
	_, err := DB.Exec(`
	UPDATE users
	SET liked_posts = ?
	WHERE uuid = ?;
	`,user.Liked_Posts,user.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func SuppUserLikePost(user *User, post_id int) {
	tab := strings.Split(user.Liked_Posts,",")
	for i := 0; i < len(tab); i++ {
		if tab[i] == strconv.Itoa(post_id) || tab[i] == "" {
			tab = append(tab[:i],tab[i+1:]... )
		}
	}
	user.Liked_Posts = strings.Join(tab,",")
	ModifyUserLikesPost(user)
}

func ModifyUserLikesComm(user *User) {
	_, err := DB.Exec(`
	UPDATE users
	SET liked_comments = ?
	WHERE uuid = ?;
	`,user.Liked_Comments,user.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func SuppUserLikeComm(user *User, comment_id int) {
	tab := strings.Split(user.Liked_Comments,",")
	for i := 0; i < len(tab); i++ {
		if tab[i] == strconv.Itoa(comment_id) || tab[i] == "" {
			tab = append(tab[:i], tab[i+1:]...)
		}
	}
	user.Liked_Comments = strings.Join(tab,",")
	ModifyUserLikesComm(user)
}

func IsPostLiked(post *Post, user_id []byte) bool {
	if strings.Contains(string(post.Likes), string(user_id)) {
		return true
	}
	return false
}

func IsPostDisliked(post *Post, user_id []byte) bool {
	if strings.Contains(string(post.Dislikes), string(user_id)) {
		return true
	}
	return false
}

func IsCommLiked(comm *Comment , user_id []byte) bool {
	if strings.Contains(string(comm.Likes), string(user_id)) {
		return true
	}
	return false
}

func IsCommDisliked(comm *Comment, user_id []byte) bool {
	if strings.Contains(string(comm.Dislikes), string(user_id)) {
		return true
	}
	return false
}

func GetLikedPosts(user *User) []map[string]interface{} {
	posts := make([]map[string]interface{},0)
	tab := strings.Split(user.Liked_Posts,",")
	for _, post_id := range tab {
		if post_id == "" {
			continue
		}
		id, _ := strconv.Atoi(post_id)
		post, err := GetPostByID(id)
		if err != nil {
			fmt.Print(err)
			return []map[string]interface{}{}
		}
		temp := map[string]interface{}{}
		temp["Title"] = post.Title
		temp["Username"] = post.User.Username
		var topics []string
		for _, topic := range post.Topic {
			topics = append(topics, topic.Title)
		}
		temp["Topic"] = topics
		if post.Images != nil {
			temp["PostImage"] = post.Images.ImageUrl
		}
		temp["Avatar"] = post.User.Avatar
		temp["Date"] = post.Date.Format("02 Jan 06 15:04")
		temp["Content"] = post.Content
		temp["ID"] = post_id
		temp["Like"] = len(strings.Split(string(post.Likes), ",")) - 1
		temp["Dislike"] = len(strings.Split(string(post.Dislikes), ",")) - 1
		posts = append(posts, temp)
	}
	return posts
}
