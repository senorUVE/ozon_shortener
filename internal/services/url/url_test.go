package url_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"ozon_shortener/internal/repository/entity"
	"ozon_shortener/internal/services/url"
	"ozon_shortener/internal/utils/bijection"
	"ozon_shortener/internal/utils/random"

	repoMocks "ozon_shortener/internal/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_CreateURL_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockDAO := repoMocks.NewMockDAO(ctrl)
	mockQuery := repoMocks.NewMockUrlQuery(ctrl)

	mockDAO.EXPECT().NewUrlQuery(gomock.Any()).Return(mockQuery).AnyTimes()

	svc := url.New(mockDAO, "localhost:8080")

	origURL := "https://example.com"

	mockQuery.EXPECT().GetUrlByOriginal(origURL).Return(entity.URL{}, errors.New("not found")).Times(1)

	newEntity := entity.URL{
		OriginalUrl: origURL,
		Token:       "",
	}
	insertedRecord := entity.URL{
		Id:          "10",
		OriginalUrl: origURL,
		Token:       "",
	}
	mockQuery.EXPECT().InsertUrlReturning(newEntity).Return(insertedRecord, nil).Times(1)

	origRandRead := random.RandRead
	random.RandRead = func(b []byte) (int, error) {
		if len(b) == 9 {
			copy(b, []byte("123456789"))
			return 9, nil
		}
		return 0, errors.New("unexpected length")
	}
	defer func() { random.RandRead = origRandRead }()

	var capturedToken string
	mockQuery.EXPECT().UpdateURL(gomock.Any()).DoAndReturn(
		func(u entity.URL) error {
			capturedToken = u.Token
			if u.Id != "10" || u.OriginalUrl != origURL {
				return errors.New("unexpected record")
			}
			return nil
		},
	).Times(1)

	result, err := svc.CreateURL(ctx, []string{origURL})
	assert.NoError(t, err)

	expectedShortURL := "http://localhost:8080/" + capturedToken
	assert.Equal(t, expectedShortURL, result[origURL])

	assert.Equal(t, 10, len(capturedToken))
	tokenPart := bijection.ConvertNumberToKey(10)
	assert.Equal(t, tokenPart, capturedToken[len(capturedToken)-len(tokenPart):])

	const allowedRandom = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
	for i, r := range capturedToken[:9] {
		assert.True(t, strings.ContainsRune(allowedRandom, r), "character %q at index %d is not allowed", r, i)
	}
}

func TestServiceGetOriginalSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockDAO := repoMocks.NewMockDAO(ctrl)
	mockQuery := repoMocks.NewMockUrlQuery(ctrl)

	mockDAO.EXPECT().NewUrlQuery(gomock.Any()).Return(mockQuery).AnyTimes()

	svc := url.New(mockDAO, "localhost:8080")

	shortURL := "http://localhost:8080/123456789a"
	expectedRecords := []entity.URL{
		{
			Id:          "10",
			OriginalUrl: "https://example.com",
			Token:       "123456789a",
		},
	}
	mockQuery.EXPECT().GetByTokens([]string{"123456789a"}).Return(expectedRecords, nil)

	res, err := svc.GetOriginal(ctx, []string{shortURL})
	assert.NoError(t, err)
	expectedMapping := map[string]string{
		shortURL: "https://example.com",
	}
	assert.Equal(t, expectedMapping, res)
}
