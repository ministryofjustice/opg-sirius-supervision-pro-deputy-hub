package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type apiTeam struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Members     []struct {
		ID          int    `json:"id"`
		DisplayName string `json:"displayName"`
		Email       string `json:"email"`
	} `json:"members"`
	TeamType *struct {
		Handle string `json:"handle"`
		Label  string `json:"label"`
	} `json:"teamType"`
}

type TeamMember struct {
	ID          int
	DisplayName string
}

type Team struct {
	Members []TeamMember
}

func (c *Client) GetProTeamsMembers(ctx Context, defaultProTeams string) ([]TeamMember, error) {
	requestURL := fmt.Sprintf("/api/v1/teams?type=%s", defaultProTeams)
	req, err := c.newRequest(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return []TeamMember{}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return []TeamMember{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return []TeamMember{}, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return []TeamMember{}, newStatusError(resp)
	}

	var v []apiTeam
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return []TeamMember{}, err
	}

	members := []TeamMember{}

	for _, k := range v {
		for _, m := range k.Members {
			members = append(members, TeamMember{
				ID:          m.ID,
				DisplayName: m.DisplayName,
			})
		}
	}

	return members, nil
}
