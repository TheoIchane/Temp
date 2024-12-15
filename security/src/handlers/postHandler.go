package handlers

import (
	"fmt"
	"forum/src/data"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var user *data.User
	cookie, err := r.Cookie("session_id")

	if err == nil && cookie.Value != "" {
		// Check if session is valid
		isValid, err := data.IsSessionValid(cookie.Value)
		if err != nil {
			fmt.Println("Error validating session:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if isValid {
			// Fetch the logged-in user
			userID, err := data.GetCurrentUserID(r)
			if err == nil {
				user, err = data.GetUserByUUID(userID)
				if err != nil {
					fmt.Println("Error fetching user:", err)
				}
			}
		} else {
			// Invalid session, clear the cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				MaxAge:   -1,
			})
		}
	} else {
		fmt.Println("No valid session, treating user as guest")
	}

	// Handle post retrieval
	post := &data.Post{}
	p := &Page{Data: make(map[string]interface{}), User: user}

	if strings.Contains(r.URL.RawQuery, "id") {
		id, err := strconv.Atoi(r.URL.RawQuery[3:])
		if err != nil {
			fmt.Println("Error parsing post ID:", err)
			ErrorHandler(w, r, http.StatusBadRequest)
			return
		}
		post, err = data.GetPostByID(id)
		if err != nil {
			fmt.Println("Error retrieving post by ID:", err)
			ErrorHandler(w, r, http.StatusNotFound)
			return
		}
		p.Title = "Post - " + post.Title
		p.Data["Posts"] = data.FormatPosts(user, post)
	} else if strings.Contains(r.URL.RawQuery, "title") {
		tempPost, err := data.GetPostByTitle(r.URL.RawQuery[6:])
		if err != nil {
			fmt.Println("Error retrieving post by title:", err)
			ErrorHandler(w, r, http.StatusNotFound)
			return
		}
		post = tempPost
		p.Title = "Post - " + post.Title
		p.Data["Posts"] = data.FormatPosts(user, post)
	} else {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	// Render the post page
	RenderTemplate(w, "post.html", p)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		// No valid cookie, user is not logged in
		fmt.Println("User is not logged in")
		RenderIndex(w, r)
		return
	}

	isValid, err := data.IsSessionValid(cookie.Value)
	if err != nil {
		// Handle database or validation error
		fmt.Println("Error validating session:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !isValid {
		// Invalid session, clear the cookie and redirect to login
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1, // Expire the cookie immediately
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "unable to process form", http.StatusBadRequest)
			return
		}
	}
	title := r.FormValue("title")
	content := verifContent(r.FormValue("content"))
	var topics []*data.Topic
	for _, topic_str := range r.Form["topic"] {
		id, _ := strconv.Atoi(topic_str)
		topic, _ := data.GetTopicByID(id)
		topics = append(topics, topic)
	}
	// topicID := 1                           //Default Value (Other Topic)
	imageIDStr := r.FormValue("image_id")  // Get the image_id from the form
	imageID, _ := strconv.Atoi(imageIDStr) // Convert to int
	// if topic_str != "" {
	// 	topicID, _ = strconv.Atoi(topic_str)
	// }
	userID, err := data.GetCurrentUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	user, _ := data.GetUserByUUID(userID)
	// topic, _ := data.GetTopicByID(topicID)
	image, err := data.GetImageByID(imageID)
	if title != "" && content != "" {
		post := data.Post{
			Title:   title,
			Content: content,
			User:    user,
			Date:    time.Now(),
		}
		if len(topics) == 0 {
			other_topic, _ := data.GetTopicByID(1)
			post.Topic = append(post.Topic, other_topic)
		} else {
			post.Topic = topics
		}
		if err == nil {
			post.Images = image
		} else {
			post.Images = nil
		}
		err = data.InsertPost(&post) // Pass the imageID to InsertPost
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	datas := map[string]interface{}{
		"Topics": make([]map[string]interface{}, 0),
	}
	datas["Topics"] = data.GetAllTopics()
	p := &Page{
		Title: "Create Post",
		User:  user,
		Data:  datas,
	}
	RenderTemplate(w, "post.html", p)
}

func verifContent(s string) string {
	tab := strings.Split(s, " ")
	final_tab := []string{}
	for _, word := range tab {
		if len(word) < 28 {
			final_tab = append(final_tab, word)
			continue
		}
		for i := 0; i <= len(word); i += 27 {
			if i+28 >= len(word) {
				final_tab = append(final_tab, word[i:len(word)-1])
				continue
			}
			final_tab = append(final_tab, word[i:i+28])
		}
	}
	return strings.Join(final_tab, " ")
}
