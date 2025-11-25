//go:build ignore
// +build ignore

package main

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ../../internal/infrastructure/database/ent/schema
