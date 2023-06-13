package order

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"workshop.4.2/internal/model"
)

const host = "http://localhost:9001"

func New() *Client {
	return &Client{
		cli: &http.Client{},
	}
}

type Client struct {
	cli *http.Client
}

func (c *Client) Get(ctx context.Context, orderID model.OrderID) (order model.Order, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("client: [order], method: [Get], err: %v", err)
		}
	}()

	req, err := http.NewRequest(http.MethodGet, host, nil)
	if err != nil {
		return model.Order{}, err
	}

	req.WithContext(ctx)

	q := req.URL.Query()
	q.Add(model.QueryParamOrderID, strconv.FormatUint(uint64(orderID), 10))
	req.URL.RawQuery = q.Encode()

	resp, err := c.cli.Do(req)
	if err != nil {
		return model.Order{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Order{}, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Order{}, err
	}

	if err = json.Unmarshal(body, &order); err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (c *Client) Create(ctx context.Context, order model.Order) error {
	body, err := json.Marshal(&order)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, host, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.WithContext(ctx)

	resp, err := c.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) Update(ctx context.Context, order model.Order) error {
	body, err := json.Marshal(&order)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, host, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.WithContext(ctx)

	resp, err := c.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	return nil
}
