package handlers

import (
	"fmt"
	"forum/src/data"
	"net/http"
	"strconv"
)

func DeletePosts(w http.ResponseWriter, r *http.Request) {
	user_id, err := data.GetCurrentUserID(r)
	if err != nil {
		RenderIndex(w, r)
		return
	}
	curr_user,_ := data.GetUserByUUID(user_id)
	if !curr_user.Admin {
		ErrorHandler(w,r,http.StatusUnauthorized)
		return
	}
	id, err := strconv.Atoi(r.URL.RawQuery[3:])
	if err != nil {
		fmt.Println(err)
		return
	}
	data.SuppPost(id)
	http.Redirect(w,r,"/",http.StatusSeeOther)
}