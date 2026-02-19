// internal/domain/voucher_test.go
package domain

import (
        "encoding/json"
        "testing"
        "time"
)

// TestVoucherBatch_TableName tests the table name method.
func TestVoucherBatch_TableName(t *testing.T) {
        batch := VoucherBatch{}
        if batch.TableName() != "voucher_batch" {
                t.Errorf("Expected table name 'voucher_batch', got '%s'", batch.TableName())
        }
}

// TestVoucher_TableName tests the table name method.
func TestVoucher_TableName(t *testing.T) {
        voucher := Voucher{}
        if voucher.TableName() != "voucher" {
                t.Errorf("Expected table name 'voucher', got '%s'", voucher.TableName())
        }
}

// TestVoucherBatch_MarshalJSON tests JSON marshaling of VoucherBatch.
func TestVoucherBatch_MarshalJSON(t *testing.T) {
        expireTime := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
        batch := VoucherBatch{
                ID:         1,
                Name:       "Test Batch",
                TotalCount: 100,
                UsedCount:  10,
                ExpireTime: expireTime,
                Status:     VoucherBatchStatusEnabled,
        }

        data, err := json.Marshal(batch)
        if err != nil {
                t.Fatalf("Failed to marshal VoucherBatch: %v", err)
        }

        var result map[string]interface{}
        if err := json.Unmarshal(data, &result); err != nil {
                t.Fatalf("Failed to unmarshal JSON: %v", err)
        }

        if result["name"] != "Test Batch" {
                t.Errorf("Expected name 'Test Batch', got '%v'", result["name"])
        }
        if result["total_count"] != float64(100) {
                t.Errorf("Expected total_count 100, got '%v'", result["total_count"])
        }
}

// TestVoucher_MarshalJSON tests JSON marshaling of Voucher.
func TestVoucher_MarshalJSON(t *testing.T) {
        now := time.Now()
        voucher := Voucher{
                ID:         1,
                BatchId:    1,
                Code:       "TEST123456",
                Status:     VoucherStatusAvailable,
                RedeemedAt: &now,
        }

        data, err := json.Marshal(voucher)
        if err != nil {
                t.Fatalf("Failed to marshal Voucher: %v", err)
        }

        var result map[string]interface{}
        if err := json.Unmarshal(data, &result); err != nil {
                t.Fatalf("Failed to unmarshal JSON: %v", err)
        }

        if result["code"] != "TEST123456" {
                t.Errorf("Expected code 'TEST123456', got '%v'", result["code"])
        }
        if result["status"] != VoucherStatusAvailable {
                t.Errorf("Expected status '%s', got '%v'", VoucherStatusAvailable, result["status"])
        }
}

// TestVoucherStatusConstants tests the voucher status constants.
func TestVoucherStatusConstants(t *testing.T) {
        tests := []struct {
                name     string
                constant string
                expected string
        }{
                {"Available", VoucherStatusAvailable, "available"},
                {"Used", VoucherStatusUsed, "used"},
                {"Expired", VoucherStatusExpired, "expired"},
                {"Disabled", VoucherStatusDisabled, "disabled"},
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        if tt.constant != tt.expected {
                                t.Errorf("Expected constant '%s', got '%s'", tt.expected, tt.constant)
                        }
                })
        }
}

// TestVoucherBatchStatusConstants tests the batch status constants.
func TestVoucherBatchStatusConstants(t *testing.T) {
        if VoucherBatchStatusEnabled != "enabled" {
                t.Errorf("Expected 'enabled', got '%s'", VoucherBatchStatusEnabled)
        }
        if VoucherBatchStatusDisabled != "disabled" {
                t.Errorf("Expected 'disabled', got '%s'", VoucherBatchStatusDisabled)
        }
}
