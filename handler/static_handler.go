package handler

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:static
var web embed.FS

func NewEmbedStaticFileServer() http.Handler {
	subFS, _ := fs.Sub(web, "static")
	return http.FileServer(http.FS(subFS))
}
