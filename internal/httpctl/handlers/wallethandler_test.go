package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	mock_handlers "walletapp/internal/httpctl/handlers/mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateHappyPath(t *testing.T) {
	t.Parallel()

	// nolint:revive // Cannot make it constant.
	req := UpdateWalletBalanceRequest{
		WalletID:      uuid.New(),
		OperationType: "WITHDRAW",
		Amount:        500,
	}

	w := httptest.NewRecorder()
	gctx, _ := gin.CreateTestContext(w)

	body, err := json.Marshal(req)
	require.NoError(t, err)

	buf := bytes.NewBuffer(body)

	reqbody := io.NopCloser(buf)

	//nolint:exhaustruct // Testcase doesnt values other fields.
	httpreq := &http.Request{
		URL:    &url.URL{}, //nolint:exhaustruct // Testcase doesnt values other fields.
		Header: make(http.Header),
		Body:   reqbody,
	}

	gctx.Request = httpreq

	l := zerolog.New(zerolog.NewTestWriter(t))
	mockucase := mock_handlers.NewMockWalletUsecase(t)
	_ = mockucase.
		On("UpdateWalletBalance", mock.Anything, mock.AnythingOfType("int")).
		Return(error(nil))

	handler := NewWalletHandler(&l, mockucase)

	handler.Update(gctx)

	rsp := w.Result()
	defer rsp.Body.Close()

	compareResult := func() bool {
		bodybytes, err := io.ReadAll(rsp.Body)
		if !assert.NoError(t, err) {

			return false
		}

		//nolint:exhaustruct // Used in unmarshalling.
		rsp := UpdateWalletBalanceResponse{}

		err = json.Unmarshal(bodybytes, &rsp)

		require.NoError(t, err)

		return assert.Equal(t, MessageSuccess, rsp.Msg)

	}

	assert.True(t, compareResult())

}

func TestGetHappyPath(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	gctx, _ := gin.CreateTestContext(w)

	//nolint:exhaustruct // Testcase doesnt values other fields.
	httpreq := &http.Request{
		URL:    &url.URL{}, //nolint:exhaustruct // Testcase doesnt values other fields.
		Header: make(http.Header),
	}

	gctx.Request = httpreq

	gctx.Params = append(gctx.Params, gin.Param{
		Key:   "wallet_uuid",
		Value: uuid.NewString(),
	})

	const expectedWalletBalance = 123

	l := zerolog.New(zerolog.NewTestWriter(t))
	mockucase := mock_handlers.NewMockWalletUsecase(t)
	_ = mockucase.
		On("GetWalletBalance", mock.Anything).
		Return(expectedWalletBalance, error(nil))

	handler := NewWalletHandler(&l, mockucase)

	handler.GetBalance(gctx)

	rsp := w.Result()
	defer rsp.Body.Close()

	compareResult := func() bool {
		bodybytes, err := io.ReadAll(rsp.Body)
		if !assert.NoError(t, err) {

			return false
		}

		//nolint:exhaustruct // Used in unmarshalling.
		rsp := GetWalletBalanceResponse{}

		err = json.Unmarshal(bodybytes, &rsp)

		require.NoError(t, err)

		return assert.Equal(t, expectedWalletBalance, rsp.Balance)
	}

	assert.True(t, compareResult())

}
