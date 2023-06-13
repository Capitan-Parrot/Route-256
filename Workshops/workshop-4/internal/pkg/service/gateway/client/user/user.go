package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"workshop.4.2/internal/model"
)

const host = "http://localhost:9000"

func New() *Client {
	return &Client{
		cli: &http.Client{},
	}
}

type Client struct {
	cli *http.Client
}

func (c *Client) Get(ctx context.Context, userID model.ClientID) (model.Client, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, host, nil)
	if err != nil {
		return model.Client{}, err
	}

	q := req.URL.Query()
	q.Add(model.QueryParamUserID, strconv.FormatUint(uint64(userID), 10))
	req.URL.RawQuery = q.Encode()

	resp, err := c.cli.Do(req)
	if err != nil {
		return model.Client{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Client{}, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Client{}, err
	}

	var user model.Client
	if err = json.Unmarshal(body, &user); err != nil {
		return model.Client{}, err
	}

	return user, nil
}
