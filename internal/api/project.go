package api

import (
	"database/sql"
	"errors"
	db "github.com/aalug/cv-backend-go/internal/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type listProjectsRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"` // profile cv id
}

type listProjectsQueryRequest struct {
	Page     int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=15"`
}

// @Schemes
// @Summary List projects for a profile cv
// @Description List projects for a profile cv with provided ID
// @Tags projects
// @Param id path integer true "CV profile ID"
// @Param page query integer true "Page number"
// @Param page_size query integer true "Page size"
// @Produce json
// @Success 200 {object} []db.ListProjectsWithTechnologiesRow
// @Failure 400 {object} ErrorResponse "Invalid ID, page or page size"
// @Failure 404 {object} ErrorResponse "CV profile with given ID does not exist"
// @Failure 500 {object} ErrorResponse "Any other server-side error"
// @Router /projects/{id} [get]
// listProjects returns a list of projects for a profile cv
func (server *Server) listProjects(ctx *gin.Context) {
	// get and validate the cv profile id
	var request listProjectsRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get and validate the query params - page and page size
	var queryRequest listProjectsQueryRequest
	if err := ctx.ShouldBindQuery(&queryRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get all projects for a profile cv
	params := db.ListProjectsWithTechnologiesParams{
		CvProfileID: request.ID,
		Limit:       queryRequest.PageSize,
		Offset:      (queryRequest.Page - 1) * queryRequest.PageSize,
	}

	projects, err := server.store.ListProjectsWithTechnologies(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, projects)
}
