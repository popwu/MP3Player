package models

// Song 结构体表示一个音乐文件
type Song struct {
	Title    string `json:"title"`    // 歌曲标题
	Artist   string `json:"artist"`   // 艺术家
	Album    string `json:"album"`    // 专辑名
	Path     string `json:"path"`     // 文件路径
	Duration int    `json:"duration"` // 歌曲时长(秒)
}

// Playlist 结构体表示一个播放列表
type Playlist struct {
	Name  string `json:"name"`  // 播放列表名称
	Songs []Song `json:"songs"` // 播放列表中的歌曲
}

// Library 结构体表示一个音乐库
type Library struct {
	Path  string `json:"path"`  // 音乐库路径
	Songs []Song `json:"songs"` // 音乐库中的所有歌曲
}
