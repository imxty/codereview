package image

import (
	"errors"
)

// GetImageSuffix
func GetImageSuffix(mime string) (string, error) {
	var (
		suffixMap = map[string]string{
			"image/jpg":  "jpg",
			"image/png":  "png",
			"image/jpeg": "jpeg",
			"image/webp": "webp",
		}
	)
	// 如果不再允许范围内则报错
	if suffix, ok := suffixMap[mime]; ok {
		return suffix, nil
	}
	return "", errors.New("invalid mime")
}
