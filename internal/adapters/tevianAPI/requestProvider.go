package tevianAPI

import (
	"encoding/json"
	"faceScanner/internal/models"
	"fmt"
	"github.com/valyala/fasthttp"
)

type TevianApiProvider struct {
	URL           string
	Authorization string
	Mimetype      string
}

func NewTevianProvider(url, authorization, mimetype string) *TevianApiProvider {
	return &TevianApiProvider{
		URL:           url,
		Authorization: authorization,
		Mimetype:      mimetype,
	}
}
func (p *TevianApiProvider) ProvideRequest(image []byte) (tevianApiResponse models.TevianApiResponse, err error) {

	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(request)

	request.Header.Set("Authorization", p.Authorization)
	request.Header.Set("Accept", p.Mimetype)
	request.Header.SetMethod("POST")
	request.Header.SetContentType(p.Mimetype)
	fmt.Println("", p.Mimetype)
	request.SetBody(image)
	request.SetRequestURI(p.URL)

	client := &fasthttp.Client{}
	if err = client.Do(request, response); err != nil {
		err = fmt.Errorf("client.Do(...): %w", err)
		return tevianApiResponse, err
	}

	fmt.Println(string(response.Body()))

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
