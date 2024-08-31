package tevianAPI

import (
	"faceScanner/internal/models"
	"fmt"
	"github.com/valyala/fasthttp"
	"image"
)

type TevianApiProvider struct {
	URL           string
	Authorization string
	Mimetype      string
}

func NewTevianProvider(url, authorization string) *TevianApiProvider {
	return &TevianApiProvider{
		URL:           url,
		Authorization: authorization,
	}
}
func (p *TevianApiProvider) ProvideRequest(image image.Image) (
	tevianApiResponse *models.TevianApiResponse, err error) {
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	p.SetHeaders(request)
	err = fasthttp.Do(request, response)
	if err != nil {
		err = fmt.Errorf("fasthttp.Do(...): %w", err)
		return nil, err
	}
	return nil, nil
}

func (p *TevianApiProvider) SetHeaders(request *fasthttp.Request) {
	request.Header.Set("Authorization", p.Authorization)
	request.Header.Set("Content-Type", p.Mimetype)
	request.Header.Set("Accept", p.Mimetype)
	return
}
