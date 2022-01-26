package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AdditionalInformationDetails struct {
	Complaints string `json:"complaints"`
	PanelDeputy  string `json:"panelDeputy"`
	AnnualBillingPreference    string `json:"annualBillingPreference"`
}

func (c *Client) UpdateAdditionalInformation(ctx Context, deputyId int, additionalInfoForm AdditionalInformationDetails) error {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(additionalInfoForm)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("/api/v1/deputies/%d/additional-information", deputyId)

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

	if resp.StatusCode == http.StatusBadRequest {
		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			return ValidationError{
				Errors: v.ValidationErrors,
			}
		}
	}

	if resp.StatusCode != http.StatusOK {
		return newStatusError(resp)
	}

	return nil
}
