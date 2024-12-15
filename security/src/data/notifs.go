package data

import (
	"fmt"
)

func GetNotifById(id int) *Notif {
	notif := &Notif{}
	rows, err := DB.Query(`SELECT * FROM notifs WHERE id = ?;`,)
	if err != nil {
		return nil
	}
	var user_id []byte
	var post_id int
	for rows.Next() {
		err = rows.Scan(&notif.ID,&user_id,&post_id,&notif.Content,&notif.Date)
		user, err := GetUserByUUID(user_id)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		post, err := GetPostByID(post_id)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		notif.User = user
		notif.Post = post
	}
	return notif
}

func (notif *Notif) InsertNotif() error {
	_, err := DB.Exec(`
	INSERT INTO notifs (
		"user_id",
		"post_id",
		"content",
		"date"
	)
		VALUES (?,?,?,?);`,notif.User.ID,notif.Post.ID,notif.Content,notif.Date)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}