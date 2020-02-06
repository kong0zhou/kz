package kz

import (
	"net/http"
	"strings"
)

// IsTextFile 如果文件内容格式为纯文本或为空，则返回true。
func IsTextFile(data []byte) bool {
	if len(data) == 0 {
		return true
	}
	return strings.Contains(http.DetectContentType(data), "text/")
}

// IsImageFile 检测数据是否为图像格式
func IsImageFile(data []byte) bool {
	return strings.Contains(http.DetectContentType(data), "image/")
}

// IsPDFFile 检测数据是否为pdf格式
func IsPDFFile(data []byte) bool {
	return strings.Contains(http.DetectContentType(data), "application/pdf")
}

// IsVideoFile 检测数据是否为视频格式
func IsVideoFile(data []byte) bool {
	return strings.Contains(http.DetectContentType(data), "video/")
}

// IsAudioFile 检测数据是否为视频格式
func IsAudioFile(data []byte) bool {
	return strings.Contains(http.DetectContentType(data), "audio/")
}
