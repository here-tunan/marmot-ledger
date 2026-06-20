package service

import (
	"errors"
	"marmot-ledger/internal/domain/entity/family"
	"marmot-ledger/internal/domain/entity/financialevent"
	"marmot-ledger/internal/domain/entity/statistics"
	"marmot-ledger/internal/domain/repository/accountdb"
	"marmot-ledger/internal/domain/repository/bucketdb"
	"marmot-ledger/internal/domain/repository/familydb"
	"marmot-ledger/internal/domain/repository/financialeventdb"
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
	fam := &familydb.Family{Name: strings.TrimSpace(req.Name), OwnerUserId: userId}
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
	return &family.Family{Id: fam.Id, Name: fam.Name, OwnerUserId: fam.OwnerUserId, Role: FamilyRoleOwner}, nil
}

func ListFamilies(userId int64) ([]family.Family, error) {
	items, err := familydb.ListFamilies(userId)
	if err != nil {
		return nil, err
	}
	result := make([]family.Family, 0, len(items))
	for _, item := range items {
		result = append(result, family.Family{Id: item.Id, Name: item.Name, OwnerUserId: item.OwnerUserId, Role: item.Role})
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
	return &family.Family{Id: fam.Id, Name: fam.Name, OwnerUserId: fam.OwnerUserId, Role: member.Role}, nil
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

func GetFamilyStatisticsSummaries(userId int64, familyId int64, query statisticsdb.StatisticsQuery) ([]statistics.Summary, error) {
	userIds, err := activeFamilyUserIds(userId, familyId)
	if err != nil {
		return nil, err
	}
	return statisticsdb.GetSummariesByUserIds(userIds, query)
}

func GetFamilyStatisticsCategoryGroups(userId int64, familyId int64, query statisticsdb.StatisticsQuery) ([]statistics.CategoryGroupStats, error) {
	userIds, err := activeFamilyUserIds(userId, familyId)
	if err != nil {
		return nil, err
	}
	return statisticsdb.GetCategoryGroupStatsListByUserIds(userIds, query)
}

func GetFamilyNetWorthTrend(userId int64, familyId int64, query statisticsdb.StatisticsQuery, granularity string) ([]statistics.NetWorthTrendPoint, error) {
	userIds, err := activeFamilyUserIds(userId, familyId)
	if err != nil {
		return nil, err
	}
	return statisticsdb.GetNetWorthTrendByUserIds(userIds, query, granularity)
}

func ListFamilyFinancialEvents(userId int64, familyId int64, query financialevent.FinancialEventQuery) (*PageResult[financialevent.FinancialEvent], error) {
	userIds, err := activeFamilyUserIds(userId, familyId)
	if err != nil {
		return nil, err
	}
	page := query.Page
	pageSize := query.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	query.Page = page
	query.PageSize = pageSize
	events, total, err := financialeventdb.ListFinancialEventsByUserIds(userIds, financialeventdb.FinancialEventQuery{
		EventType:           query.EventType,
		StartTime:           query.StartTime,
		EndTime:             query.EndTime,
		Currency:            query.Currency,
		CategoryId:          query.CategoryId,
		CategoryGroupId:     query.CategoryGroupId,
		BucketId:            query.BucketId,
		Keyword:             query.Keyword,
		IncludeInStatistics: query.IncludeInStatistics,
		Page:                query.Page,
		PageSize:            query.PageSize,
	})
	if err != nil {
		return nil, err
	}
	result := make([]financialevent.FinancialEvent, 0, len(events))
	for _, item := range events {
		result = append(result, *toFinancialEventEntity(&item, nil))
	}
	return &PageResult[financialevent.FinancialEvent]{List: result, Page: page, PageSize: pageSize, Total: total}, nil
}

func GetFamilyAssets(userId int64, familyId int64) (*statistics.FamilyAssets, error) {
	if err := ensureFamilyMember(userId, familyId); err != nil {
		return nil, err
	}
	members, err := familydb.ListMembers(familyId, false)
	if err != nil {
		return nil, err
	}
	result := &statistics.FamilyAssets{Members: make([]statistics.MemberAsset, 0, len(members))}
	for _, member := range members {
		accounts, err := accountdb.ListAccounts(member.UserId, accountdb.AccountQuery{})
		if err != nil {
			return nil, err
		}
		asset := statistics.MemberAsset{UserId: member.UserId, Account: member.Account, Name: member.Name, DisplayName: member.DisplayName, Role: member.Role, Totals: make([]statistics.MemberAssetTotal, 0), Accounts: make([]statistics.MemberAssetAccount, 0, len(accounts))}
		totals := make(map[string]*statistics.MemberAssetTotal)
		for _, account := range accounts {
			active := true
			buckets, err := bucketdb.ListBuckets(member.UserId, bucketdb.BucketQuery{AccountId: account.Id, IsActive: &active})
			if err != nil {
				return nil, err
			}
			accountAsset := statistics.MemberAssetAccount{Id: account.Id, Name: account.Name, Type: account.Type, Buckets: make([]statistics.MemberAssetBucket, 0, len(buckets))}
			for _, bucket := range buckets {
				accountAsset.Buckets = append(accountAsset.Buckets, statistics.MemberAssetBucket{Id: bucket.Id, Name: bucket.Name, Currency: bucket.Currency, Balance: bucket.Balance, BucketType: bucket.BucketType, BucketNature: bucket.BucketNature})
				total := totals[bucket.Currency]
				if total == nil {
					total = &statistics.MemberAssetTotal{Currency: bucket.Currency}
					totals[bucket.Currency] = total
				}
				if bucket.BucketNature == BucketNatureLiability {
					total.Liability = total.Liability.Add(bucket.Balance)
				} else {
					total.Asset = total.Asset.Add(bucket.Balance)
				}
				total.BucketCount++
			}
			if len(accountAsset.Buckets) > 0 {
				asset.Accounts = append(asset.Accounts, accountAsset)
			}
		}
		for _, total := range totals {
			total.NetWorth = total.Asset.Sub(total.Liability)
			asset.Totals = append(asset.Totals, *total)
		}
		if asset.Totals == nil {
			asset.Totals = []statistics.MemberAssetTotal{}
		}
		if asset.Accounts == nil {
			asset.Accounts = []statistics.MemberAssetAccount{}
		}
		result.Members = append(result.Members, asset)
	}
	return result, nil
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
func ActiveFamilyUserIds(userId int64, familyId int64) ([]int64, error) {
	return activeFamilyUserIds(userId, familyId)
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
