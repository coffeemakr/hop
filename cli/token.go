package cli

import (
	"io/ioutil"
	"os"
	"path"
)

type TokenStore interface {
	SaveToken(token string) error
	GetToken() (string, error)
}

type fileTokenStore struct {
	Path string
}

func NewFileTokenStore() (TokenStore, error) {
	folder, err := GetConfigFolderPath()
	if err != nil {
		return nil, err
	}
	path := path.Join(folder, "access-token.txt")
	return &fileTokenStore{Path: path}, nil
}

func (s *fileTokenStore) SaveToken(token string) error {
	file, err := os.Create(s.Path)
	if err != nil {
		return err
	}
	_, err = file.WriteString(token)
	if err != nil {
		return err
	}
	return nil
}

func (s *fileTokenStore) GetToken() (string, error) {
	file, err := os.Open(s.Path)
	if err != nil {
		if err == os.ErrNotExist {
			err = ErrNoTokenSaved
		}
		return "", err
	}
	token, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(token[:]), nil
}
