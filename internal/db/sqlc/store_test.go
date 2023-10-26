package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSQLStore_ListProjectsWithTechnologies(t *testing.T) {
	store := NewStore(testDB)

	// create test data in the database
	cvProfile := createRandomCvProfile(t)
	project := createRandomProject(t, cvProfile.ID)

	for i := 0; i < 5; i++ {
		technology := createRandomTechnology(t)
		p := CreateProjectTechnologyParams{
			ProjectID:    project.ID,
			TechnologyID: technology.ID,
		}
		_, err := store.CreateProjectTechnology(context.Background(), p)
		require.NoError(t, err)
	}

	params := ListProjectsWithTechnologiesParams{
		CvProfileID: cvProfile.ID,
		Limit:       5,
		Offset:      0,
	}

	projects, err := store.ListProjectsWithTechnologies(context.Background(), params)
	require.NoError(t, err)

	require.NotEmpty(t, projects[0])
	require.Len(t, projects[0].TechnologiesUsed, 5)
	require.Equal(t, projects[0].ID, project.ID)
}
