package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

// createTestProjectSkill create and return project skill object for testing purposes
func createTestProjectSkill(t *testing.T, projectID, skillID int32) ProjectSkill {
	params := CreateProjectSkillParams{
		ProjectID: projectID,
		SkillID:   skillID,
	}
	projectSkill, err := testQueries.CreateProjectSkill(context.Background(), params)
	require.NoError(t, err)

	require.NotEmpty(t, projectSkill)
	require.Equal(t, projectID, projectSkill.ProjectID)
	require.Equal(t, skillID, projectSkill.SkillID)

	return projectSkill
}

func TestQueries_CreateProjectSkill(t *testing.T) {
	project := createRandomProject(t, 0)
	skill := createRandomSkill(t, project.CvProfileID)

	createTestProjectSkill(t, project.ID, skill.ID)
}
