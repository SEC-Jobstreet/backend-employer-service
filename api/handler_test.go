package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	mockdb "github.com/SEC-Jobstreet/backend-employer-service/db/mock"
	db "github.com/SEC-Jobstreet/backend-employer-service/db/sqlc"
	"github.com/SEC-Jobstreet/backend-employer-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateEnterpriseAPI(t *testing.T) {
	enterprise := RandomEnterprise()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockQuerier)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "Valid Enterprise Creation Request",
			body: gin.H{
				"id":            enterprise.ID,
				"name":          enterprise.Name.String,
				"country":       enterprise.Country.String,
				"address":       enterprise.Address.String,
				"field":         enterprise.Field.String,
				"size":          enterprise.Size.String,
				"url":           enterprise.Url.String,
				"license":       enterprise.License.String,
				"employer_id":   enterprise.EmployerID.String,
				"employer_role": enterprise.EmployerRole.String,
			},
			buildStubs: func(store *mockdb.MockQuerier) {
				arg := db.CreateEnterpriseParams{
					ID:           enterprise.ID,
					Name:         enterprise.Name,
					Country:      enterprise.Country,
					Address:      enterprise.Address,
					Field:        enterprise.Field,
					Size:         enterprise.Size,
					Url:          enterprise.Url,
					EmployerID:   enterprise.EmployerID,
					EmployerRole: enterprise.EmployerRole,
				}

				store.EXPECT().
					CreateEnterprise(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(enterprise, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchEnterprise(t, recorder.Body, enterprise)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockQuerier(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/v1/create_enterprise"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func requireBodyMatchEnterprise(t *testing.T, body *bytes.Buffer, account db.Enterprise) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotEnterprise db.Enterprise
	err = json.Unmarshal(data, &gotEnterprise)
	require.NoError(t, err)
	require.Equal(t, account, gotEnterprise)
}

func RandomEnterprise() db.Enterprise {
	return db.Enterprise{
		ID: uuid.New(),
		Name: pgtype.Text{
			String: utils.RandomString(10),
			Valid:  true,
		},
		Country: pgtype.Text{
			String: utils.RandomString(10),
			Valid:  true,
		},
		Address: pgtype.Text{
			String: utils.RandomString(10),
			Valid:  true,
		},
		Field: pgtype.Text{
			String: utils.RandomString(10),
			Valid:  true,
		},
		Size: pgtype.Text{
			String: utils.RandomString(10),
			Valid:  true,
		},
		Url: pgtype.Text{
			String: utils.RandomString(10),
			Valid:  true,
		},

		EmployerID: pgtype.Text{
			String: uuid.New().String(),
			Valid:  true,
		},
		EmployerRole: pgtype.Text{
			String: strconv.Itoa(utils.RandomInt(1, 4)),
			Valid:  true,
		},
	}
}
