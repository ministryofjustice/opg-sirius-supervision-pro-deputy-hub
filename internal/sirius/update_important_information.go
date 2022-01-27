package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ImportantInformationDetails struct {
	Complaints bool `json:"complaints"`
	PanelDeputy  bool `json:"panelDeputy"`
	AnnualBillingInvoice  string `json:"annualBillingInvoice"`
	OtherImportantInformation string `json:"otherImportantInformation"`
}

func (c *Client) UpdateImportantInformation(ctx Context, deputyId int, importantInfoForm ImportantInformationDetails) error {
	var body bytes.Buffer

	fmt.Println("sirius important info")
	fmt.Println(importantInfoForm)

	err := json.NewEncoder(&body).Encode(importantInfoForm)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("/api/v1/deputies/%d/important-information", deputyId)

	req, err := c.newRequest(ctx, http.MethodPut, requestURL, &body)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)

	if err != nil {
		return err
	}

	io.Copy(os.Stdout, resp.Body)

	defer resp.Body.Close()



	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if resp.StatusCode == http.StatusBadRequest {

		var v struct {
			ValidationErrors ValidationErrors `json:"validation_errors"`
		}
		fmt.Println("errors in internal sirius")
		fmt.Println(v.ValidationErrors)

		if err := json.NewDecoder(resp.Body).Decode(&v); err == nil {
			fmt.Println("errors in internal sirius")
			fmt.Println(v.ValidationErrors)
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
