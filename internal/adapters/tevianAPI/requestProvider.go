package tevianAPI

import (
	"context"
	"encoding/json"
	"faceScanner/internal/models"
	"fmt"
	"github.com/valyala/fasthttp"
	"golang.org/x/time/rate"
	"net/url"
)

const (
	OrientationClassifier = "orientation_classifier"
	RotateUntilFacesFound = "rotate_until_faces_found"
	Demographics          = "demographics"
)

type TevianApiProvider struct {
	URL           string
	Authorization string
	Mimetype      string
	RateLimiter   *rate.Limiter
}

func NewTevianProvider(url, authorization, mimetype string) *TevianApiProvider {
	return &TevianApiProvider{
		URL:           url,
		Authorization: authorization,
		Mimetype:      mimetype,
		RateLimiter:   rate.NewLimiter(1, 1),
	}
}

func (p *TevianApiProvider) ProvideRequest(image []byte) (tevianApiResponse models.TevianApiResponse, err error) {
	var (
		request  = fasthttp.AcquireRequest()
		response = fasthttp.AcquireResponse()
	)

	requestURL, err := p.GetURL()
	if err != nil {
		err = fmt.Errorf("p.GetURL(): %w", err)
		return tevianApiResponse, err
	}

	defer fasthttp.ReleaseRequest(request)

	request.Header.Set("Authorization", p.Authorization)
	request.Header.Set("Accept", p.Mimetype)
	request.Header.SetMethod("POST")
	request.Header.SetContentType(p.Mimetype)
	request.SetBody(image)
	request.SetRequestURI(requestURL)

	err = p.RateLimiter.Wait(context.Background())
	if err != nil {
		err = fmt.Errorf("p.RateLimiter.Wait(...): %w", err)
		return tevianApiResponse, err
	}

	client := &fasthttp.Client{}
	if err = client.Do(request, response); err != nil {
		err = fmt.Errorf("client.Do(...): %w", err)
		return tevianApiResponse, err
	}

	if response.StatusCode() != 200 {
		err = fmt.Errorf("response.StatusCode() != 200")
		return tevianApiResponse, err
	}

	err = json.Unmarshal(response.Body(), &tevianApiResponse)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal(...): %w", err)
		return tevianApiResponse, err
	}
	tevianApiResponse.BodyRaw = string(response.Body())
	return tevianApiResponse, err
}

func (p *TevianApiProvider) GetURL() (string, error) {
	values := url.Values{}
	values.Add(OrientationClassifier, "true")
	values.Add(RotateUntilFacesFound, "true")
	values.Add(Demographics, "true")

	parsedURL, err := url.Parse(p.URL)
	if err != nil {
		err = fmt.Errorf("url.Parse(...): %w", err)
		return "", err
	}
	parsedURL.RawQuery = values.Encode()

	return parsedURL.String(), nil
}

func (p *TevianApiProvider) GetResponse(body []byte) (tevianApiResponse models.TevianApiResponse, err error) {

	err = json.Unmarshal(body, &tevianApiResponse)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal(...): %w", err)
		return tevianApiResponse, err
	}
	tevianApiResponse.BodyRaw = string(body)
	return tevianApiResponse, err
}
