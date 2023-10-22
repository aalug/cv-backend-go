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
				params := db.ListProjectsParams{
					CvProfileID: cvProfile.ID,
					Limit:       10,
					Offset:      0,
				}
				store.EXPECT().
					ListProjects(gomock.Any(), gomock.Eq(params)).
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
					ListProjects(gomock.Any(), gomock.Any()).
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
					ListProjects(gomock.Any(), gomock.Any()).
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
					ListProjects(gomock.Any(), gomock.Any()).
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
					ListProjects(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.ListProjectsRow{}, sql.ErrNoRows)
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
					ListProjects(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.ListProjectsRow{}, sql.ErrConnDone)
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

// generateRandomProjectRows generates and returns a slice of random list project rows
func generateRandomProjectRows() []db.ListProjectsRow {
	var projects []db.ListProjectsRow
	for i := 0; i < 10; i++ {
		projects = append(projects, db.ListProjectsRow{
			ID:               int32(i),
			Title:            utils.RandomString(6),
			ShortDescription: utils.RandomString(5),
			Image:            utils.RandomString(6),
			TechnologiesUsed: []string{
				utils.RandomString(4),
				utils.RandomString(5),
			},
			Description:   utils.RandomString(10),
			HexThemeColor: utils.RandomString(6),
			ProjectUrl:    utils.RandomString(6),
		})
	}

	return projects
}

// requireBodyMatchProjects asserts that the response body matches the provided projects
func requireBodyMatchProjects(t *testing.T, body *bytes.Buffer, projects []db.ListProjectsRow) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotProjects []db.ListProjectsRow
	err = json.Unmarshal(data, &gotProjects)
	require.NoError(t, err)

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].ID < projects[j].ID
	})

	sort.Slice(gotProjects, func(i, j int) bool {
		return gotProjects[i].ID < gotProjects[j].ID
	})

	for i := 0; i < len(projects); i++ {
		// sort technologies used
		sort.Strings(projects[i].TechnologiesUsed)
		sort.Strings(gotProjects[i].TechnologiesUsed)

		require.Equal(t, projects[i].ID, gotProjects[i].ID)
		require.Equal(t, projects[i].Title, gotProjects[i].Title)
		require.Equal(t, projects[i].ShortDescription, gotProjects[i].ShortDescription)
		require.Equal(t, projects[i].Description, gotProjects[i].Description)
		require.Equal(t, projects[i].Image, gotProjects[i].Image)
		require.Equal(t, projects[i].TechnologiesUsed, gotProjects[i].TechnologiesUsed)
		require.Equal(t, projects[i].HexThemeColor, gotProjects[i].HexThemeColor)
		require.Equal(t, projects[i].ProjectUrl, gotProjects[i].ProjectUrl)
	}
}
