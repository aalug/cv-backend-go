package db

import (
	"context"
	"github.com/aalug/cv-backend-go/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

// createRandomSkill create and return a random project
func createRandomSkill(t *testing.T, cvProfileID int32) Skill {
	if cvProfileID == 0 {
		cvProfileID = createRandomCvProfile(t).ID
	}

	params := CreateSkillParams{
		Name:          utils.RandomString(7),
		Description:   utils.RandomString(10),
		Category:      utils.RandomString(7),
		Importance:    utils.RandomInt(1, 100),
		Image:         utils.RandomString(5),
		HexThemeColor: utils.RandomString(6),
		CvProfileID:   cvProfileID,
	}
	skill, err := testQueries.CreateSkill(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, skill)
	require.Equal(t, params.Name, skill.Name)
	require.Equal(t, params.Description, skill.Description)
	require.Equal(t, params.Category, skill.Category)
	require.Equal(t, params.Importance, skill.Importance)
	require.Equal(t, params.Image, skill.Image)
	require.Equal(t, params.HexThemeColor, skill.HexThemeColor)
	require.Equal(t, params.CvProfileID, skill.CvProfileID)

	return skill
}

func TestQueries_CreateSkill(t *testing.T) {
	createRandomSkill(t, 0)
}

func TestQueries_GetSkill(t *testing.T) {
	skill := createRandomSkill(t, 0)
	skill2, err := testQueries.GetSkill(context.Background(), skill.ID)
	require.NoError(t, err)
	require.NotEmpty(t, skill2)
	require.Equal(t, skill.Name, skill2.Name)
	require.Equal(t, skill.Description, skill2.Description)
	require.Equal(t, skill.Category, skill2.Category)
	require.Equal(t, skill.Importance, skill2.Importance)
	require.Equal(t, skill.Image, skill2.Image)
	require.Equal(t, skill.HexThemeColor, skill2.HexThemeColor)
	require.Equal(t, skill.CvProfileID, skill2.CvProfileID)
}

func TestQueries_ListSkills(t *testing.T) {
	cvProfile := createRandomCvProfile(t)
	cvProfileID := cvProfile.ID
	for i := 0; i < 5; i++ {
		createRandomSkill(t, cvProfileID)
	}
	params := ListSkillsParams{
		CvProfileID: cvProfileID,
		Limit:       5,
		Offset:      0,
	}

	skills, err := testQueries.ListSkills(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, skills, 5)

	for _, skill := range skills {
		require.Equal(t, cvProfileID, skill.CvProfileID)
		require.NotEmpty(t, skill)
	}
}
