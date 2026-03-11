package client

import (
	"github.com/lumi4x/news/src/common/api"
	"github.com/sirupsen/logrus"
)

const (
	UnknownError api.ErrorCode = iota
)

var codeMessages = map[api.ErrorCode]string{
	UnknownError: "Unknown Error",
}

func CodeText(code int) string {
	if text, ok := codeMessages[code]; ok {
		return text
	}

	logrus.Error("unsupported code, defaulting to unknown")
	return codeMessages[UnknownError]
}
