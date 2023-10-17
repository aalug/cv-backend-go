package api

import (
	"database/sql"
	"errors"
	db "github.com/aalug/cv-backend-go/internal/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type listSkillsRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"` // profile cv id
}

// @Schemes
// @Summary List skills for a profile cv
// @Description List skills for a profile cv with provided ID
// @Tags skills
// @Param id path integer true "CV profile ID"
// @Produce json
// @Success 200 {object} []db.Skill
// @Failure 400 {object} ErrorResponse "Invalid ID"
// @Failure 404 {object} ErrorResponse "CV profile with given ID does not exist"
// @Failure 500 {object} ErrorResponse "Any other server-side error"
// @Router /skills/{id} [get]
// listSkills returns all skills for a profile cv
func (server *Server) listSkills(ctx *gin.Context) {
	var request listSkillsRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get all skills for a profile cv
	params := db.ListSkillsParams{
		CvProfileID: request.ID,
		Limit:       50,
		Offset:      0,
	}

	skills, err := server.store.ListSkills(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, skills)
}
