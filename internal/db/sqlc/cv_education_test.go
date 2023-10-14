package db

import (
	"context"
	"github.com/aalug/cv-backend-go/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// createRandomCvEducation create and return a random cv education
func createRandomCvEducation(t *testing.T, cvProfileID int32) CvEducation {
	if cvProfileID == 0 {
		cvProfileID = createRandomCvProfile(t).ID
	}
	params := CreateCvEducationParams{
		Institution: utils.RandomString(5),
		Degree:      utils.RandomString(5),
		StartDate:   time.Now(),
		EndDate:     time.Now(),
		CvProfileID: cvProfileID,
	}

	cvEducation, err := testQueries.CreateCvEducation(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, cvEducation)
	require.Equal(t, params.Institution, cvEducation.Institution)
	require.Equal(t, params.Degree, cvEducation.Degree)
	require.Equal(t, params.StartDate.Format("2006-01-02"), cvEducation.StartDate.Format("2006-01-02"))
	require.Equal(t, params.EndDate.Format("2006-01-02"), cvEducation.EndDate.Format("2006-01-02"))
	require.Equal(t, params.CvProfileID, cvEducation.CvProfileID)
	require.NotZero(t, cvEducation.ID)

	return cvEducation
}

func TestQueries_CreateCvEducation(t *testing.T) {
	createRandomCvEducation(t, 0)
}

func TestQueries_GetCvEducation(t *testing.T) {
	cvEducation := createRandomCvEducation(t, 0)
	cvEducation2, err := testQueries.GetCvEducation(context.Background(), cvEducation.ID)
	require.NoError(t, err)
	require.NotEmpty(t, cvEducation2)
	require.Equal(t, cvEducation.ID, cvEducation2.ID)
	require.Equal(t, cvEducation.Institution, cvEducation2.Institution)
	require.Equal(t, cvEducation.Degree, cvEducation2.Degree)
	require.WithinDuration(t, cvEducation.StartDate, cvEducation2.StartDate, 24*time.Hour)
	require.WithinDuration(t, cvEducation.EndDate, cvEducation2.EndDate, 24*time.Hour)
	require.Equal(t, cvEducation.StartDate.Format("2006-01-02"), cvEducation2.StartDate.Format("2006-01-02"))
	require.Equal(t, cvEducation.EndDate.Format("2006-01-02"), cvEducation2.EndDate.Format("2006-01-02"))
	require.Equal(t, cvEducation.CvProfileID, cvEducation2.CvProfileID)
}

func TestQueries_ListCvEducations(t *testing.T) {
	cvProfile := createRandomCvProfile(t)
	cvProfileID := cvProfile.ID
	for i := 0; i < 5; i++ {
		createRandomCvEducation(t, cvProfileID)
	}

	params := ListCvEducationsParams{
		CvProfileID: cvProfileID,
		Limit:       5,
		Offset:      0,
	}

	cvEducations, err := testQueries.ListCvEducations(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, cvEducations, 5)

	for _, cvEducation := range cvEducations {
		require.NotEmpty(t, cvEducation)
		require.Equal(t, cvProfileID, cvEducation.CvProfileID)
	}
}
