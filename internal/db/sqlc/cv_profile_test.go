package db

import (
	"context"
	"database/sql"
	"github.com/aalug/cv-backend-go/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

// createRandomCvProfile creates and return a random cv profile
func createRandomCvProfile(t *testing.T) CvProfile {
	params := CreateCvProfileParams{
		Name:    utils.RandomString(5),
		Email:   utils.RandomEmail(),
		Phone:   utils.RandomString(9),
		Address: utils.RandomString(5),
		LinkedinUrl: sql.NullString{
			String: utils.RandomString(5),
			Valid:  utils.RandomInt(0, 1) == 1,
		},
		GithubUrl:      utils.RandomString(5),
		Bio:            utils.RandomString(10),
		ProfilePicture: utils.RandomString(6),
	}

	cvProfile, err := testQueries.CreateCvProfile(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, cvProfile)
	require.Equal(t, params.Name, cvProfile.Name)
	require.Equal(t, params.Email, cvProfile.Email)
	require.Equal(t, params.Phone, cvProfile.Phone)
	require.Equal(t, params.Address, cvProfile.Address)
	require.Equal(t, params.GithubUrl, cvProfile.GithubUrl)
	require.Equal(t, params.Bio, cvProfile.Bio)
	require.Equal(t, params.ProfilePicture, cvProfile.ProfilePicture)
	require.NotZero(t, cvProfile.CreatedAt)
	require.NotZero(t, cvProfile.ID)

	if params.LinkedinUrl.Valid {
		require.Equal(t, params.LinkedinUrl.String, cvProfile.LinkedinUrl.String)
	}

	return cvProfile
}

func TestQueries_CreateCvProfile(t *testing.T) {
	createRandomCvProfile(t)
}

func TestQueries_GetCvProfile(t *testing.T) {
	cvProfile := createRandomCvProfile(t)
	cvProfile2, err := testQueries.GetCvProfile(context.Background(), cvProfile.ID)

	require.NoError(t, err)
	require.NotEmpty(t, cvProfile2)
	require.Equal(t, cvProfile.ID, cvProfile2.ID)
	require.Equal(t, cvProfile.Name, cvProfile2.Name)
	require.Equal(t, cvProfile.Email, cvProfile2.Email)
	require.Equal(t, cvProfile.Phone, cvProfile2.Phone)
	require.Equal(t, cvProfile.Address, cvProfile2.Address)
	require.Equal(t, cvProfile.GithubUrl, cvProfile2.GithubUrl)
	require.Equal(t, cvProfile.Bio, cvProfile2.Bio)
	require.Equal(t, cvProfile.ProfilePicture, cvProfile2.ProfilePicture)
	require.Equal(t, cvProfile.CreatedAt, cvProfile2.CreatedAt)

	if cvProfile.LinkedinUrl.Valid {
		require.Equal(t, cvProfile.LinkedinUrl.String, cvProfile2.LinkedinUrl.String)
	}
}
