package imgbb

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"social-backend/internal/domain/image"
	e "social-backend/internal/infrastructure/errors"

	"github.com/goccy/go-json"
	"golang.org/x/sync/errgroup"
)

type uploadResponseBody struct {
	Success bool `json:"success"`
	Data    struct {
		URL       string `json:"url"`
		DeleteUrl string `json:"delete_url"`
	} `json:"data"`
}

type ImgBBService struct {
	apiKey string
	apiURL string
}

func NewImgBBService(apiKey string, apiURL string) *ImgBBService {
	return &ImgBBService{apiKey, apiURL}
}

func (service *ImgBBService) Upload(targetImage io.Reader) (image.Image, error) {
	imgBytes, err := io.ReadAll(targetImage)
	if err != nil {
		return image.Image{}, err
	}

	encoded := base64.StdEncoding.EncodeToString(imgBytes)

	data := url.Values{}
	data.Set("key", service.apiKey)
	data.Set("image", encoded)

	resp, err := http.PostForm(service.apiURL, data)
	if err != nil {
		return image.Image{}, err
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	}()

	body, _ := io.ReadAll(resp.Body)

	var res uploadResponseBody

	if err = json.Unmarshal(body, &res); err != nil {
		return image.Image{}, err
	}

	if !res.Success {
		return image.Image{}, e.ImgBBUploadingError
	}

	resultImage := image.Image{
		Url:       res.Data.URL,
		DeleteUrl: res.Data.DeleteUrl,
	}

	return resultImage, nil
}

func (service *ImgBBService) UploadImages(ctx context.Context, images []io.Reader) ([]image.Image, error) {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(5)

	uploadedImages := make([]image.Image, len(images))

	for i, img := range images {
		i, img := i, img
		g.Go(func() error {
			uploaded, err := service.Upload(img)
			if err != nil {
				return err
			}

			uploadedImages[i] = uploaded
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return uploadedImages, nil
}

func (service *ImgBBService) DeleteImage(deleteUrl string) (bool, error) {
	resp, err := http.Get(deleteUrl)
	if err != nil {
		return false, err
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}
