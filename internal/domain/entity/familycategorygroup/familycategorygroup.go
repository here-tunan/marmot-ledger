package familycategorygroup

type FamilyCategoryGroup struct {
	Id        int64  `json:"id"`
	FamilyId  int64  `json:"familyId"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Icon      string `json:"icon"`
	Color     string `json:"color"`
	CreatedBy int64  `json:"createdBy"`
	Sort      int    `json:"sort"`
	IsActive  bool   `json:"isActive"`
}

type CategoryGroupMember struct {
	Id            int64 `json:"id"`
	FamilyGroupId int64 `json:"familyGroupId"`
	CategoryId    int64 `json:"categoryId"`
	AddedBy       int64 `json:"addedBy"`
}

type CategoryGroupWithMembers struct {
	FamilyCategoryGroup
	MemberCount int     `json:"memberCount"`
	MemberIds   []int64 `json:"memberIds"`
}

type CreateGroupRequest struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
	Sort  int    `json:"sort"`
}

type UpdateGroupRequest struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
	Sort  int    `json:"sort"`
}

type AddMembersRequest struct {
	CategoryIds []int64 `json:"categoryIds"`
}

type CategoryWithGroups struct {
	CategoryId   int64   `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	Type         string  `json:"type"`
	GroupIds     []int64 `json:"groupIds"`
}
