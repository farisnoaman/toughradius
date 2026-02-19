// internal/adminapi/vouchers_test.go
package adminapi

import (
        "bytes"
        "encoding/json"
        "net/http"
        "net/http/httptest"
        "testing"

        "github.com/labstack/echo/v4"
        "github.com/stretchr/testify/assert"
        "github.com/talkincode/toughradius/v9/internal/domain"
)

// TestVoucherBatchRequest_Validation tests VoucherBatchRequest validation.
func TestVoucherBatchRequest_Validation(t *testing.T) {
        tests := []struct {
                name    string
                request VoucherBatchRequest
                wantErr bool
        }{
                {
                        name: "valid request",
                        request: VoucherBatchRequest{
                                Name:       "Test Batch",
                                ProfileId:  float64(1),
                                TotalCount: 10,
                                ExpireTime: "2025-12-31 23:59:59",
                        },
                        wantErr: false,
                },
                {
                        name: "missing name",
                        request: VoucherBatchRequest{
                                ProfileId:  float64(1),
                                TotalCount: 10,
                                ExpireTime: "2025-12-31 23:59:59",
                        },
                        wantErr: true,
                },
                {
                        name: "invalid total count",
                        request: VoucherBatchRequest{
                                Name:       "Test Batch",
                                ProfileId:  float64(1),
                                TotalCount: 0,
                                ExpireTime: "2025-12-31 23:59:59",
                        },
                        wantErr: true,
                },
                {
                        name: "total count exceeds limit",
                        request: VoucherBatchRequest{
                                Name:       "Test Batch",
                                ProfileId:  float64(1),
                                TotalCount: 100001,
                                ExpireTime: "2025-12-31 23:59:59",
                        },
                        wantErr: true,
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        // Create echo context for validation
                        e := echo.New()
                        body, _ := json.Marshal(tt.request)
                        req := httptest.NewRequest(http.MethodPost, "/api/v1/voucher-batches", bytes.NewReader(body))
                        req.Header.Set("Content-Type", "application/json")
                        rec := httptest.NewRecorder()
                        c := e.NewContext(req, rec)

                        err := c.Validate(&tt.request)
                        if tt.wantErr {
                                assert.Error(t, err)
                        } else {
                                assert.NoError(t, err)
                        }
                })
        }
}

// TestVoucherBatchRequest_toVoucherBatch tests the conversion method.
func TestVoucherBatchRequest_toVoucherBatch(t *testing.T) {
        tests := []struct {
                name     string
                request  VoucherBatchRequest
                expected domain.VoucherBatch
        }{
                {
                        name: "string profile_id",
                        request: VoucherBatchRequest{
                                Name:       "Test",
                                ProfileId:  "123",
                                TotalCount: 10,
                                CodeLength: 10,
                        },
                        expected: domain.VoucherBatch{
                                Name:       "Test",
                                ProfileId:  123,
                                TotalCount: 10,
                                CodeLength: 10,
                        },
                },
                {
                        name: "float profile_id",
                        request: VoucherBatchRequest{
                                Name:       "Test",
                                ProfileId:  float64(456),
                                TotalCount: 20,
                                CodeLength: 8,
                        },
                        expected: domain.VoucherBatch{
                                Name:       "Test",
                                ProfileId:  456,
                                TotalCount: 20,
                                CodeLength: 8,
                        },
                },
                {
                        name: "boolean status true",
                        request: VoucherBatchRequest{
                                Name:       "Test",
                                ProfileId:  float64(1),
                                TotalCount: 10,
                                Status:     true,
                        },
                        expected: domain.VoucherBatch{
                                Name:       "Test",
                                ProfileId:  1,
                                TotalCount: 10,
                                Status:     domain.VoucherBatchStatusEnabled,
                        },
                },
                {
                        name: "boolean status false",
                        request: VoucherBatchRequest{
                                Name:       "Test",
                                ProfileId:  float64(1),
                                TotalCount: 10,
                                Status:     false,
                        },
                        expected: domain.VoucherBatch{
                                Name:       "Test",
                                ProfileId:  1,
                                TotalCount: 10,
                                Status:     domain.VoucherBatchStatusDisabled,
                        },
                },
                {
                        name: "default code length",
                        request: VoucherBatchRequest{
                                Name:       "Test",
                                ProfileId:  float64(1),
                                TotalCount: 10,
                                CodeLength: 5, // Below minimum
                        },
                        expected: domain.VoucherBatch{
                                Name:       "Test",
                                ProfileId:  1,
                                TotalCount: 10,
                                CodeLength: 10, // Should be set to default
                        },
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        result := tt.request.toVoucherBatch()
                        assert.Equal(t, tt.expected.Name, result.Name)
                        assert.Equal(t, tt.expected.ProfileId, result.ProfileId)
                        assert.Equal(t, tt.expected.TotalCount, result.TotalCount)
                        assert.Equal(t, tt.expected.Status, result.Status)
                        assert.GreaterOrEqual(t, result.CodeLength, 6) // Minimum code length
                })
        }
}

// TestGenerateVoucherCode tests voucher code generation.
func TestGenerateVoucherCode(t *testing.T) {
        tests := []struct {
                name       string
                prefix     string
                length     int
                wantLen    int
                wantPrefix string
        }{
                {
                        name:       "with prefix",
                        prefix:     "VIP",
                        length:     10,
                        wantLen:    13, // prefix + code
                        wantPrefix: "VIP",
                },
                {
                        name:       "without prefix",
                        prefix:     "",
                        length:     10,
                        wantLen:    10,
                        wantPrefix: "",
                },
                {
                        name:       "long code",
                        prefix:     "TEST",
                        length:     20,
                        wantLen:    24,
                        wantPrefix: "TEST",
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        code, err := generateVoucherCode(tt.prefix, tt.length)
                        assert.NoError(t, err)
                        assert.Len(t, code, tt.wantLen)
                        if tt.wantPrefix != "" {
                                assert.True(t, len(code) >= len(tt.wantPrefix))
                        }
                })
        }
}

// TestGenerateVoucherCode_Uniqueness tests that generated codes are unique.
func TestGenerateVoucherCode_Uniqueness(t *testing.T) {
        codes := make(map[string]bool)
        for i := 0; i < 1000; i++ {
                code, err := generateVoucherCode("", 10)
                assert.NoError(t, err)
                assert.False(t, codes[code], "Generated duplicate code: %s", code)
                codes[code] = true
        }
}
