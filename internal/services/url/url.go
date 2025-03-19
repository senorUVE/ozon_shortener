package url

import (
	"context"
	"fmt"
	"ozon_shortener/internal/repository"
	"ozon_shortener/internal/repository/entity"
	"ozon_shortener/internal/utils/bijection"
	"ozon_shortener/internal/utils/extract"
	"ozon_shortener/internal/utils/random"
	"strconv"
	"strings"
)

type Service interface {
	CreateURL(ctx context.Context, originalUrls []string) (map[string]string, error)
	GetOriginal(ctx context.Context, shortUrls []string) (map[string]string, error)
}

type service struct {
	dao    repository.DAO
	domain string
}

func New(
	dao repository.DAO,
	domain string,
) Service {
	domain = strings.TrimRight(domain, "/")
	return &service{
		dao:    dao,
		domain: domain,
	}
}

func (s *service) CreateURL(ctx context.Context, originalUrls []string) (map[string]string, error) {

	result := make(map[string]string)
	urlQuery := s.dao.NewUrlQuery(ctx)

	for _, origURL := range originalUrls {
		existing, err := urlQuery.GetUrlByOriginal(origURL)
		if err == nil {
			result[origURL] = fmt.Sprintf("%s/%s", s.domain, existing.Token)
			continue
		}

		newEntity := entity.URL{
			OriginalUrl: origURL,
		}
		if err := urlQuery.InsertUrl(newEntity); err != nil {
			return nil, fmt.Errorf("insert url: %w", err)
		}

		inserted, err := urlQuery.GetUrlByOriginal(origURL)
		if err != nil {
			return nil, fmt.Errorf("get inserted url: %w", err)
		}

		idNum, err := strconv.ParseInt(inserted.Id, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse id: %w", err)
		}

		token := bijection.ConvertNumberToKey(idNum)
		if len(token) < 10 {
			randomPartLength := 10 - len(token)
			randomPart, err := random.GenerateRandom(randomPartLength)
			if err != nil {
				return nil, fmt.Errorf("generate random string: %w", err)
			}
			token = randomPart + token
		}
		if len(token) > 10 {
			token = token[:10]
		}

		inserted.Token = token
		if err := urlQuery.UpdateURL(inserted); err != nil {
			return nil, fmt.Errorf("update url with token: %w", err)
		}

		result[origURL] = fmt.Sprintf("%s/%s", s.domain, token)
	}

	return result, nil
}

func (s *service) GetOriginal(ctx context.Context, shortUrls []string) (map[string]string, error) {
	tokens := make([]string, 0, len(shortUrls))
	for _, fullURL := range shortUrls {
		token, err := extract.ExtractToken(s.domain, fullURL)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	records, err := s.dao.NewUrlQuery(ctx).GetByTokens(tokens)
	if err != nil {
		return nil, fmt.Errorf("get urls by tokens: %w", err)
	}

	tokenToOriginal := make(map[string]string)
	for _, rec := range records {
		tokenToOriginal[rec.Token] = rec.OriginalUrl
	}
	mapping := make(map[string]string)
	for _, fullURL := range shortUrls {
		token, _ := extract.ExtractToken(s.domain, fullURL)
		if orig, ok := tokenToOriginal[token]; ok {
			mapping[fullURL] = orig
		} else {
			mapping[fullURL] = ""
		}
	}
	return mapping, nil
}
