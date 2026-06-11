package family

type Member struct {
	Id              int64  `json:"id"`
	FamilyId        int64  `json:"familyId"`
	UserId          int64  `json:"userId"`
	Account         string `json:"account"`
	Name            string `json:"name"`
	Role            string `json:"role"`
	Status          string `json:"status"`
	DisplayName     string `json:"displayName"`
	InvitedByUserId int64  `json:"invitedByUserId"`
	FamilyName      string `json:"familyName"`
}

type InviteRequest struct {
	Account     string `json:"account"`
	DisplayName string `json:"displayName"`
}
