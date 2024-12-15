package data

import (
	"errors"
	"fmt"
	"strings"
)

func GetTopicByID(id int) (*Topic,error) {
	topic := Topic{}
	user_rows, err := DB.Query("SELECT * FROM topics WHERE id = ?", id)
	if err != nil {
		return &Topic{}, errors.New("Non Existing Topic")
	}
	for user_rows.Next() {
		err = user_rows.Scan(&topic.ID, &topic.Title, &topic.Content)
		if err != nil {
			fmt.Print(err)
			return &Topic{}, errors.New("Error Parsing Data")
		}
	}
	return &topic, nil
}

func InsertTopic(topic *Topic) error {
	_, err := DB.Exec(`
	INSERT INTO posts (
		"title",
		"content",
	)
		VALUES (?,?);`, topic.Title, topic.Content)

	if err != nil {
		if strings.Contains(err.Error(),"UNIQUE constraint failed: topics.title") {
			return errors.New("Existing Topic")
		}
	}
	return nil
}

func GetAllTopics() []map[string]interface{} {
	topics := make([]map[string]interface{},0)
	rows, _ := DB.Query("SELECT id FROM topics")
	for rows.Next() {
		var topic_id int
		temp := map[string]interface{}{}
		rows.Scan(&topic_id)
		topic, err := GetTopicByID(topic_id)
		if err != nil {
			return topics
		}
		temp["Title"] = topic.Title
		temp["ID"] = topic_id
		temp["Content"] = topic.Content
		topics = append(topics, temp)
	}
	return topics
}

func GetTopicsPost(id int, user *User) []map[string]interface{} {
	posts := make([]map[string]interface{},0)
	id_string := fmt.Sprintf(",%d,",id)
	rows, err := DB.Query("SELECT id FROM posts WHERE topic_id LIKE '%' || ? || '%' ORDER BY created_at DESC",id_string)
	if err != nil {
		fmt.Println("SQL Query Error")
		return posts
	}
	for rows.Next() {
		temp := make(map[string]interface{})
		var post_id int
		err := rows.Scan(&post_id)
		if err != nil {
			fmt.Println("Scan Error")
			return posts
		}
		post, err := GetPostByID(post_id)
		if err != nil {
			fmt.Println("get post error")
			return posts
		}
		temp["Title"] = post.Title
		temp["Username"] = post.User.Username
		temp["Avatar"] = post.User.Avatar
		var topics []string
		for _, topic := range post.Topic {
			topics = append(topics, topic.Title)
		}
		temp["Topic"] = topics
		if post.Images != nil {
			temp["PostImage"] = post.Images.ImageUrl
		}
		if user != nil {
			temp["User_Like"] = IsPostLiked(post, user.ID)
			temp["User_Dislike"] = IsPostDisliked(post, user.ID)
		}else {
			temp["Popup"] = true
		}
		temp["Date"] = post.Date.Format("02 Jan 06 15:04")
		temp["Content"] = post.Content
		temp["ID"] = post.ID
		temp["Like"] = len(strings.Split(string(post.Likes), ",")) - 1
		temp["Dislike"] = len(strings.Split(string(post.Dislikes), ",")) - 1
		posts = append(posts, temp)
	}
	return posts
}