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

type firm struct {
	FirmName string `json:"firmName"`
	FirmId   int    `json:"id"`
}

type deputySubType struct {
	SubType string `json:"handle"`
}

//"deputyImportantInformation":{"id":18,"complaints":false,"panelDeputy":false,"annualBillingInvoice":{"handle":"SCHEDULE AND INVOICE","label":"Schedule and Invoice"},"otherImportantInformation":""},"deputySubType":{"handle":"ORGANISATION","label":"Organisation"}}

type ProDeputyDetails struct {
	ID                               int                  `json:"id"`
	DeputyFirstName                  string               `json:"firstname"`
	DeputySurname                    string               `json:"surname"`
	DeputyNumber                     int                  `json:"deputyNumber"`
	DeputySubType                    deputySubType        `json:"deputySubType"`
	DeputyImportantInformation struct {
		Id                   int  `json:"id"`
		Complaints           bool `json:"complaints"`
		PanelDeputy          bool `json:"panelDeputy"`
		AnnualBillingInvoice struct {
			Handle string `json:"handle"`
			Label  string `json:"label"`
		} `json:"annualBillingInvoice"`
		OtherImportantInformation string `json:"otherImportantInformation"`
	} `json:"deputyImportantInformation"`
	OrganisationName                 string               `json:"organisationName"`
	OrganisationTeamOrDepartmentName string               `json:"organisationTeamOrDepartmentName"`
	ExecutiveCaseManager             executiveCaseManager `json:"executiveCaseManager"`
	Firm                             firm                 `json:"firm"`
	Email                            string               `json:"email"`
	PhoneNumber                      string               `json:"phoneNumber"`
	AddressLine1                     string               `json:"addressLine1"`
	AddressLine2                     string               `json:"addressLine2"`
	AddressLine3                     string               `json:"addressLine3"`
	Town                             string               `json:"town"`
	County                           string               `json:"county"`
	Postcode                         string               `json:"postcode"`
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

	fmt.Println("kate")
	fmt.Println(v.DeputyImportantInformation)

	return v, err
}
