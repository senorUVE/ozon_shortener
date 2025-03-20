package url_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"ozon_shortener/internal/repository/entity"
	"ozon_shortener/internal/services/url"

	repoMocks "ozon_shortener/internal/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestServiceCreateURLSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockDAO := repoMocks.NewMockDAO(ctrl)
	mockQuery := repoMocks.NewMockUrlQuery(ctrl)

	mockDAO.EXPECT().NewUrlQuery(gomock.Any()).Return(mockQuery).AnyTimes()

	svc := url.New(mockDAO, "localhost:8080")

	origURL := "https://example.com"

	mockQuery.EXPECT().GetUrlByOriginal(origURL).Return(entity.URL{}, errors.New("not found")).Times(1)

	mockQuery.EXPECT().InsertUrl(gomock.Any()).Return(nil).Times(1)

	insertedRecord := entity.URL{
		Id:          "10",
		OriginalUrl: origURL,
		Token:       "",
	}
	mockQuery.EXPECT().GetUrlByOriginal(origURL).Return(insertedRecord, nil).Times(1)

	mockQuery.EXPECT().UpdateURL(gomock.Any()).DoAndReturn(
		func(u entity.URL) error {
			if u.Id != "10" || u.OriginalUrl != origURL {
				return errors.New("unexpected record")
			}
			if len(u.Token) != 10 {
				return errors.New("token length is not 10")
			}
			if u.Token[len(u.Token)-1:] != "a" {
				return errors.New("token does not end with 'a'")
			}
			return nil
		},
	).Times(1)

	result, err := svc.CreateURL(ctx, []string{origURL})
	assert.NoError(t, err)

	shortURL := result[origURL]
	prefix := "http://localhost:8080/"
	assert.True(t, strings.HasPrefix(shortURL, prefix), "shortURL should start with %s", prefix)

	token := shortURL[len(prefix):]
	assert.Equal(t, 10, len(token), "token length should be 10")
	assert.Equal(t, "a", token[len(token)-1:], "token should end with 'a'")

	const allowedRandom = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
	for i, r := range token[:9] {
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
