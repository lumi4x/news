package client

import (
	"github.com/lumi4x/news/src/common/api"
	"github.com/sirupsen/logrus"
)

const (
	UnknownError api.ErrorCode = iota
	UserNotFound
	UserAlreadyBanned
	UserAlreadyUnbanned
	YouAreBanned
)

var codeMessages = map[api.ErrorCode]string{
	UnknownError:        "Unknown Error",
	UserNotFound:        "User not found",
	UserAlreadyBanned:   "User already banned",
	UserAlreadyUnbanned: "User is already unbanned",
	YouAreBanned:        "You've been banned from the game. Please contact us if you feel this is a mistake.",
}

func CodeText(code int) string {
	if text, ok := codeMessages[code]; ok {
		return text
	}

	logrus.Error("unsupported code, defaulting to unknown")
	return codeMessages[UnknownError]
}
