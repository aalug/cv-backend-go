package db

import (
	"context"
	"github.com/aalug/cv-backend-go/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

// createRandomProject create and return a random project
func createRandomProject(t *testing.T) Project {
	cvProfile := createRandomCvProfile(t)

	technologiesUsed := []string{
		utils.RandomString(2),
		utils.RandomString(3),
	}

	params := CreateProjectParams{
		Title:            utils.RandomString(5),
		ShortDescription: utils.RandomString(5),
		Description:      utils.RandomString(10),
		Image:            utils.RandomString(5),
		TechnologiesUsed: technologiesUsed,
		HexThemeColor:    utils.RandomString(5),
		ProjectUrl:       utils.RandomString(5),
		CvProfileID:      cvProfile.ID,
	}

	project, err := testQueries.CreateProject(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, project)
	require.Equal(t, params.Title, project.Title)
	require.Equal(t, params.ShortDescription, project.ShortDescription)
	require.Equal(t, params.Description, project.Description)
	require.Equal(t, params.Image, project.Image)
	require.Equal(t, params.TechnologiesUsed, project.TechnologiesUsed)
	require.Equal(t, params.HexThemeColor, project.HexThemeColor)
	require.Equal(t, params.ProjectUrl, project.ProjectUrl)
	require.Equal(t, params.CvProfileID, project.CvProfileID)
	require.NotZero(t, project.ID)

	return project
}

func TestQueries_CreateProject(t *testing.T) {
	createRandomProject(t)
}

func TestQueries_GetProject(t *testing.T) {
	project := createRandomProject(t)
	project2, err := testQueries.GetProject(context.Background(), project.ID)
	require.NoError(t, err)
	require.NotEmpty(t, project2)
	require.Equal(t, project.ID, project2.ID)
	require.Equal(t, project.Title, project2.Title)
	require.Equal(t, project.ShortDescription, project2.ShortDescription)
	require.Equal(t, project.Description, project2.Description)
	require.Equal(t, project.Image, project2.Image)
	require.Equal(t, project.TechnologiesUsed, project2.TechnologiesUsed)
	require.Equal(t, project.HexThemeColor, project2.HexThemeColor)
	require.Equal(t, project.ProjectUrl, project2.ProjectUrl)
	require.Equal(t, project.CvProfileID, project2.CvProfileID)
}

func TestQueries_ListProjects(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomProject(t)
	}

	params := ListProjectsParams{
		Limit:  5,
		Offset: 0,
	}

	projects, err := testQueries.ListProjects(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, projects, 5)

	for _, project := range projects {
		require.NotEmpty(t, project)
	}
}
