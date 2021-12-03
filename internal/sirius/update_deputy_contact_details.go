package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DeputyContactDetails struct {
	DeputyFirstName string `json:"firstname"`
	DeputySurname   string `json:"surname"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phoneNumber"`
	AddressLine1    string `json:"addressLine1"`
	AddressLine2    string `json:"addressLine2"`
	AddressLine3    string `json:"addressLine3"`
	Town            string `json:"town"`
	County          string `json:"county"`
	Postcode        string `json:"postcode"`
}

func (c *Client) UpdateDeputyContactDetails(ctx Context, deputyId int, deputyDetails DeputyContactDetails) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(deputyDetails)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("/api/v1/deputies/%d", deputyId)

	req, err := c.newRequest(ctx, http.MethodPut, requestURL, &body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return newStatusError(resp)
	}

	return nil
}
