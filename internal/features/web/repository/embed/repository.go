package web_embed_repository

import (
	"embed"
)

//go:embed *.html
var files embed.FS

type WebRepository struct{}

func NewWebRepository() *WebRepository {
	return &WebRepository{}
}
