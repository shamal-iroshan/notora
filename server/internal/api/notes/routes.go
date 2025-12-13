package notes

import "github.com/gin-gonic/gin"

func RegisterNoteRoutes(router *gin.RouterGroup, handler *NoteHandler, authMiddleware gin.HandlerFunc) {
	router.Use(authMiddleware)

	router.POST("/notes", handler.Create)
	router.GET("/notes", handler.GetAll)
	router.GET("/notes/meta", handler.Metadata)
	router.GET("/notes/:id", handler.Get)
	router.PUT("/notes/:id", handler.Update)
	router.PATCH("/notes/:id", handler.UpdateFlags)
	router.POST("/notes/:id/duplicate", handler.Duplicate)
	router.DELETE("/notes/:id", handler.DeleteForever)
	router.POST("/notes/search", handler.Search)

}
