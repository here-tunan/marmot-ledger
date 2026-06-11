package familydb

import (
	"errors"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"

	"xorm.io/xorm"
)

type Member struct {
	Id              int64            `json:"id" xorm:"pk autoincr 'id'"`
	FamilyId        int64            `json:"familyId" xorm:"'family_id'"`
	UserId          int64            `json:"userId" xorm:"'user_id'"`
	Role            string           `json:"role" xorm:"'role'"`
	Status          string           `json:"status" xorm:"'status'"`
	DisplayName     string           `json:"displayName" xorm:"'display_name'"`
	InvitedByUserId *int64           `json:"invitedByUserId" xorm:"'invited_by_user_id'"`
	InvitedAt       *model.LocalTime `json:"invitedAt" xorm:"'invited_at'"`
	JoinedAt        *model.LocalTime `json:"joinedAt" xorm:"'joined_at'"`
	LeftAt          *model.LocalTime `json:"leftAt" xorm:"'left_at'"`
	CreatedAt       model.LocalTime  `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt       model.LocalTime  `json:"updatedAt" xorm:"updated 'updated_at'"`
}

type MemberView struct {
	Member     `xorm:"extends"`
	Account    string `json:"account" xorm:"'account'"`
	Name       string `json:"name" xorm:"'name'"`
	FamilyName string `json:"familyName" xorm:"'family_name'"`
}

func (Member) TableName() string { return "family_member" }

func InsertMember(session *xorm.Session, member *Member) error {
	_, err := session.InsertOne(member)
	return err
}

func GetMember(familyId int64, userId int64) (*Member, error) {
	m := &Member{}
	has, err := infrastructure.Mysql.Where("family_id = ? AND user_id = ?", familyId, userId).Get(m)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("family member not found")
	}
	return m, nil
}

func ListMembers(familyId int64, includeInvited bool) ([]MemberView, error) {
	items := make([]MemberView, 0)
	session := infrastructure.Mysql.Table("family_member").Alias("fm").
		Select("fm.*, u.account AS account, u.name AS name").
		Join("INNER", "user u", "fm.user_id = u.id").
		Where("fm.family_id = ?", familyId)
	if includeInvited {
		session.And("fm.status IN ('active', 'invited')")
	} else {
		session.And("fm.status = 'active'")
	}
	err := session.Asc("fm.id").Find(&items)
	return items, err
}

func ListInvitations(userId int64) ([]MemberView, error) {
	items := make([]MemberView, 0)
	err := infrastructure.Mysql.Table("family_member").Alias("fm").
		Select("fm.*, f.name AS family_name, u.account AS account, u.name AS name").
		Join("INNER", "family f", "fm.family_id = f.id").
		Join("INNER", "user u", "fm.user_id = u.id").
		Where("fm.user_id = ? AND fm.status = 'invited' AND f.is_deleted = 0", userId).
		Asc("fm.id").Find(&items)
	return items, err
}

func UpdateMemberStatus(id int64, userId int64, status string, joinedAt *model.LocalTime, leftAt *model.LocalTime) error {
	_, err := infrastructure.Mysql.Where("id = ? AND user_id = ?", id, userId).Cols("status", "joined_at", "left_at").Update(&Member{Status: status, JoinedAt: joinedAt, LeftAt: leftAt})
	return err
}
