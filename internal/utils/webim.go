package utils

import (
	"chaoxing/internal/models"
	"strings"
)

func ParseIMParams(data string) models.IMParamsType {
	var imParams models.IMParamsType

	index := strings.Index(data, "id=\"myName\"")
	if index == -1 {
		return imParams
	}
	start := index + 35
	end := strings.Index(data[start:], "<")
	if end == -1 {
		return imParams
	}
	imParams.MyName = data[start : start+end]

	// 提取 myToken
	index = strings.Index(data, "id=\"myToken\"")
	if index == -1 {
		return imParams
	}
	start = index + 36
	end = strings.Index(data[start:], "<")
	if end == -1 {
		return imParams
	}
	imParams.MyToken = data[start : start+end]

	// 提取 myTuid
	index = strings.Index(data, "id=\"myTuid\"")
	if index == -1 {
		return imParams
	}
	start = index + 35
	end = strings.Index(data[start:], "<")
	if end == -1 {
		return imParams
	}
	imParams.MyTuid = data[start : start+end]

	return imParams
}
