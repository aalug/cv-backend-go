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

func TestListProjectsAPI(t *testing.T) {
	cvProfile := generateRandomCvProfile()
	projects := generateRandomProjectRows()

	type Query struct {
		page     int32
		pageSize int32
	}

	testCases := []struct {
		name          string
		id            int32
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			id:   cvProfile.ID,
			query: Query{
				page:     1,
				pageSize: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				params := db.ListProjectsWithTechnologiesParams{
					CvProfileID: cvProfile.ID,
					Limit:       10,
					Offset:      0,
				}
				store.EXPECT().
					ListProjectsWithTechnologies(gomock.Any(), gomock.Eq(params)).
					Times(1).
					Return(projects, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchProjects(t, recorder.Body, projects)
			},
		},
		{
			name: "Invalid ID",
			id:   0,
			query: Query{
				page:     1,
				pageSize: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListProjectsWithTechnologies(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Page",
			id:   cvProfile.ID,
			query: Query{
				page:     0,
				pageSize: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListProjectsWithTechnologies(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Page Size",
			id:   cvProfile.ID,
			query: Query{
				page:     1,
				pageSize: 30,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListProjectsWithTechnologies(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Not Found",
			id:   cvProfile.ID,
			query: Query{
				page:     1,
				pageSize: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListProjectsWithTechnologies(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.ListProjectsWithTechnologiesRow{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			id:   cvProfile.ID,
			query: Query{
				page:     1,
				pageSize: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListProjectsWithTechnologies(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.ListProjectsWithTechnologiesRow{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("%s/projects/%d", baseUrl, tc.id)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query params
			q := req.URL.Query()
			q.Add("page", fmt.Sprintf("%d", tc.query.page))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			req.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(recorder)
		})
	}
}

func TestListProjectsBySkillNameAPI(t *testing.T) {
	cvProfile := generateRandomCvProfile()
	projects := generateRandomProjectBySkillNameRows()
	skillName := utils.RandomString(5)

	type Query struct {
		page     int32
		pageSize int32
	}

	testCases := []struct {
		name          string
		id            int32
		skillName     string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			id:        cvProfile.ID,
			skillName: skillName,
			query: Query{
				page:     1,
				pageSize: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				params := db.ListProjectsWithTechnologiesBySkillNameParams{
					SkillName:   skillName,
					CvProfileID: cvProfile.ID,
					Limit:       10,
					Offset:      0,
				}
				store.EXPECT().
					ListProjectsWithTechnologiesBySkillName(gomock.Any(), gomock.Eq(params)).
					Times(1).
					Return(projects, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchProjects(t, recorder.Body, projects)
			},
		},
		{
			name:      "Invalid ID",
			id:        0,
			skillName: skillName,
			query: Query{
				page:     1,
				pageSize: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListProjectsWithTechnologiesBySkillName(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Invalid Skill Name",
			id:        cvProfile.ID,
			skillName: "123",
			query: Query{
				page:     1,
				pageSize: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListProjectsWithTechnologiesBySkillName(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Invalid Page",
			id:        cvProfile.ID,
			skillName: skillName,
			query: Query{
				page:     0,
				pageSize: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListProjectsWithTechnologiesBySkillName(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Invalid Page Size",
			id:        cvProfile.ID,
			skillName: skillName,
			query: Query{
				page:     1,
				pageSize: 50,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListProjectsWithTechnologiesBySkillName(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Not Found",
			id:        cvProfile.ID,
			skillName: skillName,
			query: Query{
				page:     1,
				pageSize: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListProjectsWithTechnologiesBySkillName(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.ListProjectsWithTechnologiesBySkillNameRow{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Internal Server Error",
			id:        cvProfile.ID,
			skillName: skillName,
			query: Query{
				page:     1,
				pageSize: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListProjectsWithTechnologiesBySkillName(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.ListProjectsWithTechnologiesBySkillNameRow{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("%s/projects/skill/%d/%s", baseUrl, tc.id, tc.skillName)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query params
			q := req.URL.Query()
			q.Add("page", fmt.Sprintf("%d", tc.query.page))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			req.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(recorder)
		})
	}
}

// generateRandomProjectRows generates and returns a slice of random list project rows
func generateRandomProjectRows() []db.ListProjectsWithTechnologiesRow {
	var projects []db.ListProjectsWithTechnologiesRow
	var technologiesUsed []db.ListTechnologiesForProjectRow

	for i := 0; i < 5; i++ {
		technologiesUsed = append(technologiesUsed, db.ListTechnologiesForProjectRow{
			ID:   int32(i),
			Name: utils.RandomString(6),
			Url:  utils.RandomString(5),
		})
	}

	for i := 0; i < 10; i++ {
		projects = append(projects, db.ListProjectsWithTechnologiesRow{
			ID:               int32(i),
			Title:            utils.RandomString(6),
			ShortDescription: utils.RandomString(5),
			Description:      utils.RandomString(10),
			Image:            utils.RandomString(6),
			HexThemeColor:    utils.RandomString(6),
			ProjectUrl:       utils.RandomString(6),
			Significance:     utils.RandomInt(1, 50),
			TechnologiesUsed: technologiesUsed,
		})
	}

	return projects
}

// generateRandomProjectBySkillNameRows generates and returns a slice of random list project by skill name rows
func generateRandomProjectBySkillNameRows() []db.ListProjectsWithTechnologiesBySkillNameRow {
	var projects []db.ListProjectsWithTechnologiesBySkillNameRow
	var technologiesUsed []db.ListTechnologiesForProjectRow

	for i := 0; i < 5; i++ {
		technologiesUsed = append(technologiesUsed, db.ListTechnologiesForProjectRow{
			ID:   int32(i),
			Name: utils.RandomString(6),
			Url:  utils.RandomString(5),
		})
	}

	for i := 0; i < 10; i++ {
		projects = append(projects, db.ListProjectsWithTechnologiesBySkillNameRow{
			ID:               int32(i),
			Title:            utils.RandomString(6),
			ShortDescription: utils.RandomString(5),
			Description:      utils.RandomString(10),
			Image:            utils.RandomString(6),
			HexThemeColor:    utils.RandomString(6),
			ProjectUrl:       utils.RandomString(6),
			Significance:     utils.RandomInt(1, 50),
			TechnologiesUsed: technologiesUsed,
		})
	}

	return projects
}

// requireBodyMatchProjects asserts that the response body matches the provided projects
func requireBodyMatchProjects(t *testing.T, body *bytes.Buffer, projects interface{}) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	switch p := projects.(type) {
	case []db.ListProjectsWithTechnologiesRow:
		var gotProjects []db.ListProjectsWithTechnologiesRow
		err = json.Unmarshal(data, &gotProjects)
		require.NoError(t, err)

		// sort projects to compare each result
		sort.Slice(p, func(i, j int) bool {
			return p[i].ID < p[j].ID
		})

		sort.Slice(gotProjects, func(i, j int) bool {
			return gotProjects[i].ID < gotProjects[j].ID
		})

		for i := 0; i < len(p); i++ {
			project := p[i]
			gotProject := gotProjects[i]

			require.Equal(t, project.ID, gotProject.ID)
			require.Equal(t, project.Title, gotProject.Title)
			require.Equal(t, project.ShortDescription, gotProject.ShortDescription)
			require.Equal(t, project.Description, gotProject.Description)
			require.Equal(t, project.Image, gotProject.Image)
			require.Equal(t, project.HexThemeColor, gotProject.HexThemeColor)
			require.Equal(t, project.ProjectUrl, gotProject.ProjectUrl)

			for j := 0; j < len(project.TechnologiesUsed); j++ {
				require.Equal(t, project.TechnologiesUsed[j].ID, gotProject.TechnologiesUsed[j].ID)
				require.Equal(t, project.TechnologiesUsed[j].Name, gotProject.TechnologiesUsed[j].Name)
				require.Equal(t, project.TechnologiesUsed[j].Url, gotProject.TechnologiesUsed[j].Url)
			}
		}
	case []db.ListProjectsWithTechnologiesBySkillNameRow:
		var gotProjects []db.ListProjectsWithTechnologiesBySkillNameRow
		err = json.Unmarshal(data, &gotProjects)
		require.NoError(t, err)

		// sort projects to compare each result
		sort.Slice(p, func(i, j int) bool {
			return p[i].ID < p[j].ID
		})

		sort.Slice(gotProjects, func(i, j int) bool {
			return gotProjects[i].ID < gotProjects[j].ID
		})

		for i := 0; i < len(p); i++ {
			project := p[i]
			gotProject := gotProjects[i]

			require.Equal(t, project.ID, gotProject.ID)
			require.Equal(t, project.Title, gotProject.Title)
			require.Equal(t, project.ShortDescription, gotProject.ShortDescription)
			require.Equal(t, project.Description, gotProject.Description)
			require.Equal(t, project.Image, gotProject.Image)
			require.Equal(t, project.HexThemeColor, gotProject.HexThemeColor)
			require.Equal(t, project.ProjectUrl, gotProject.ProjectUrl)

			for j := 0; j < len(project.TechnologiesUsed); j++ {
				require.Equal(t, project.TechnologiesUsed[j].ID, gotProject.TechnologiesUsed[j].ID)
				require.Equal(t, project.TechnologiesUsed[j].Name, gotProject.TechnologiesUsed[j].Name)
				require.Equal(t, project.TechnologiesUsed[j].Url, gotProject.TechnologiesUsed[j].Url)
			}
		}
	default:
		t.Fatalf("unsupported type %T", projects)
	}
}
