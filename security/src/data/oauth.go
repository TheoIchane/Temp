package data

import (
	"errors"
	"fmt"
)

func GetCredentialsByID(id []byte) (*Credentials, error) {
	return nil, errors.New("test")
}

func (cred *Credentials) Insert() error {
	_, err := DB.Exec(`
		INSERT INTO credentials(
		ID,
		user_id,
		credential
		)
		VALUES (?,?,?);`,cred.ID,cred.User.ID,cred.Credential)
	if err != nil {
		return err
	}
	return nil
}

func IsCredExist(id []byte) bool {
	rows, err := DB.Query(`
	SELECT count(*) FROM credentials
	WHERE id = ?
	;`,id)
	if err != nil {
		fmt.Println(err)
		return false
	}
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
	}
	return count == 1
}