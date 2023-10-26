package db

import (
	"context"
	"github.com/aalug/cv-backend-go/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

// createRandomTechnology create and return a random technology
func createRandomTechnology(t *testing.T) Technology {
	params := CreateTechnologyParams{
		Name:       utils.RandomString(5),
		Url:        utils.RandomString(5),
		OrderField: utils.RandomInt(0, 100),
	}
	technology, err := testQueries.CreateTechnology(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, technology)

	require.Equal(t, params.Name, technology.Name)
	require.Equal(t, params.Url, technology.Url)
	require.Equal(t, params.OrderField, technology.OrderField)

	return technology
}

// createRandomProjectTechnology create and return a random project technology
func createRandomProjectTechnology(t *testing.T) ProjectTechnology {
	technology := createRandomTechnology(t)
	project := createRandomProject(t, 0)

	params := CreateProjectTechnologyParams{
		ProjectID:    project.ID,
		TechnologyID: technology.ID,
	}

	projectTechnology, err := testQueries.CreateProjectTechnology(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, projectTechnology)

	require.Equal(t, params.ProjectID, projectTechnology.ProjectID)
	require.Equal(t, params.TechnologyID, projectTechnology.TechnologyID)

	return projectTechnology
}

func TestQueries_CreateTechnology(t *testing.T) {
	createRandomTechnology(t)
}

func TestQueries_CreateProjectTechnology(t *testing.T) {
	createRandomProjectTechnology(t)
}

func TestQueries_ListTechnologiesForProject(t *testing.T) {
	project := createRandomProject(t, 0)

	// create 5 random technologies and connect them to the project
	for i := 0; i < 5; i++ {
		technology := createRandomTechnology(t)
		params := CreateProjectTechnologyParams{
			ProjectID:    project.ID,
			TechnologyID: technology.ID,
		}
		_, err := testQueries.CreateProjectTechnology(context.Background(), params)
		require.NoError(t, err)
	}

	technologies, err := testQueries.ListTechnologiesForProject(context.Background(), project.ID)
	require.NoError(t, err)
	require.NotEmpty(t, technologies)

	for _, technology := range technologies {
		require.NotEmpty(t, technology)
	}
}
