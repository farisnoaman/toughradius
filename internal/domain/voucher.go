// internal/domain/voucher.go
package domain

import (
        "time"
)

// VoucherBatch represents a batch of prepaid vouchers with common settings.
// It is used to generate multiple vouchers at once for ISP prepaid services.
//
// Database table: voucher_batch
// GORM features: Auto-migration, soft delete (DeletedAt), timestamps
//
// Lifecycle:
//  1. Created via Admin API POST /api/v1/voucher-batches
//  2. Vouchers are generated based on batch settings
//  3. Vouchers can be distributed to customers
//  4. Batch can be disabled to prevent voucher redemption
//
// Business Rules:
//  - Voucher codes are auto-generated with configurable format
//  - Each batch creates vouchers linked to this batch record
//  - Batch can be linked to a profile for automatic user creation
type VoucherBatch struct {
        // ID is the auto-incrementing primary key.
        ID int64 `json:"id,string" gorm:"primaryKey" form:"id"`

        // NodeId references the network node this batch belongs to.
        NodeId int64 `json:"node_id,string" form:"node_id"`

        // Name is the display name for this voucher batch.
        // Must be unique across the system.
        Name string `json:"name" gorm:"uniqueIndex;size:100" form:"name"`

        // ProfileId references the billing profile for voucher users.
        // When a voucher is redeemed, a user is created with this profile.
        ProfileId int64 `json:"profile_id,string" gorm:"index" form:"profile_id"`

        // TotalCount is the total number of vouchers in this batch.
        TotalCount int `json:"total_count" form:"total_count"`

        // UsedCount is the number of vouchers that have been redeemed.
        UsedCount int `json:"used_count" form:"used_count"`

        // ExpireTime is the expiration timestamp for all vouchers in this batch.
        // Vouchers cannot be redeemed after this time.
        ExpireTime time.Time `json:"expire_time" gorm:"index" form:"expire_time"`

        // ValidDays is the number of days a voucher is valid after redemption.
        // If 0, the voucher uses the batch expire time.
        ValidDays int `json:"valid_days" form:"valid_days"`

        // Prefix is the prefix for generated voucher codes.
        Prefix string `json:"prefix" gorm:"size:10" form:"prefix"`

        // CodeLength is the length of random characters in voucher codes.
        CodeLength int `json:"code_length" form:"code_length"`

        // Status indicates the batch status.
        // Possible values: "enabled", "disabled"
        // Default: "enabled"
        Status string `json:"status" gorm:"default:'enabled';size:20;index" form:"status"`

        // Remark is an optional description for this batch.
        Remark string `json:"remark" form:"remark"`

        // CreatedAt is automatically set by GORM on INSERT.
        CreatedAt time.Time `json:"created_at"`

        // UpdatedAt is automatically updated by GORM on UPDATE.
        UpdatedAt time.Time `json:"updated_at"`

        // DeletedAt enables GORM soft delete. Non-null means record is deleted.
        DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

// TableName returns the database table name for VoucherBatch.
func (VoucherBatch) TableName() string {
        return "voucher_batch"
}

// Voucher represents a single prepaid voucher code.
// Each voucher can be redeemed to create or extend a user account.
//
// Database table: voucher
// GORM features: Auto-migration, soft delete (DeletedAt), timestamps
//
// Lifecycle:
//  1. Generated as part of a VoucherBatch
//  2. Can be distributed to customers (status: "available")
//  3. Redeemed by customer (status: "used")
//  4. May expire if not used in time (status: "expired")
//
// Business Rules:
//  - Voucher code must be unique
//  - Once used, voucher is linked to the created user
//  - Expired vouchers cannot be redeemed
type Voucher struct {
        // ID is the auto-incrementing primary key.
        ID int64 `json:"id,string" gorm:"primaryKey" form:"id"`

        // BatchId references the VoucherBatch this voucher belongs to.
        BatchId int64 `json:"batch_id,string" gorm:"index" form:"batch_id"`

        // Code is the unique voucher code used for redemption.
        Code string `json:"code" gorm:"uniqueIndex;size:50" form:"code"`

        // Password is an optional password for additional security.
        // If set, both code and password must be provided for redemption.
        Password string `json:"password,omitempty" form:"password"`

        // ProfileId references the billing profile (inherited from batch).
        ProfileId int64 `json:"profile_id,string" gorm:"index" form:"profile_id"`

        // Status indicates the voucher status.
        // Possible values: "available", "used", "expired", "disabled"
        // Default: "available"
        Status string `json:"status" gorm:"default:'available';size:20;index" form:"status"`

        // UserId references the RadiusUser created when this voucher was redeemed.
        // Null if the voucher has not been redeemed.
        UserId int64 `json:"user_id,string" form:"user_id"`

        // RedeemedAt is the timestamp when the voucher was redeemed.
        RedeemedAt *time.Time `json:"redeemed_at" form:"redeemed_at"`

        // ExpireTime is the voucher expiration timestamp.
        // If null, the batch expiration time is used.
        ExpireTime *time.Time `json:"expire_time" form:"expire_time"`

        // Remark is an optional note for this voucher.
        Remark string `json:"remark" form:"remark"`

        // CreatedAt is automatically set by GORM on INSERT.
        CreatedAt time.Time `json:"created_at"`

        // UpdatedAt is automatically updated by GORM on UPDATE.
        UpdatedAt time.Time `json:"updated_at"`

        // DeletedAt enables GORM soft delete. Non-null means record is deleted.
        DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

// TableName returns the database table name for Voucher.
func (Voucher) TableName() string {
        return "voucher"
}

// Voucher status constants
const (
        VoucherStatusAvailable = "available"
        VoucherStatusUsed      = "used"
        VoucherStatusExpired   = "expired"
        VoucherStatusDisabled  = "disabled"
)

// Voucher batch status constants
const (
        VoucherBatchStatusEnabled  = "enabled"
        VoucherBatchStatusDisabled = "disabled"
)
