package service

import (
	"errors"
	"marmot-ledger/internal/domain/entity/family"
	"marmot-ledger/internal/domain/entity/statistics"
	"marmot-ledger/internal/domain/repository/familydb"
	"marmot-ledger/internal/domain/repository/statisticsdb"
	"marmot-ledger/internal/domain/repository/userdb"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"
	"strings"
	"time"
)

const (
	FamilyRoleOwner  = "owner"
	FamilyRoleAdmin  = "admin"
	FamilyRoleMember = "member"

	FamilyMemberStatusInvited  = "invited"
	FamilyMemberStatusActive   = "active"
	FamilyMemberStatusRejected = "rejected"
	FamilyMemberStatusLeft     = "left"
)

func CreateFamily(userId int64, req *family.CreateFamilyRequest) (*family.Family, error) {
	if req == nil || strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("family name is required")
	}
	baseCurrency := strings.ToUpper(strings.TrimSpace(req.BaseCurrency))
	if baseCurrency == "" {
		baseCurrency = "CNY"
	}
	session := infrastructure.Mysql.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return nil, err
	}
	committed := false
	defer func() {
		if !committed {
			_ = session.Rollback()
		}
	}()
	fam := &familydb.Family{Name: strings.TrimSpace(req.Name), BaseCurrency: baseCurrency, OwnerUserId: userId}
	if err := familydb.InsertFamily(session, fam); err != nil {
		return nil, err
	}
	now := model.LocalTime(time.Now())
	member := &familydb.Member{FamilyId: fam.Id, UserId: userId, Role: FamilyRoleOwner, Status: FamilyMemberStatusActive, JoinedAt: &now}
	if err := familydb.InsertMember(session, member); err != nil {
		return nil, err
	}
	if err := session.Commit(); err != nil {
		return nil, err
	}
	committed = true
	return &family.Family{Id: fam.Id, Name: fam.Name, BaseCurrency: fam.BaseCurrency, OwnerUserId: fam.OwnerUserId, Role: FamilyRoleOwner}, nil
}

func ListFamilies(userId int64) ([]family.Family, error) {
	items, err := familydb.ListFamilies(userId)
	if err != nil {
		return nil, err
	}
	result := make([]family.Family, 0, len(items))
	for _, item := range items {
		result = append(result, family.Family{Id: item.Id, Name: item.Name, BaseCurrency: item.BaseCurrency, OwnerUserId: item.OwnerUserId, Role: item.Role})
	}
	return result, nil
}

func GetFamily(userId int64, familyId int64) (*family.Family, error) {
	ok, err := familydb.IsActiveMember(familyId, userId)
	if err != nil || !ok {
		return nil, errors.New("family not found")
	}
	fam, err := familydb.GetFamily(familyId)
	if err != nil {
		return nil, err
	}
	member, _ := familydb.GetMember(familyId, userId)
	return &family.Family{Id: fam.Id, Name: fam.Name, BaseCurrency: fam.BaseCurrency, OwnerUserId: fam.OwnerUserId, Role: member.Role}, nil
}

func ListFamilyMembers(userId int64, familyId int64, includeInvited bool) ([]family.Member, error) {
	if err := ensureFamilyMember(userId, familyId); err != nil {
		return nil, err
	}
	items, err := familydb.ListMembers(familyId, includeInvited)
	if err != nil {
		return nil, err
	}
	result := make([]family.Member, 0, len(items))
	for _, item := range items {
		result = append(result, toFamilyMemberEntity(&item))
	}
	return result, nil
}

func InviteFamilyMember(userId int64, familyId int64, req *family.InviteRequest) (*family.Member, error) {
	if err := ensureFamilyAdmin(userId, familyId); err != nil {
		return nil, err
	}
	if req == nil || strings.TrimSpace(req.Account) == "" {
		return nil, errors.New("account is required")
	}
	user := &userdb.User{Account: strings.TrimSpace(req.Account)}
	if err := user.GetUser(); err != nil {
		return nil, err
	}
	if user.Id == userId {
		return nil, errors.New("cannot invite yourself")
	}
	if old, err := familydb.GetMember(familyId, user.Id); err == nil && old.Status != FamilyMemberStatusLeft && old.Status != FamilyMemberStatusRejected {
		return nil, errors.New("member already exists")
	}
	now := model.LocalTime(time.Now())
	inviter := userId
	member := &familydb.Member{FamilyId: familyId, UserId: user.Id, Role: FamilyRoleMember, Status: FamilyMemberStatusInvited, DisplayName: strings.TrimSpace(req.DisplayName), InvitedByUserId: &inviter, InvitedAt: &now}
	session := infrastructure.Mysql.NewSession()
	defer session.Close()
	if err := familydb.InsertMember(session, member); err != nil {
		return nil, err
	}
	view := &familydb.MemberView{Member: *member, Account: user.Account, Name: user.Name}
	entity := toFamilyMemberEntity(view)
	return &entity, nil
}

func ListFamilyInvitations(userId int64) ([]family.Member, error) {
	items, err := familydb.ListInvitations(userId)
	if err != nil {
		return nil, err
	}
	result := make([]family.Member, 0, len(items))
	for _, item := range items {
		result = append(result, toFamilyMemberEntity(&item))
	}
	return result, nil
}

func AcceptFamilyInvitation(userId int64, invitationId int64) error {
	now := model.LocalTime(time.Now())
	return familydb.UpdateMemberStatus(invitationId, userId, FamilyMemberStatusActive, &now, nil)
}
func RejectFamilyInvitation(userId int64, invitationId int64) error {
	return familydb.UpdateMemberStatus(invitationId, userId, FamilyMemberStatusRejected, nil, nil)
}

func GetFamilyStatisticsSummary(userId int64, familyId int64, query statisticsdb.StatisticsQuery) (*statistics.Summary, error) {
	userIds, err := activeFamilyUserIds(userId, familyId)
	if err != nil {
		return nil, err
	}
	return statisticsdb.GetSummaryByUserIds(userIds, query)
}
func GetFamilyStatisticsCategoryGroup(userId int64, familyId int64, query statisticsdb.StatisticsQuery) (*statistics.CategoryGroupStats, error) {
	userIds, err := activeFamilyUserIds(userId, familyId)
	if err != nil {
		return nil, err
	}
	return statisticsdb.GetCategoryGroupStatsByUserIds(userIds, query)
}

func ensureFamilyMember(userId int64, familyId int64) error {
	ok, err := familydb.IsActiveMember(familyId, userId)
	if err != nil || !ok {
		return errors.New("not family member")
	}
	return nil
}
func ensureFamilyAdmin(userId int64, familyId int64) error {
	member, err := familydb.GetMember(familyId, userId)
	if err != nil || member.Status != FamilyMemberStatusActive {
		return errors.New("not family member")
	}
	if member.Role != FamilyRoleOwner && member.Role != FamilyRoleAdmin {
		return errors.New("no family permission")
	}
	return nil
}
func activeFamilyUserIds(userId int64, familyId int64) ([]int64, error) {
	if err := ensureFamilyMember(userId, familyId); err != nil {
		return nil, err
	}
	members, err := familydb.ListMembers(familyId, false)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(members))
	for _, m := range members {
		ids = append(ids, m.UserId)
	}
	return ids, nil
}
func toFamilyMemberEntity(item *familydb.MemberView) family.Member {
	invitedBy := int64(0)
	if item.InvitedByUserId != nil {
		invitedBy = *item.InvitedByUserId
	}
	return family.Member{Id: item.Id, FamilyId: item.FamilyId, UserId: item.UserId, Account: item.Account, Name: item.Name, Role: item.Role, Status: item.Status, DisplayName: item.DisplayName, InvitedByUserId: invitedBy, FamilyName: item.FamilyName}
}
