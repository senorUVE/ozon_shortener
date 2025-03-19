package url

import (
	"context"
	"ozon_shortener/internal/services/url"
)

type UrlAdapters interface {
	GenerateUrls(ctx context.Context, originalUrls []string) (map[string]string, error)
	GetOriginal(ctx context.Context, shortUrls []string) (map[string]string, error)
	PublicURL() string
}

type adapter struct {
	service url.Service
}

func New(
	service url.Service,
) UrlAdapters {
	return &adapter{
		service: service,
	}
}

func (a *adapter) GenerateUrls(ctx context.Context, originalUrls []string) (map[string]string, error) {
	return a.service.CreateURL(ctx, originalUrls)
}

func (a *adapter) GetOriginal(ctx context.Context, shortUrls []string) (map[string]string, error) {
	return a.service.GetOriginal(ctx, shortUrls)
}

func (a *adapter) PublicURL() string {
	return a.service.PublicURL()
}
