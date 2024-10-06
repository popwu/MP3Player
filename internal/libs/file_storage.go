package libs

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/qiniu/go-sdk/v7/storage"
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

func (fs *FileStorage) GetSongs() ([]Song, error) {
	libraries, err := fs.GetLibraries()
	if err != nil {
		return nil, err
	}

	var songs []Song
	for _, library := range libraries {
		err := filepath.Walk(library, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && isSupportedAudioFile(path) {
				songs = append(songs, Song{
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

func (fs *FileStorage) SaveFileList(fileList *[]storage.ListItem) error {
	// 将文件列表转换为 JSON
	jsonData, err := json.Marshal(fileList)
	if err != nil {
		return err
	}

	// 获取当前可执行文件的路径
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	// 获取应用程序的根目录（假设可执行文件在根目录下）
	rootDir := filepath.Dir(execPath)

	// 确定保存文件列表的路径
	filePath := filepath.Join(rootDir, "file_list.json")

	// 写入文件
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
