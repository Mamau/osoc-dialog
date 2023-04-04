package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"osoc-dialog/internal/api/http/v1/request"
	"osoc-dialog/pkg/log"
	"osoc-dialog/pkg/router/middleware/auth/jwt"
)

type dialogRoutes struct {
	logger         log.Logger
	dialogProvider DialogProvider
}

func newDialogRoutes(group *gin.RouterGroup, l log.Logger, dp DialogProvider) {
	d := &dialogRoutes{
		logger:         l,
		dialogProvider: dp,
	}

	group.GET("/:user_id/list", d.list)
	group.POST("/:user_id/send", d.send)
}

func (d *dialogRoutes) list(c *gin.Context) {
	var req request.DialogList
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetInt(jwt.XUserIDKey)
	messageList, err := d.dialogProvider.Messages(c.Request.Context(), userID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": messageList})
}

func (d *dialogRoutes) send(c *gin.Context) {
	var req request.DialogMessage
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dialogMessage request.SendDialogMessage
	if err := c.ShouldBindJSON(&dialogMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt(jwt.XUserIDKey)

	if err := d.dialogProvider.SaveMessage(c.Request.Context(), req.UserID, userID, dialogMessage.Text); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, http.StatusText(http.StatusCreated))
}
