package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scch94/ins_log"
)

type Handler struct {
	utfi string
}

func NewHandler() *Handler {
	handler := &Handler{}
	handler.utfi = ins_log.GenerateUTFI()
	return handler
}
func (h *Handler) Welcome(c *gin.Context) {
	ctx := c.Request.Context()
	ctx = ins_log.SetPackageNameInContext(ctx, "handlerWelcome")
	ins_log.Info(ctx, "starting handler welcome")

	c.JSON(http.StatusOK, "bienvenidos")

}
