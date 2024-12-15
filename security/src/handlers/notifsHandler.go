package handlers

import (
	"fmt"
	"forum/src/data"
	"time"
)

func CreateNotif(user *data.User, post *data.Post, action string) *data.Notif {
	var content string
	switch action {
	case "like":
		content = fmt.Sprintf("%s a laissé un avis positif sur votre post: %s",user.Username,post.Title)
	case "dislike":
		content = fmt.Sprintf("%s a laissé un avis négatif sur votre post: %s",user.Username,post.Title)
	case "comment":
		content = fmt.Sprintf("%s a commenté votre post: %s",user.Username,post.Title)
	default :
		return nil
	}
	return &data.Notif{
		User: post.User,
		Post: post,
		Content: content,
		Date: time.Now(),
	}
}