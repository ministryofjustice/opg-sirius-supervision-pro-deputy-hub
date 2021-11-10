package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type executiveCaseManager struct {
	EcmId   int    `json:"id"`
	EcmName string `json:"displayName"`
}

type ProDeputyDetails struct {
	ID                   int                  `json:"id"`
	DeputyFirstName      string               `json:"firstname"`
	DeputySurname        string               `json:"surname"`
	DeputyNumber         int                  `json:"deputyNumber"`
	OrganisationName     string               `json:"organisationName"`
	ExecutiveCaseManager executiveCaseManager `json:"executiveCaseManager"`
}

func (c *Client) GetProDeputyDetails(ctx Context, deputyId int) (ProDeputyDetails, error) {
	var v ProDeputyDetails

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deputies/%d", deputyId), nil)
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
