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
	params := db.ListProjectsParams{
		CvProfileID: request.ID,
		Limit:       queryRequest.PageSize,
		Offset:      (queryRequest.Page - 1) * queryRequest.PageSize,
	}

	projects, err := server.store.ListProjects(ctx, params)
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

type getProjectDetailsRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"` // project id
}

// getProjectDetails returns project details for a profile cv
func (server *Server) getProjectDetails(ctx *gin.Context) {
	var request getProjectDetailsRequest
	// validate the project id
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get project details
	project, err := server.store.GetProject(ctx, request.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, project)
}
