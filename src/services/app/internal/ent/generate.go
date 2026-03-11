package ent

import (
	_ "github.com/lumi4x/news/src/services/app/internal/ent/build/runtime"
)

// NOTE: Specifying '--target build' requires it to point to a go package, so if this is the
// 			 first time generating, create a build/temp.go file with "package build" as its contents.

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/versioned-migration --feature sql/execquery --feature sql/upsert --feature intercept --feature sql/lock ./schema --target ./build
