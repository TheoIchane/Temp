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

const MaxUploadSize = 20 * 1024 * 1024

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)

	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		http.Error(w, "File is too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err == http.ErrMissingFile {
		http.Redirect(w, r, "/post", http.StatusSeeOther)
		return
	}
	if err != nil {
		http.Error(w, "Could not get the uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadDir := "./uploads/posts"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	imagePath := filepath.Join(uploadDir, filename)
	outFile, err := os.Create(imagePath)
	if err != nil {
		http.Error(w, "Unable to save the file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Failed to save the file", http.StatusInternalServerError)
		return
	}

	imageUrl := "/uploads/posts/" + filename
	images := &data.Images{
		ImageUrl: imageUrl,
	}

	// Insert image URL into the database
	imageID, err := data.InsertImage(images)
	if err != nil {
		http.Error(w, "Failed to insert image into database", http.StatusInternalServerError)
		return
	}

	// Redirect to post creation with the imageID (you may want to modify the URL handling as needed)
	http.Redirect(w, r, fmt.Sprintf("/post/create?image_id=%d", imageID), http.StatusSeeOther)
}
