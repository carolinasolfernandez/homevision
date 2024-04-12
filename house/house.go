package house

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
)

type houses struct {
	Houses []House `json:"houses"`
}

type House struct {
	Id        int    `json:"id"`
	Address   string `json:"address"`
	Homeowner string `json:"homeowner"`
	Price     int    `json:"price"`
	PhotoURL  string `json:"photoURL"`
}

type houseService struct {
	httpClient *http.Client
	url        string
	photosDir  string
}

func NewHouseService(httpClient *http.Client, url string, photosDir string) *houseService {
	return &houseService{
		httpClient: httpClient,
		url:        url,
		photosDir:  photosDir,
	}
}

// GetHouses calls the house service to get the houses
func (h houseService) GetHouses(pages, numPerPage int, houseCh chan<- []House, errorCh chan<- error) {
	var wg = &sync.WaitGroup{}

	for i := 1; i <= pages; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := h.getPage(i, numPerPage, houseCh)
			if err != nil {
				errorCh <- err
			}
		}(i)
	}
	wg.Wait()
	close(houseCh)
}

// SavePhotos gets each house photo and saves the image to a file
func (h houseService) SavePhotos(houseCh <-chan []House, doneCh chan<- struct{}, errorCh chan<- error) {
	var wg = &sync.WaitGroup{}
	for hc := range houseCh {
		for _, ho := range hc {
			wg.Add(1)
			go func(house House) {
				defer wg.Done()
				fileName := fmt.Sprintf("%d-%s.%s", house.Id, house.Address, path.Ext(house.PhotoURL))
				if err := h.savePhoto(fileName, house.PhotoURL); err != nil {
					errorCh <- err
				}
			}(ho)
		}
	}
	wg.Wait()
	doneCh <- struct{}{}
}

// getPage retrieves a single page of house data
func (h houseService) getPage(page, perPage int, housesCh chan<- []House) error {
	url := fmt.Sprintf("%s?page=%d&per_page=%d", h.url, page, perPage)
	log.Println("get houses: ", url)

	resp, err := h.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("getPage: error getting url : %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("getPage: failed to close response body: %v\n", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("getPage: error reading body : %w", err)
	}

	var houses houses
	if err = json.Unmarshal(body, &houses); err != nil {
		return fmt.Errorf("getPage: error unmarshaling response body: %w", err)
	}

	if houses.Houses == nil {
		return fmt.Errorf("getPage: no houses found")
	}
	housesCh <- houses.Houses
	return nil
}

// savePhoto saves the photo to a file
func (h houseService) savePhoto(fileName string, photoURL string) (err error) {
	resp, err := h.httpClient.Get(photoURL)
	if err != nil {
		return fmt.Errorf("savePhoto: HTTP error: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("savePhoto: failed to close response body: %v\n", err)
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("savePhoto: error reading response body: %w", err)
	}

	file, err := os.Create(h.photosDir + "/" + fileName)
	if err != nil {
		return fmt.Errorf("savePhoto: error creating file: %w", err)
	}

	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			log.Printf("savePhoto: failed to close file: %v\n", err)
		}
	}(file)

	if _, err = file.Write(respBody); err != nil {
		return fmt.Errorf("savePhoto: error writing photo to file: %w", err)
	}

	return err
}
