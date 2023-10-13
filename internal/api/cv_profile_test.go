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
	"testing"
	"time"
)

func TestGetCvProfileAPI(t *testing.T) {
	cvProfile := generateRandomCvProfile()

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
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCvProfile(t, recorder.Body, cvProfile)
			},
		},
		{
			name: "Invalid ID",
			id:   0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCvProfile(gomock.Any(), gomock.Any()).
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
					GetCvProfile(gomock.Any(), gomock.Eq(cvProfile.ID)).
					Times(1).
					Return(db.CvProfile{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/cv-profiles/%d", tc.id)
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
func requireBodyMatchCvProfile(t *testing.T, body *bytes.Buffer, cvProfile db.CvProfile) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotCvProfile db.CvProfile
	err = json.Unmarshal(data, &gotCvProfile)
	require.NoError(t, err)

	require.Equal(t, cvProfile.ID, gotCvProfile.ID)
	require.Equal(t, cvProfile.Name, gotCvProfile.Name)
	require.Equal(t, cvProfile.Email, gotCvProfile.Email)
	require.Equal(t, cvProfile.Phone, gotCvProfile.Phone)
	require.Equal(t, cvProfile.Address, gotCvProfile.Address)
	require.Equal(t, cvProfile.LinkedinUrl, gotCvProfile.LinkedinUrl)
	require.Equal(t, cvProfile.GithubUrl, gotCvProfile.GithubUrl)
	require.Equal(t, cvProfile.Bio, gotCvProfile.Bio)
	require.WithinDuration(t, cvProfile.CreatedAt, gotCvProfile.CreatedAt, 1*time.Second)
}
