package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/aalug/cv-backend-go/internal/db/mock"
	db "github.com/aalug/cv-backend-go/internal/db/sqlc"
	"github.com/aalug/cv-backend-go/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
)

const baseUrl = "/api/v1"

func TestListSkillsAPI(t *testing.T) {
	cvProfile := generateRandomCvProfile()
	skills := generateRandomSkills()

	testCases := []struct {
		name          string
		id            int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			id:   cvProfile.ID,
			buildStubs: func(store *mockdb.MockStore) {
				params := db.ListSkillsParams{
					CvProfileID: cvProfile.ID,
					Limit:       50,
					Offset:      0,
				}
				store.EXPECT().
					ListSkills(gomock.Any(), gomock.Eq(params)).
					Times(1).
					Return(skills, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchSkills(t, recorder.Body, skills)
			},
		},
		{
			name: "Invalid ID",
			id:   0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListSkills(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Not Found",
			id:   cvProfile.ID,
			buildStubs: func(store *mockdb.MockStore) {
				params := db.ListSkillsParams{
					CvProfileID: cvProfile.ID,
					Limit:       50,
					Offset:      0,
				}
				store.EXPECT().
					ListSkills(gomock.Any(), gomock.Eq(params)).
					Times(1).
					Return([]db.Skill{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			id:   cvProfile.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListSkills(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Skill{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("%s/skills/%d", baseUrl, tc.id)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(recorder)
		})
	}
}

// generateRandomSkills generates and returns a slice of random skills
func generateRandomSkills() []db.Skill {
	category := utils.RandomString(5)
	var skills []db.Skill
	for i := 0; i < 10; i++ {
		skills = append(skills, db.Skill{
			ID:            int32(i),
			Name:          utils.RandomString(6),
			Description:   utils.RandomString(10),
			Category:      category,
			Image:         utils.RandomString(6),
			HexThemeColor: utils.RandomString(6),
			CvProfileID:   1,
			Importance:    int32(i),
		})
	}

	return skills
}

// requireBodyMatchSkills asserts that the response body matches the provided skills
func requireBodyMatchSkills(t *testing.T, body *bytes.Buffer, skills []db.Skill) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotSkills []db.Skill
	err = json.Unmarshal(data, &gotSkills)
	require.NoError(t, err)

	sort.Slice(skills, func(i, j int) bool {
		return skills[i].ID < skills[j].ID
	})

	sort.Slice(gotSkills, func(i, j int) bool {
		return gotSkills[i].ID < gotSkills[j].ID
	})

	for i := 0; i < len(skills); i++ {
		require.Equal(t, skills[i].ID, gotSkills[i].ID)
		require.Equal(t, skills[i].Name, gotSkills[i].Name)
		require.Equal(t, skills[i].Description, gotSkills[i].Description)
		require.Equal(t, skills[i].Category, gotSkills[i].Category)
		require.Equal(t, skills[i].Image, gotSkills[i].Image)
		require.Equal(t, skills[i].HexThemeColor, gotSkills[i].HexThemeColor)
		require.Equal(t, skills[i].CvProfileID, gotSkills[i].CvProfileID)
		require.Equal(t, skills[i].Importance, gotSkills[i].Importance)
	}
}
