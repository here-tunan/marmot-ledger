package bucketdb

import (
	"errors"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"
	"strings"

	"github.com/shopspring/decimal"
	"xorm.io/xorm"
)

type Bucket struct {
	Id             int64           `json:"id" xorm:"pk autoincr 'id'"`
	UserId         int64           `json:"userId" xorm:"'user_id'"`
	AccountId      int64           `json:"accountId" xorm:"'account_id'"`
	Name           string          `json:"name" xorm:"'name'"`
	Currency       string          `json:"currency" xorm:"'currency'"`
	Balance        decimal.Decimal `json:"balance" xorm:"'balance'"`
	InitialBalance decimal.Decimal `json:"initialBalance" xorm:"'initial_balance'"`
	BucketType     string          `json:"bucketType" xorm:"'bucket_type'"`
	BucketNature   string          `json:"bucketNature" xorm:"'bucket_nature'"`
	BucketGroupKey string          `json:"bucketGroupKey" xorm:"'bucket_group_key'"`
	IsActive       bool            `json:"isActive" xorm:"'is_active'"`
	IsDeleted      bool            `json:"isDeleted" xorm:"'is_deleted'"`
	CreatedAt      model.LocalTime `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt      model.LocalTime `json:"updatedAt" xorm:"updated 'updated_at'"`
}

type BucketQuery struct {
	AccountId    int64
	Currency     string
	BucketType   string
	BucketNature string
	IsActive     *bool
}

func (Bucket) TableName() string {
	return "bucket"
}

func InsertBucket(session *xorm.Session, bucket *Bucket) error {
	_, err := session.InsertOne(bucket)
	return err
}

func ListBuckets(userId int64, query BucketQuery) ([]Bucket, error) {
	buckets := make([]Bucket, 0)
	session := infrastructure.Mysql.Where("user_id = ? AND is_deleted = ?", userId, 0)

	if query.AccountId != 0 {
		session.And("account_id = ?", query.AccountId)
	}
	if strings.TrimSpace(query.Currency) != "" {
		session.And("currency = ?", strings.ToUpper(strings.TrimSpace(query.Currency)))
	}
	if strings.TrimSpace(query.BucketType) != "" {
		session.And("bucket_type = ?", strings.TrimSpace(query.BucketType))
	}
	if strings.TrimSpace(query.BucketNature) != "" {
		session.And("bucket_nature = ?", strings.TrimSpace(query.BucketNature))
	}
	if query.IsActive != nil {
		session.And("is_active = ?", *query.IsActive)
	}

	err := session.Asc("id").Find(&buckets)
	return buckets, err
}

func GetBucketByIdForUser(session *xorm.Session, id int64, userId int64) (*Bucket, error) {
	buckets := make([]Bucket, 0, 1)
	err := session.Where("id = ? AND user_id = ? AND is_deleted = ?", id, userId, 0).Limit(1).Find(&buckets)
	if err != nil {
		return nil, err
	}
	if len(buckets) == 0 {
		return nil, errors.New("bucket not found")
	}
	return &buckets[0], nil
}

func GetBucketByIdForUserForUpdate(session *xorm.Session, id int64, userId int64) (*Bucket, error) {
	bucket := &Bucket{}
	has, err := session.SQL("SELECT * FROM bucket WHERE id = ? AND user_id = ? AND is_deleted = 0 FOR UPDATE", id, userId).Get(bucket)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("bucket not found")
	}
	return bucket, nil
}

func GetBucket(id int64, userId int64) (*Bucket, error) {
	session := infrastructure.Mysql.NewSession()
	defer session.Close()
	return GetBucketByIdForUser(session, id, userId)
}

func UpdateBucket(bucket *Bucket) error {
	_, err := infrastructure.Mysql.
		Where("id = ? AND user_id = ? AND is_deleted = ?", bucket.Id, bucket.UserId, 0).
		Cols("name", "bucket_type", "bucket_nature", "bucket_group_key", "is_active").
		Update(bucket)
	return err
}

func UpdateBucketBalance(session *xorm.Session, id int64, userId int64, balance decimal.Decimal) error {
	_, err := session.
		Where("id = ? AND user_id = ? AND is_deleted = ?", id, userId, 0).
		Cols("balance").
		Update(&Bucket{Balance: balance})
	return err
}
