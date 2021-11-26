package sirius

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type FirmDetails struct {
	ID           int    `json:"id"`
	FirmName     string `json:"firmName"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	AddressLine3 string `json:"addressLine3"`
	Town         string `json:"town"`
	County       string `json:"county"`
	Postcode     string `json:"postcode"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phoneNumber"`
}

func (c *Client) AddFirmDetails(ctx Context, addFirmForm FirmDetails) (int, error) {
	var k FirmDetails
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(FirmDetails{
		FirmName:     addFirmForm.FirmName,
		AddressLine1: addFirmForm.AddressLine1,
		AddressLine2: addFirmForm.AddressLine2,
		AddressLine3: addFirmForm.AddressLine3,
		Town:         addFirmForm.Town,
		County:       addFirmForm.County,
		Postcode:     addFirmForm.Postcode,
		Email:        addFirmForm.Email,
		PhoneNumber:  addFirmForm.PhoneNumber,
	})
	if err != nil {
		return 0, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/v1/firm", &body)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return 0, ErrUnauthorized
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300

	if !statusOK {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return 0, ValidationError{
				Errors: v.ValidationErrors,
			}
		}

		return 0, newStatusError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&k)
	return k.ID, err
}
