package route

import (
	"album.com/handlers"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

//Config is a func provides routing configuration
func Config() (router *httprouter.Router) {

	router = httprouter.New()
	handler := alice.New(loggingHandler, recoverHandler)

	/* Album api's */
	// Create a album with album name as query param
	router.POST("/album", wrapHandler(handler.ThenFunc(handlers.InsertAlbum)))
	// Delete an ablum with ablum name as query param
	router.DELETE("/album", wrapHandler(handler.ThenFunc(handlers.DeleteAlbum)))

	/* Image api's */
	// Create an img in an ablum with ablum name as query param
	router.POST("/image", wrapHandler(handler.ThenFunc(handlers.InsertImage)))
	// Delete an ablum with ablum, img name as query param
	router.DELETE("/image", wrapHandler(handler.ThenFunc(handlers.DeleteImage)))
	// Fetch an img in an ablum with ablum, img name as query param
	router.GET("/image", wrapHandler(handler.ThenFunc(handlers.GetImage)))
	// Fetch an img in an ablum with ablum name as query param
	router.GET("/images", wrapHandler(handler.ThenFunc(handlers.GetImages)))

	return
}
