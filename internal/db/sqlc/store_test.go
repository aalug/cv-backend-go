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

func TestSQLStore_ListProjectsWithTechnologiesBySkillName(t *testing.T) {
	store := NewStore(testDB)

	// create test data in the database
	cvProfile := createRandomCvProfile(t)
	skill := createRandomSkill(t, cvProfile.ID)

	for i := 0; i < 5; i++ {
		project := createRandomProject(t, cvProfile.ID)

		params := CreateProjectSkillParams{
			ProjectID: project.ID,
			SkillID:   skill.ID,
		}
		_, err := store.CreateProjectSkill(context.Background(), params)
		require.NoError(t, err)
	}

	params := ListProjectsWithTechnologiesBySkillNameParams{
		CvProfileID: cvProfile.ID,
		SkillName:   skill.Name,
		Limit:       5,
		Offset:      0,
	}

	projects, err := store.ListProjectsWithTechnologiesBySkillName(context.Background(), params)
	require.NoError(t, err)

	require.Len(t, projects, 5)
	for _, project := range projects {
		require.NotEmpty(t, project)
	}

	p := ListProjectsBySkillNameParams{
		CvProfileID: cvProfile.ID,
		Limit:       5,
		Offset:      0,
		SkillName:   skill.Name,
	}

	x, err := store.ListProjectsBySkillName(context.Background(), p)
	require.NoError(t, err)
	require.Len(t, x, 5)
	for i := 0; i < 5; i++ {
		require.Equal(t, x[i].ID, projects[i].ID)
	}
}
