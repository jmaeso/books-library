package classify

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Service interface {
	FindByID(id string) (ClassifyBookResponse, error)
	SearchByTitle(title string) ([]SearchResult, error)
}

type Client struct {
	baseURL         string
	baseQueryParams string
}

type SearchResult struct {
	Title  string `xml:"title,attr"`
	Author string `xml:"author,attr"`
	Year   string `xml:"hyr,attr"`
	ID     string `xml:"owi,attr"`
}

type classifyResponse struct {
	Results []SearchResult `xml:"works>work"`
}

type ClassifyBookResponse struct {
	BookData struct {
		Title  string `xml:"title,attr"`
		Author string `xml:"author,attr"`
		ID     string `xml:"owi,attr"`
	} `xml:"work"`
	Classification struct {
		MostPopular string `xml:"sfa,attr"`
	} `xml:"recommendations>ddc>mostPopular"`
}

func NewClient() *Client {
	return &Client{
		baseURL:         "http://classify.oclc.org/classify2/Classify",
		baseQueryParams: "?&summary=true",
	}
}

func (c *Client) FindByID(id string) (ClassifyBookResponse, error) {
	body, err := c.query("&owi=" + url.QueryEscape(id))
	if err != nil {
		return ClassifyBookResponse{}, err
	}

	var resp ClassifyBookResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return ClassifyBookResponse{}, err
	}

	return resp, nil
}

func (c *Client) SearchByTitle(title string) ([]SearchResult, error) {
	body, err := c.query("&title=" + url.QueryEscape(title))
	if err != nil {
		return []SearchResult{}, err
	}

	var resp classifyResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return []SearchResult{}, err
	}

	return resp.Results, nil
}

func (c *Client) query(query string) ([]byte, error) {
	resp, err := http.Get(c.baseURL + c.baseQueryParams + query)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
