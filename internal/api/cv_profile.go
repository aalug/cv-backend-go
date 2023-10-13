package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getCvProfileRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

// getCvProfile handles getting cv profile details
func (server *Server) getCvProfile(ctx *gin.Context) {
	var request getCvProfileRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cvProfile, err := server.store.GetCvProfile(ctx, request.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, cvProfile)
}
