package data

import (
	"errors"
	"fmt"
	"strings"
)

func GetCommentByID(id int) (*Comment, error) {
	comment := Comment{}
	comment_rows, err := DB.Query("SELECT * FROM comments WHERE id = ?", id)
	if err != nil {
		return nil, errors.New("Non Existing Post")
	}
	var user_id []byte
	var post_id int
	for comment_rows.Next() {
		err = comment_rows.Scan(&comment.ID, &comment.Content, &post_id, &user_id,&comment.Likes,&comment.Dislikes, &comment.Date)
		if err != nil {
			return nil, err
		}
	}
	comment.User, err = GetUserByUUID(user_id)
	if err != nil {
		return nil, errors.New("Error Parsing Data2")
	}
	comment.Post, err = GetPostByID(post_id)
	if err != nil {
		return nil, errors.New("Error Parsing Data3")
	}
	return &comment, nil
}

func InsertComment(comment *Comment) error {
	_, err := DB.Exec(`
	INSERT INTO comments (
		"content",
		"user_id",
		"post_id",
        "created_at"
	)
		VALUES (?,?,?,?);`, comment.Content, comment.User.ID, comment.Post.ID, comment.Date)

	if err != nil {
		return err
	}
	return nil
}

func GetPostComments(post *Post, user *User) []map[string]interface{} {
	comments := []map[string]interface{}{}
	rows, _ := DB.Query(`
	SELECT id 
	FROM comments 
	WHERE post_id = ?
	ORDER BY created_at DESC;
	`, post.ID)
	for rows.Next() {
		temp := map[string]interface{}{}
		var comment_id int
		err := rows.Scan(&comment_id)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		comment, err := GetCommentByID(comment_id)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		temp["ID"] = comment_id
		temp["Post_ID"] = comment.Post.ID
		temp["Date"] = comment.Date.Format("02 Jan 06 15:04")
		if user != nil {
			temp["User_Like"] = IsCommLiked(comment, user.ID)
			temp["User_Dislike"] = IsCommDisliked(comment, user.ID)
		}else {
			temp["Popup"] = true
		}
		temp["Like"] = len(strings.Split(string(comment.Likes), ",")) - 1
		temp["Dislike"] = len(strings.Split(string(comment.Dislikes), ",")) - 1
		temp["Username"] = comment.User.Username
		temp["Avatar"] = comment.User.Avatar
		temp["Content"] = comment.Content
		comments = append(comments, temp)
	}
	return comments
}


func ModifyCommentLikes(comm *Comment) {
	_, err := DB.Exec(`
	UPDATE comments
	SET likes = ?, dislikes = ?
	WHERE id = ?;
	`, comm.Likes, comm.Dislikes, comm.ID)
	if err != nil {
		return
	}
}

func SuppCommentLike(comm *Comment, user_id []byte) {
	tab := strings.Split(string(comm.Likes), ",")
	for i := 0; i < len(tab); i++ {
		if tab[i] == string(user_id) || tab[i] == "" {
			tab = append(tab[:i], tab[i+1:]...)
		}
	}
	comm.Likes = []byte(strings.Join(tab, ","))
	ModifyCommentLikes(comm)
}

func SuppCommentDislike(comm *Comment, user_id []byte) {
	tab := strings.Split(string(comm.Likes), ",")
	for i := 0; i < len(tab); i++ {
		if tab[i] == string(user_id) || tab[i] == "" {
			tab = append(tab[:i], tab[i+1:]...)
		}
	}
	comm.Dislikes = []byte(strings.Join(tab, ","))
	ModifyCommentLikes(comm)
}