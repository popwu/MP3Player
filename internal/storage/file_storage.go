package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"mp3-player/internal/models"
)

type FileStorage struct {
	filename string
}

func NewFileStorage(filename string) *FileStorage {
	return &FileStorage{filename: filename}
}

func (fs *FileStorage) GetLibraries() ([]string, error) {
	data, err := ioutil.ReadFile(fs.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var libraries []string
	err = json.Unmarshal(data, &libraries)
	return libraries, err
}

func (fs *FileStorage) AddLibrary(path string) error {
	libraries, err := fs.GetLibraries()
	if err != nil {
		return err
	}

	libraries = append(libraries, path)
	data, err := json.Marshal(libraries)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fs.filename, data, 0644)
}

func (fs *FileStorage) GetSongs() ([]models.Song, error) {
	libraries, err := fs.GetLibraries()
	if err != nil {
		return nil, err
	}

	var songs []models.Song
	for _, library := range libraries {
		err := filepath.Walk(library, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && isSupportedAudioFile(path) {
				songs = append(songs, models.Song{
					Title: filepath.Base(path),
					Path:  path,
				})
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return songs, nil
}

func isSupportedAudioFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".mp3" || ext == ".flac"
}

func (fs *FileStorage) IsValidSongPath(path string) bool {
	libraries, err := fs.GetLibraries()
	if err != nil {
		return false
	}

	// 检查路径是否在已知的音乐库中
	for _, library := range libraries {
		if strings.HasPrefix(path, library) {
			// 检查文件扩展名
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".mp3" || ext == ".flac" {
				return true
			}
		}
	}

	return false
}
