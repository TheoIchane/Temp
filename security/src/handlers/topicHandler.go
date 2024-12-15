package handlers

import (
	"forum/src/data"
	"net/http"
	"strconv"
)

func TopicHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := data.GetCurrentUser(r)
	p := Page{
		Data: map[string]interface{}{},
	}
	if user != nil {
		p.User = user
	}
	if r.URL.Path != "/topics/" {
		id, err := strconv.Atoi(r.URL.Path[8:])
		if err != nil {
			ErrorHandler(w,r,http.StatusBadRequest)
			return
		}
		topic, err := data.GetTopicByID(id)
		if err != nil {
			ErrorHandler(w,r,http.StatusBadRequest)
			return
		}
		p.Title = topic.Title
		p.Data["Topic"] = topic
		p.Data["Posts"] = data.GetTopicsPost(id,user)
		RenderTemplate(w,"topic.html",p)
		return
	}
	p.Title = "Catégories"
	p.Data["Topics"] = data.GetAllTopics()
	RenderTemplate(w,"topic.html",p)
}