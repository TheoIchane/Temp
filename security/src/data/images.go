package data

import (
	"errors"
	"fmt"
	"strings"
)

func InsertImage(images *Images) (int, error) {
	result, err := DB.Exec(`
    INSERT INTO images (imageUrl)
    VALUES (?);`, images.ImageUrl)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: images.ImageUrl") {
			return 0, errors.New("Existing Image")
		}
		return 0, err
	}

	// Get the last inserted ID
	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, errors.New("Failed to retrieve image ID")
	}

	return int(lastID), nil
}

func GetImageByID(id int) (*Images, error) {
	if id == 0 {
		return nil,errors.New("No Image")
	}
	images := Images{}
	user_rows, err := DB.Query("SELECT * FROM images WHERE id = ?", id)
	if err != nil {
		return &Images{}, errors.New("Non Existing Images")
	}
	for user_rows.Next() {
		err = user_rows.Scan(&images.ID, &images.ImageUrl)
		if err != nil {
			fmt.Print(err)
			return &Images{}, errors.New("Error Parsing Data")
		}
	}
	return &images, nil
}
