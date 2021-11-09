package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type FirmDetails struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Client) GetFirmDetails(ctx Context, deputyId int) (FirmDetails, error) {
	var v FirmDetails

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/%d/firm", deputyId), nil)
	if err != nil {
		return v, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return v, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return v, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return v, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v, err
}
