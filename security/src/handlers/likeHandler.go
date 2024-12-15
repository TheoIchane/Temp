package handlers

import (
	"fmt"
	"forum/src/data"
	"net/http"
	"strconv"
	"strings"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := data.GetCurrentUserID(r)
	if err != nil {
		RenderIndex(w, r)
		return
	}
	cur_user, err := data.GetUserByUUID(user_id)
	if err != nil {
		RenderIndex(w, r)
		return
	}
	comment_id := ""
	if len(r.URL.Query()["comment_id"]) != 0 {
	comment_id = r.URL.Query()["comment_id"][0]
	}
	post_id := ""
	if len(r.URL.Query()["id"]) != 0 {
		post_id = r.URL.Query()["id"][0]
		}
	if comment_id != "" {
		LikeComment(comment_id,"like",cur_user)
	}else if post_id != "" {
		LikePost(post_id,"like",cur_user)
	} else {
		ErrorHandler(w,r,http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post?id=%s",post_id), http.StatusSeeOther)
}

func DislikeHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := data.GetCurrentUserID(r)
	if err != nil {
		RenderIndex(w, r)
		return
	}
	curr_user, err := data.GetUserByUUID(user_id)
	if err != nil {
		RenderIndex(w, r)
		return
	}
	comment_id := ""
	if len(r.URL.Query()["comment_id"]) != 0 {
	comment_id = r.URL.Query()["comment_id"][0]
	}
	post_id := ""
	if len(r.URL.Query()["id"]) != 0 {
		post_id = r.URL.Query()["id"][0]
		}
	if comment_id != "" {
		err = LikeComment(comment_id,"dislike",curr_user)
		if err != nil {
			ErrorHandler(w,r,http.StatusInternalServerError)
		}
	}else if post_id != "" {
		err = LikePost(post_id,"dislike",curr_user)
		if err != nil {
			ErrorHandler(w,r,http.StatusInternalServerError)
		}
	} else {
		ErrorHandler(w,r,http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post?id=%s",post_id), http.StatusSeeOther)
}

//Handling of Like/Dislike for the post with id, the s variable stand for the like or dislike.
func LikePost(post_id, s string, user *data.User) error {
	user_id := user.ID
	id, err := strconv.Atoi(post_id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	Post, err := data.GetPostByID(id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if s == "like" {
		if !strings.Contains(string(Post.Likes), string(user_id)) {
			if strings.Contains(string(Post.Dislikes), string(user_id)) {
				data.SuppPostDislike(Post, user_id)
			}
			user_id = append(user_id, ',')
			user.Liked_Posts += strconv.Itoa(Post.ID) + ","
			Post.Likes = append(Post.Likes, user_id...)
			data.ModifyUserLikesPost(user)
			data.ModifyPostLikes(Post)
			notif := CreateNotif(user,Post,"like")
			err = notif.InsertNotif()
		if err != nil {
			fmt.Println(err)
			return err
		}
		} else {
			data.SuppPostLike(Post, user.ID)
			data.SuppUserLikePost(user, Post.ID)
		}
	} else {
		if !strings.Contains(string(Post.Dislikes), string(user_id)) {
			if strings.Contains(string(Post.Likes), string(user_id)) {
				data.SuppPostLike(Post, user_id)
				data.SuppUserLikePost(user,Post.ID)
			}
			user_id = append(user_id, ',')
			Post.Dislikes = append(Post.Dislikes, user_id...)
			data.ModifyPostLikes(Post)
			notif := CreateNotif(user,Post,"dislike")
			err = notif.InsertNotif()
		if err != nil {
			fmt.Println(err)
			return err
		}
		} else {
			data.SuppPostDislike(Post, user_id)
		}
	}
	return nil
}


func LikeComment(comment_id, s string, user *data.User) error{
	user_id := user.ID
	id, err := strconv.Atoi(comment_id)
		if err != nil {
			fmt.Println(err)
			return err
		}
		Comment, err := data.GetCommentByID(id)
		if err != nil {
			fmt.Println(err)
			return err 
		}
	if s == "like" {
		if !strings.Contains(string(Comment.Likes), string(user_id)) {
			if strings.Contains(string(Comment.Dislikes), string(user_id)) {
				data.SuppCommentDislike(Comment, user_id)
			}
			user_id = append(user_id, ',')
			user.Liked_Comments += strconv.Itoa(Comment.ID) + ","
			Comment.Likes = append(Comment.Likes, user_id...)
			data.ModifyUserLikesComm(user)
			data.ModifyCommentLikes(Comment)
		} else {
			data.SuppCommentLike(Comment, user.ID)
			data.SuppUserLikeComm(user,Comment.ID)
		}
	} else {
		if !strings.Contains(string(Comment.Dislikes), string(user_id)) {
			if strings.Contains(string(Comment.Likes), string(user_id)) {
				data.SuppCommentLike(Comment, user_id)
				data.SuppUserLikeComm(user,Comment.ID)
			}
			user_id = append(user_id, ',')
			Comment.Dislikes = append(Comment.Dislikes, user_id...)
			data.ModifyCommentLikes(Comment)
		} else {
			data.SuppCommentDislike(Comment, user_id)
		}
	}
	return nil
}