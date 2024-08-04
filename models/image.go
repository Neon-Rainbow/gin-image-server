package models

import (
	"os"
)

// Image 图片模型
type Image struct {
	Filename string `json:"filename"`
	Filepath string `json:"filepath"`
}

// GetImage 检查图片是否存在
func GetImage(filepath string) (*Image, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, err
	}
	return &Image{
		Filepath: filepath,
	}, nil
}
