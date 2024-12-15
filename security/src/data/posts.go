package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func GetPostByID(id int) (*Post, error) {
	post := Post{}
	user_rows, err := DB.Query("SELECT * FROM posts WHERE id = ?", id)
	if err != nil {
		return &Post{}, errors.New("Non Existing Post")
	}
	var user_id []byte
	var topics string
	var images_id int
	for user_rows.Next() {
		err = user_rows.Scan(&post.ID, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &user_id, &topics, &images_id, &post.Date)
		if err != nil {
			fmt.Print(err)
			return &Post{}, errors.New("Error Parsing Data 1")
		}
	}

	post.Images, err = GetImageByID(images_id)

	post.User, err = GetUserByUUID(user_id)
	if err != nil {
		return &Post{}, errors.New("Error Parsing Data 3")
	}
	for _, topic_str := range strings.Split(topics,",") {
		if topic_str == "" {
			continue
		}
		topic_id, _ := strconv.Atoi(topic_str)
		topic, _:= GetTopicByID(topic_id)
		post.Topic = append(post.Topic, topic)
	}
	if err != nil {
		// return &Post{}, errors.New("Error Parsing Data")
	}
	return &post, nil
}

func GetPostByTitle(title string) (*Post, error) {
	post := Post{}
	user_rows, err := DB.Query("SELECT * FROM posts WHERE title = ?", title)
	if err != nil {
		return &Post{}, errors.New("Non Existing Post")
	}
	var user_id []byte
	var topics string
	var images_id int
	for user_rows.Next() {
		err = user_rows.Scan(&post.ID, &post.Title, &post.Content, &post.Likes, &post.Dislikes, &user_id, &topics, &images_id, &post.Date)
		if err != nil {
			fmt.Print(err)
			return &Post{}, errors.New("Error Parsing Data 1")
		}
	}
	post.Images, err = GetImageByID(images_id)
	post.User, err = GetUserByUUID(user_id)
	if err != nil {
		return &Post{}, errors.New("Error Parsing Data 3")
	}
	for _, topic_str := range strings.Split(topics,",") {
		if topic_str == "" {
			continue
		}
		topic_id, _ := strconv.Atoi(topic_str)
		topic, _:= GetTopicByID(topic_id)
		post.Topic = append(post.Topic, topic)
	}
	return &post, nil
}

func InsertPost(post *Post) error {
	// var imageID interface{}
	// if post.Images != nil && post.Images.ID != 0 {
	// 	imageID = post.Images.ID
	// } else {
	// 	imageID = nil // No image uploaded, set to NULL
	// }
	var topics string
	for _, topic := range post.Topic {
		topics += strconv.Itoa(topic.ID) + ","
	}
	if post.Images != nil {
		_, err := DB.Exec(`
	INSERT INTO posts (
		"title",
		"content",
		"user_id",
		"topic_id",
		"image_id",
        "created_at"
	)
		VALUES (?,?,?,?,?,?);`, post.Title, post.Content, post.User.ID, topics, post.Images.ID, post.Date)

		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed: posts.title") {
				return errors.New("An existing post has the same title")
			}
			fmt.Println(err)
		}
		fmt.Println("Post Succesfuly Inserted")
	} else {
		_, err := DB.Exec(`
	INSERT INTO posts (
		"title",
		"content",
		"user_id",
		"topic_id",
        "created_at"
	)
		VALUES (?,?,?,?,?);`, post.Title, post.Content, post.User.ID, topics, post.Date)

		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed: posts.title") {
				return errors.New("An existing post has the same title")
			}
			fmt.Println(err)
		}
		fmt.Println("Post Succesfuly Inserted")
	}

	return nil
}

func GetLastPosts(user *User) []map[string]interface{} {
	posts := []map[string]interface{}{}
	rows, err := DB.Query("SELECT id FROM posts ORDER BY created_at DESC LIMIT 5")
	if err != nil {
		return posts
	}
	for rows.Next() {
		var post_id int
		temp := map[string]interface{}{}
		err = rows.Scan(&post_id)
		if err != nil {
			fmt.Print(err)
			return posts
		}
		post, err := GetPostByID(post_id)
		if err != nil {
			fmt.Println(err)
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
		temp["ID"] = post_id
		temp["Like"] = len(strings.Split(string(post.Likes), ",")) - 1
		temp["Dislike"] = len(strings.Split(string(post.Dislikes), ",")) - 1
		posts = append(posts, temp)
	}
	return posts
}

func GetUserPosts(user *User) []map[string]interface{} {
	posts := []map[string]interface{}{}
	rows, err := DB.Query("SELECT id FROM posts WHERE user_id = ? ", user.ID)
	if err != nil {
		return posts
	}
	for rows.Next() {
		var post_id int
		temp := map[string]interface{}{}
		err = rows.Scan(&post_id)
		if err != nil {
			fmt.Print(err)
			return posts
		}
		post, err := GetPostByID(post_id)
		if err != nil {
			fmt.Println(err)
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
		}
		temp["Date"] = post.Date.Format("02 Jan 06 15:04")
		temp["Content"] = post.Content
		temp["ID"] = post_id
		temp["Like"] = len(strings.Split(string(post.Likes), ",")) - 1
		temp["Dislike"] = len(strings.Split(string(post.Dislikes), ",")) - 1
		posts = append(posts, temp)
	}
	return posts
}

func FormatPosts(user *User, post *Post) []map[string]interface{} {
	var tab []map[string]interface{}
	temp := map[string]interface{}{}
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
		temp["IsUser"] = true
		temp["User_Like"] = IsPostLiked(post, user.ID)
		temp["User_Dislike"] = IsPostDisliked(post, user.ID)
	} else {
		temp["Popup"] = true
	}
	temp["Date"] = post.Date.Format("02 Jan 06 15:04")
	temp["Content"] = post.Content
	temp["ID"] = post.ID
	temp["Like"] = len(strings.Split(string(post.Likes), ",")) - 1
	temp["Dislike"] = len(strings.Split(string(post.Dislikes), ",")) - 1
	temp["Comments"] = GetPostComments(post,user)
	return append(tab, temp)
}

func ModifyPostLikes(post *Post) {
	_, err := DB.Exec(`
	UPDATE posts
	SET likes = ?, dislikes = ?
	WHERE id = ?;
	`, post.Likes, post.Dislikes, post.ID)
	if err != nil {
		return
	}
}

func SuppPostLike(post *Post, user_id []byte) {
	tab := strings.Split(string(post.Likes), ",")
	for i := 0; i < len(tab); i++ {
		if tab[i] == string(user_id) || tab[i] == "" {
			tab = append(tab[:i], tab[i+1:]...)
		}
	}
	post.Likes = []byte(strings.Join(tab, ","))
	ModifyPostLikes(post)
}

func SuppPostDislike(post *Post, user_id []byte) {
	tab := strings.Split(string(post.Likes), ",")
	for i := 0; i < len(tab); i++ {
		if tab[i] == string(user_id) || tab[i] == "" {
			tab = append(tab[:i], tab[i+1:]...)
		}
	}
	post.Dislikes = []byte(strings.Join(tab, ","))
	ModifyPostLikes(post)
}
func SuppPost(id int) {
	_, err := DB.Exec(`
	DELETE FROM posts
	WHERE id = ?;
	`, id)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println("Post deleted sucessfully")
}
