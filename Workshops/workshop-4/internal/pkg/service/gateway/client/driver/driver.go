package driver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"workshop.4.2/internal/model"
)

const host = "http://localhost:9002"

func New() *Client {
	return &Client{
		cli: &http.Client{},
	}
}

type Client struct {
	cli *http.Client
}

// 1. Валидация входных данных
// 2. Подготовить запрос
// 3. Валидация результата
// 4. Декодировать результат

func (c *Client) Get(ctx context.Context, driverID model.DriverID) (model.Driver, error) {
	if driverID == 0 {
		return model.Driver{}, errors.New("invalid")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, host+"/order", nil)
	if err != nil {
		return model.Driver{}, err
	}

	q := req.URL.Query()
	q.Add(model.QueryParamUserID, strconv.FormatUint(uint64(driverID), 10))
	req.URL.RawQuery = q.Encode() // ?userId=1

	resp, err := c.cli.Do(req)
	if err != nil {
		return model.Driver{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Driver{}, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Driver{}, err
	}

	var driver model.Driver
	if err = json.Unmarshal(body, &driver); err != nil {
		return model.Driver{}, err
	}

	return driver, nil
}

func (c *Client) Free(ctx context.Context) (model.Driver, error) {
	req, err := http.NewRequest(http.MethodGet, host+"/free", nil)
	if err != nil {
		return model.Driver{}, err
	}

	req.WithContext(ctx)

	resp, err := c.cli.Do(req)
	if err != nil {
		return model.Driver{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Driver{}, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Driver{}, err
	}

	var driver model.Driver
	if err = json.Unmarshal(body, &driver); err != nil {
		return model.Driver{}, err
	}

	return driver, nil
}
