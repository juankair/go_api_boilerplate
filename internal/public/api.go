package public

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/uptrace/bunrouter"
	"net/http"
)

func RegisterHandler(router *bunrouter.Router) {
	attachmentBox := packr.New("attachments", "./public")
	fileServer := http.FileServer(attachmentBox)

	fileHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	})

	router.GET("/public/*path", bunrouter.HTTPHandlerFunc(fileHandler))

}
