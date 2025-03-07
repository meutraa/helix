package helix

import "errors"

// ExpiresAt must be parsed manually since an empty string means perma ban
type Ban struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
	ExpiresAt Time   `json:"expires_at"`
}

type ManyBans struct {
	Bans       []Ban      `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type BannedUsersResponse struct {
	ResponseCommon
	Data ManyBans
}

// BroadcasterID must match the auth tokens user_id
type BannedUsersParams struct {
	BroadcasterID string `query:"broadcaster_id"`
	UserID        string `query:"user_id"`
	After         string `query:"after"`
	Before        string `query:"before"`
}

// GetBannedUsers returns all banned and timed-out users in a channel.
//
// Required scope: moderation:read
func (c *Client) GetBannedUsers(params *BannedUsersParams) (*BannedUsersResponse, error) {
	resp, err := c.get("/moderation/banned", &ManyBans{}, params)
	if err != nil {
		return nil, err
	}

	bans := &BannedUsersResponse{}
	resp.HydrateResponseCommon(&bans.ResponseCommon)
	bans.Data.Bans = resp.Data.(*ManyBans).Bans
	bans.Data.Pagination = resp.Data.(*ManyBans).Pagination

	return bans, nil
}

type BanUserParams struct {
	BroadcasterID string             `json:"broadcaster_id"`
	ModeratorId   string             `json:"moderator_id"`
	Body          BanUserRequestBody `json:"data"`
}

type BanUserRequestBody struct {
	Duration int    `json:"duration"` // optional
	Reason   string `json:"reason"`   // required
	UserId   string `json:"user_id"`  // required
}

type BanUserResponse struct {
	ResponseCommon
	Data ManyBanUser
}

type ManyBanUser struct {
	Bans []BanUser `json:"data"`
}

type BanUser struct {
	BoardcasterId string `json:"broadcaster_id"`
	CreatedAt     Time   `json:"created_at"`
	EndTime       Time   `json:"end_time"`
	ModeratorId   string `json:"moderator_id"`
	UserId        string `json:"user_id"`
}

// BanUser Bans a user from participating in a broadcaster’s chat room, or puts them in a timeout.
// Required scope: moderator:manage:banned_users
func (c *Client) BanUser(params *BanUserParams) (*BanUserResponse, error) {
	resp, err := c.postAsJSON("/moderation/bans", &ManyBanUser{}, params)
	if err != nil {
		return nil, err
	}

	banResp := &BanUserResponse{}
	resp.HydrateResponseCommon(&banResp.ResponseCommon)
	banResp.Data.Bans = resp.Data.(*ManyBanUser).Bans

	return banResp, nil
}

type UnbanUserParams struct {
	BroadcasterID string `json:"broadcaster_id"`
	ModeratorID   string `json:"moderator_id"`
	UserID        string `json:"user_id"`
}

type UnbanUserResponse struct {
	ResponseCommon
}

// UnbanUser Removes the ban or timeout that was placed on the specified user
// Required scope: moderator:manage:banned_users
func (c *Client) UnbanUser(params *UnbanUserParams) (*UnbanUserResponse, error) {
	resp, err := c.delete("/moderation/bans", nil, params)
	if err != nil {
		return nil, err
	}

	unbanResp := &UnbanUserResponse{}
	resp.HydrateResponseCommon(&unbanResp.ResponseCommon)
	return unbanResp, nil
}

type BlockedTermsParams struct {
	// Required
	BroadcasterID string `json:"broadcaster_id"`
	ModeratorID   string `json:"moderator_id"`

	// Optional
	After string `json:"after"`
	First int    `json:"first"`
}

type BlockedTermsResponse struct {
	ResponseCommon
	Data ManyBlockedTerms
}

type ManyBlockedTerms struct {
	Terms      []BlockedTerm `json:"data"`
	Pagination Pagination    `json:"pagination"`
}

type BlockedTerm struct {
	BroadcasterID string `json:"broadcaster_id"`
	CreatedAt     Time   `json:"created_at"`
	ExpiresAt     Time   `json:"expires_at"`
	ID            string `json:"id"`
	ModeratorID   string `json:"moderator_id"`
	Text          string `json:"text"`
	UpdatedAt     Time   `json:"updated_at"`
}

// GetBlockedTerms Gets the broadcaster’s list of non-private, blocked words or phrases.
// These are the terms that the broadcaster or moderator added manually, or that were denied by AutoMod.
// Required scope: moderator:read:blocked_terms
func (c *Client) GetBlockedTerms(params *BlockedTermsParams) (*BlockedTermsResponse, error) {
	if params.BroadcasterID == "" || params.ModeratorID == "" {
		return nil, errors.New("broadcaster id and moderator id must be provided")
	}

	resp, err := c.get("/moderation/blocked_terms", &ManyBlockedTerms{}, params)
	if err != nil {
		return nil, err
	}

	blockedTermsResp := &BlockedTermsResponse{}
	resp.HydrateResponseCommon(&blockedTermsResp.ResponseCommon)
	blockedTermsResp.Data.Terms = resp.Data.(*ManyBlockedTerms).Terms
	blockedTermsResp.Data.Pagination = resp.Data.(*ManyBlockedTerms).Pagination

	return blockedTermsResp, nil
}

type AddBlockedTermParams struct {
	BroadcasterID string `json:"broadcaster_id"`
	ModeratorID   string `json:"moderator_id"`
	Text          string `json:"text"`
}

type AddBlockedTermResponse struct {
	ResponseCommon
	Data ManyAddBlockedTerms
}

type ManyAddBlockedTerms struct {
	Terms []BlockedTerm `json:"data"`
}

// AddBlockedTerm Adds a word or phrase to the broadcaster’s list of blocked terms.
// These are the terms that broadcasters don’t want used in their chat room.
// Required scope: moderator:manage:blocked_terms
func (c *Client) AddBlockedTerm(params *AddBlockedTermParams) (*AddBlockedTermResponse, error) {
	if params.BroadcasterID == "" || params.ModeratorID == "" {
		return nil, errors.New("broadcaster id and moderator id must be provided")
	}
	if len(params.Text) < 2 || len(params.Text) > 500 {
		return nil, errors.New("the term len must be between 2 and 500")
	}

	resp, err := c.post("/moderation/blocked_terms", &ManyAddBlockedTerms{}, params)
	if err != nil {
		return nil, err
	}

	addTermResp := &AddBlockedTermResponse{}
	resp.HydrateResponseCommon(&addTermResp.ResponseCommon)
	addTermResp.Data.Terms = resp.Data.(*ManyAddBlockedTerms).Terms

	return addTermResp, nil
}

type RemoveBlockedTermParams struct {
	BroadcasterID string `json:"broadcaster_id"`
	ModeratorID   string `json:"moderator_id"`
	ID            string `json:"id"`
}

type RemoveBlockedTermResponse struct {
	ResponseCommon
}

// RemoveBlockedTerm Removes the word or phrase that the broadcaster is blocking users from using in their chat room.
// Required scope: moderator:manage:blocked_terms
func (c *Client) RemoveBlockedTerm(params *RemoveBlockedTermParams) (*RemoveBlockedTermResponse, error) {
	if params.BroadcasterID == "" || params.ModeratorID == "" {
		return nil, errors.New("broadcaster id and moderator id must be provided")
	}

	if params.ID == "" {
		return nil, errors.New("id must be provided")
	}

	resp, err := c.delete("/moderation/blocked_terms", nil, params)
	if err != nil {
		return nil, err
	}

	blockedTermResp := &RemoveBlockedTermResponse{}
	resp.HydrateResponseCommon(&blockedTermResp.ResponseCommon)

	return blockedTermResp, nil
}
