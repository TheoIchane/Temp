package handlers

import (
	"fmt"
	"forum/src/data"
	"net/http"
	"strconv"
	"time"
)

//
func CommentHandler(w http.ResponseWriter, r *http.Request) {
	user_id, err := data.GetCurrentUserID(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	user, err := data.GetUserByUUID(user_id)
	if err != nil {
		fmt.Println(err)
		return
	}
	id, err := strconv.Atoi(r.URL.RawQuery[3:])
	if err != nil {
		fmt.Println(err)
		return
	}
	post, err := data.GetPostByID(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	comment := data.Comment{
		User:    user,
		Post:    post,
		Content: r.FormValue("content"),
		Date:    time.Now(),
	}
	notif := CreateNotif(user,post,"comment")
	err = notif.InsertNotif()
	if err != nil {
		fmt.Println(err)
		ErrorHandler(w, r, http.StatusInternalServerError)
	}
	err = data.InsertComment(&comment)
	if err != nil {
		fmt.Println(err)
		ErrorHandler(w, r, http.StatusInternalServerError)
	}
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", post.ID), http.StatusSeeOther)
}
