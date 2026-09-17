package xbox

type authTokens struct {
	UserHash  string `json:"userHash"`
	XSTSToken string `json:"XSTSToken"`
}

type profileResponse struct {
	Users []profileUser `json:"profileUsers"`
}

type profileUser struct {
	ID string `json:"id"`
}

type friendList struct {
	People []Person `json:"people"`
}

type Person struct {
	XUID          string `json:"xuid"`
	Gamertag      string `json:"gamertag"`
	DisplayName   string `json:"displayName"`
	GamerScore    string `json:"gamerScore"`
	DisplayPicRaw string `json:"displayPicRaw"`
}
