package clients

import (
	"encoding/json"
	core "fl/my-portfolio/internal/app"
	"fl/my-portfolio/internal/configs"
	"io"
	"net/http"
)

type CMCClient struct {
	config configs.CMCConfig
}

func NewCMCClient() *CMCClient {
	return &CMCClient{
		config: core.GetConfig().CMC,
	}
}

func (c *CMCClient) GetCryptoListingsLatest() (*CryptoListingsLatestResponse, error) {
	var obj *CryptoListingsLatestResponse

	response, err := c.sendRequest(cryptoListingLatest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (c *CMCClient) sendRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-CMC_PRO_API_KEY", c.config.APIKey)

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
