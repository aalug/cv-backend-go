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
	"time"
)

func TestGetCvProfileAPI(t *testing.T) {
	cvProfile := generateRandomCvProfile()
	education := []db.CvEducation{
		{
			ID:          utils.RandomInt(1, 1000),
			Institution: utils.RandomString(5),
			Degree:      utils.RandomString(5),
			StartDate:   time.Now().Add(-time.Hour * 24 * 3 * 365),
			EndDate:     time.Now(),
			CvProfileID: cvProfile.ID,
		},
	}

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
				store.EXPECT().
					GetCvProfile(gomock.Any(), gomock.Eq(cvProfile.ID)).
					Times(1).
					Return(cvProfile, nil)
				params := db.ListCvEducationsParams{
					CvProfileID: cvProfile.ID,
					Limit:       5,
					Offset:      0,
				}
				store.EXPECT().
					ListCvEducations(gomock.Any(), gomock.Eq(params)).
					Times(1).
					Return(education, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCvProfile(t, recorder.Body, cvProfile, education)
			},
		},
		{
			name: "Invalid ID",
			id:   0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCvProfile(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					ListCvEducations(gomock.Any(), gomock.Any()).
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
				store.EXPECT().
					GetCvProfile(gomock.Any(), gomock.Eq(cvProfile.ID)).
					Times(1).
					Return(db.CvProfile{}, sql.ErrNoRows)
				store.EXPECT().
					ListCvEducations(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Internal Server Error GetCvProfile",
			id:   cvProfile.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCvProfile(gomock.Any(), gomock.Eq(cvProfile.ID)).
					Times(1).
					Return(db.CvProfile{}, sql.ErrConnDone)
				store.EXPECT().
					ListCvEducations(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Internal Server Error ListCvEducations",
			id:   cvProfile.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCvProfile(gomock.Any(), gomock.Eq(cvProfile.ID)).
					Times(1).
					Return(cvProfile, nil)
				store.EXPECT().
					ListCvEducations(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.CvEducation{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("%s/cv-profiles/%d", baseUrl, tc.id)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(recorder)
		})
	}
}

// generateRandomCvProfile generates and returns a random cv profile
func generateRandomCvProfile() db.CvProfile {
	return db.CvProfile{
		ID:      utils.RandomInt(1, 1000),
		Name:    utils.RandomString(5),
		Email:   utils.RandomEmail(),
		Phone:   utils.RandomString(9),
		Address: utils.RandomString(6),
		LinkedinUrl: sql.NullString{
			String: utils.RandomString(10),
			Valid:  true,
		},
		GithubUrl:      utils.RandomString(10),
		Bio:            utils.RandomString(8),
		CreatedAt:      time.Now(),
		ProfilePicture: utils.RandomString(5),
	}
}

// requireBodyMatchCvProfile asserts that the response body matches the provided cv profile
func requireBodyMatchCvProfile(t *testing.T, body *bytes.Buffer, cvProfile db.CvProfile, education []db.CvEducation) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotCvProfile getCvProfileResponse
	err = json.Unmarshal(data, &gotCvProfile)
	require.NoError(t, err)

	require.Equal(t, cvProfile.Name, gotCvProfile.Name)
	require.Equal(t, cvProfile.Email, gotCvProfile.Email)
	require.Equal(t, cvProfile.Phone, gotCvProfile.Phone)
	require.Equal(t, cvProfile.Address, gotCvProfile.Address)
	require.Equal(t, cvProfile.GithubUrl, gotCvProfile.GithubUrl)
	require.Equal(t, cvProfile.Bio, gotCvProfile.Bio)
	require.Equal(t, cvProfile.ProfilePicture, gotCvProfile.ProfilePicture)

	sort.Slice(education, func(i, j int) bool {
		return education[i].ID < education[j].ID
	})

	sort.Slice(gotCvProfile.Education, func(i, j int) bool {
		return gotCvProfile.Education[i].ID < gotCvProfile.Education[j].ID
	})

	for i := 0; i < len(education); i++ {
		require.Equal(t, education[i].ID, gotCvProfile.Education[i].ID)
		require.Equal(t, education[i].Institution, gotCvProfile.Education[i].Institution)
		require.Equal(t, education[i].Degree, gotCvProfile.Education[i].Degree)
		require.WithinDuration(t, education[i].StartDate, gotCvProfile.Education[i].StartDate, 1*time.Second)
		require.WithinDuration(t, education[i].EndDate, gotCvProfile.Education[i].EndDate, 1*time.Second)
		require.Equal(t, education[i].CvProfileID, gotCvProfile.Education[i].CvProfileID)
	}

}
