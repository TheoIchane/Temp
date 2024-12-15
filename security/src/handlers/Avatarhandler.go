package handlers

import (
	"fmt"
	"forum/src/data"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const MaxAvatarUploadSize = 20 * 1024 * 1024 // Limit avatar size to 20mb

func UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := data.GetCurrentUserID(r)
	if err != nil {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	// Limit the upload size
	r.Body = http.MaxBytesReader(w, r.Body, MaxAvatarUploadSize)

	if err := r.ParseMultipartForm(MaxAvatarUploadSize); err != nil {
		fmt.Println(err)
		http.Error(w, "File is too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "Could not get the uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ensure the directory exists
	avatarDir := "./uploads/avatars"
	if _, err := os.Stat(avatarDir); os.IsNotExist(err) {
		os.Mkdir(avatarDir, os.ModePerm)
	}

	// Create a unique filename for the avatar
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	avatarPath := filepath.Join(avatarDir, filename)

	outFile, err := os.Create(avatarPath)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to save the file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Copy the uploaded file to the destination
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Failed to save the file", http.StatusInternalServerError)
		return
	}

	// Generate the avatar URL
	avatarUrl := "/uploads/avatars/" + filename

	// Update the user's avatar in the database
	err = data.UpdateUserAvatar(userID, avatarUrl)
	if err != nil {
		http.Error(w, "Failed to update avatar", http.StatusInternalServerError)
		return
	}

	// Redirect the user back to their profile
	http.Redirect(w, r, "/profil", http.StatusSeeOther)
}
