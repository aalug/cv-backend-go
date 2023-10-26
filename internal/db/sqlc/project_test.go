package db

import (
	"context"
	"github.com/aalug/cv-backend-go/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

// createRandomProject create and return a random project
func createRandomProject(t *testing.T, cvProfileID int32) Project {
	if cvProfileID == 0 {
		cvProfileID = createRandomCvProfile(t).ID
	}

	params := CreateProjectParams{
		Title:            utils.RandomString(5),
		ShortDescription: utils.RandomString(5),
		Description:      utils.RandomString(10),
		Image:            utils.RandomString(5),
		HexThemeColor:    utils.RandomString(5),
		ProjectUrl:       utils.RandomString(5),
		CvProfileID:      cvProfileID,
	}

	project, err := testQueries.CreateProject(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, project)
	require.Equal(t, params.Title, project.Title)
	require.Equal(t, params.ShortDescription, project.ShortDescription)
	require.Equal(t, params.Description, project.Description)
	require.Equal(t, params.Image, project.Image)
	require.Equal(t, params.HexThemeColor, project.HexThemeColor)
	require.Equal(t, params.ProjectUrl, project.ProjectUrl)
	require.Equal(t, params.CvProfileID, project.CvProfileID)
	require.NotZero(t, project.ID)

	return project
}

func TestQueries_CreateProject(t *testing.T) {
	createRandomProject(t, 0)
}

func TestQueries_ListProjects(t *testing.T) {
	cvProfile := createRandomCvProfile(t)
	for i := 0; i < 5; i++ {
		createRandomProject(t, cvProfile.ID)
	}

	params := ListProjectsParams{
		CvProfileID: cvProfile.ID,
		Limit:       5,
		Offset:      0,
	}

	projects, err := testQueries.ListProjects(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, projects, 5)

	for _, project := range projects {
		require.NotEmpty(t, project)
	}
}
