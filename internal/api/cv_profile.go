package api

import (
	"database/sql"
	"errors"
	db "github.com/aalug/cv-backend-go/internal/db/sqlc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getCvProfileRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type getCvProfileResponse struct {
	CvProfileID    int32            `json:"cv_profile_id"`
	Name           string           `json:"name"`
	Email          string           `json:"email"`
	Phone          string           `json:"phone"`
	Address        string           `json:"address"`
	LinkedinUrl    string           `json:"linkedin_url"`
	GithubUrl      string           `json:"github_url"`
	Bio            string           `json:"bio"`
	ProfilePicture string           `json:"profile_picture"`
	Education      []db.CvEducation `json:"education"`
}

// getCvProfile handles getting cv profile details
func (server *Server) getCvProfile(ctx *gin.Context) {
	var request getCvProfileRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get cv profile
	cvProfile, err := server.store.GetCvProfile(ctx, request.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// get cv education
	params := db.ListCvEducationsParams{
		CvProfileID: cvProfile.ID,
		Limit:       5,
		Offset:      0,
	}
	cvEducation, err := server.store.ListCvEducations(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// create a response
	cvProfileResponse := getCvProfileResponse{
		CvProfileID:    cvProfile.ID,
		Name:           cvProfile.Name,
		Email:          cvProfile.Email,
		Phone:          cvProfile.Phone,
		Address:        cvProfile.Address,
		GithubUrl:      cvProfile.GithubUrl,
		Bio:            cvProfile.Bio,
		ProfilePicture: cvProfile.ProfilePicture,
		Education:      cvEducation,
	}

	if cvProfile.LinkedinUrl.Valid {
		cvProfileResponse.LinkedinUrl = cvProfile.LinkedinUrl.String
	}

	ctx.JSON(http.StatusOK, cvProfileResponse)
}
