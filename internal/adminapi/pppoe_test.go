// internal/adminapi/pppoe_test.go
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

// TestPppoeProfileRequest_Validation tests PppoeProfileRequest validation.
func TestPppoeProfileRequest_Validation(t *testing.T) {
        tests := []struct {
                name    string
                request PppoeProfileRequest
                wantErr bool
        }{
                {
                        name: "valid request",
                        request: PppoeProfileRequest{
                                Name:           "Test PPPoE",
                                AddrPool:       "pppoe-pool",
                                SessionTimeout: 86400,
                                UpRate:         10240,
                                DownRate:       20480,
                        },
                        wantErr: false,
                },
                {
                        name: "missing name",
                        request: PppoeProfileRequest{
                                AddrPool: "pppoe-pool",
                        },
                        wantErr: true,
                },
                {
                        name: "invalid VLAN ID",
                        request: PppoeProfileRequest{
                                Name:    "Test",
                                Vlanid1: 5000, // Exceeds 4096
                        },
                        wantErr: true,
                },
                {
                        name: "valid VLAN IDs",
                        request: PppoeProfileRequest{
                                Name:    "Test",
                                Vlanid1: 100,
                                Vlanid2: 200,
                        },
                        wantErr: false,
                },
                {
                        name: "invalid priority",
                        request: PppoeProfileRequest{
                                Name:      "Test",
                                Priority:  10, // Exceeds 7
                        },
                        wantErr: true,
                },
                {
                        name: "valid priority",
                        request: PppoeProfileRequest{
                                Name:      "Test",
                                Priority:  5,
                        },
                        wantErr: false,
                },
                {
                        name: "invalid active_num",
                        request: PppoeProfileRequest{
                                Name:      "Test",
                                ActiveNum: 200, // Exceeds 100
                        },
                        wantErr: true,
                },
                {
                        name: "valid active_num",
                        request: PppoeProfileRequest{
                                Name:      "Test",
                                ActiveNum: 5,
                        },
                        wantErr: false,
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        e := echo.New()
                        body, _ := json.Marshal(tt.request)
                        req := httptest.NewRequest(http.MethodPost, "/api/v1/pppoe-profiles", bytes.NewReader(body))
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

// TestPppoeProfileRequest_toPppoeProfile tests the conversion method.
func TestPppoeProfileRequest_toPppoeProfile(t *testing.T) {
        tests := []struct {
                name     string
                request  PppoeProfileRequest
                expected domain.PppoeProfile
        }{
                {
                        name: "basic conversion",
                        request: PppoeProfileRequest{
                                Name:           "Test Profile",
                                AddrPool:       "pppoe-pool",
                                SessionTimeout: 86400,
                                UpRate:         10240,
                                DownRate:       20480,
                        },
                        expected: domain.PppoeProfile{
                                Name:           "Test Profile",
                                AddrPool:       "pppoe-pool",
                                SessionTimeout: 86400,
                                UpRate:         10240,
                                DownRate:       20480,
                        },
                },
                {
                        name: "with VLAN and PVC settings",
                        request: PppoeProfileRequest{
                                Name:    "Test",
                                Vlanid1: 100,
                                Vlanid2: 200,
                                PvcVPI:  8,
                                PvcVCI:  35,
                        },
                        expected: domain.PppoeProfile{
                                Name:    "Test",
                                Vlanid1: 100,
                                Vlanid2: 200,
                                PvcVPI:  8,
                                PvcVCI:  35,
                        },
                },
                {
                        name: "with burst settings",
                        request: PppoeProfileRequest{
                                Name:           "Test",
                                UpBurstRate:    20480,
                                DownBurstRate:  40960,
                                UpBurstSize:    1024,
                                DownBurstSize:  2048,
                        },
                        expected: domain.PppoeProfile{
                                Name:           "Test",
                                UpBurstRate:    20480,
                                DownBurstRate:  40960,
                                UpBurstSize:    1024,
                                DownBurstSize:  2048,
                        },
                },
                {
                        name: "with bind_mac true",
                        request: PppoeProfileRequest{
                                Name:    "Test",
                                BindMac: true,
                        },
                        expected: domain.PppoeProfile{
                                Name:    "Test",
                                BindMac: 1,
                        },
                },
                {
                        name: "with bind_vlan true",
                        request: PppoeProfileRequest{
                                Name:     "Test",
                                BindVlan: true,
                        },
                        expected: domain.PppoeProfile{
                                Name:     "Test",
                                BindVlan: 1,
                        },
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        result := tt.request.toPppoeProfile()
                        assert.Equal(t, tt.expected.Name, result.Name)
                        assert.Equal(t, tt.expected.AddrPool, result.AddrPool)
                        assert.Equal(t, tt.expected.SessionTimeout, result.SessionTimeout)
                        assert.Equal(t, tt.expected.UpRate, result.UpRate)
                        assert.Equal(t, tt.expected.DownRate, result.DownRate)
                        assert.Equal(t, tt.expected.Vlanid1, result.Vlanid1)
                        assert.Equal(t, tt.expected.Vlanid2, result.Vlanid2)
                        assert.Equal(t, tt.expected.PvcVPI, result.PvcVPI)
                        assert.Equal(t, tt.expected.PvcVCI, result.PvcVCI)
                        assert.Equal(t, tt.expected.BindMac, result.BindMac)
                        assert.Equal(t, tt.expected.BindVlan, result.BindVlan)
                })
        }
}

// TestPppoeUserRequest_Validation tests PppoeUserRequest validation.
func TestPppoeUserRequest_Validation(t *testing.T) {
        tests := []struct {
                name    string
                request PppoeUserRequest
                wantErr bool
        }{
                {
                        name: "valid request",
                        request: PppoeUserRequest{
                                ProfileId: float64(1),
                                Username:  "pppoeuser@example.com",
                                Password:  "password123",
                        },
                        wantErr: false,
                },
                {
                        name: "missing profile_id",
                        request: PppoeUserRequest{
                                Username: "testuser",
                                Password: "password123",
                        },
                        wantErr: true,
                },
                {
                        name: "missing username",
                        request: PppoeUserRequest{
                                ProfileId: float64(1),
                                Password:  "password123",
                        },
                        wantErr: true,
                },
                {
                        name: "invalid email",
                        request: PppoeUserRequest{
                                ProfileId: float64(1),
                                Username:  "testuser",
                                Email:     "not-an-email",
                        },
                        wantErr: true,
                },
                {
                        name: "invalid VLAN IDs",
                        request: PppoeUserRequest{
                                ProfileId: float64(1),
                                Username:  "testuser",
                                Vlanid1:   5000, // Exceeds 4096
                        },
                        wantErr: true,
                },
                {
                        name: "valid with IP and MAC",
                        request: PppoeUserRequest{
                                ProfileId: float64(1),
                                Username:  "testuser",
                                IpAddr:    "192.168.1.100",
                                MacAddr:   "AA:BB:CC:DD:EE:FF",
                        },
                        wantErr: false,
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        e := echo.New()
                        body, _ := json.Marshal(tt.request)
                        req := httptest.NewRequest(http.MethodPost, "/api/v1/pppoe-users", bytes.NewReader(body))
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

// TestPppoeUserRequest_toPppoeUser tests the conversion method.
func TestPppoeUserRequest_toPppoeUser(t *testing.T) {
        tests := []struct {
                name     string
                request  PppoeUserRequest
                expected domain.PppoeUser
        }{
                {
                        name: "basic conversion",
                        request: PppoeUserRequest{
                                ProfileId: float64(1),
                                Username:  "testuser@example.com",
                                Password:  "password123",
                        },
                        expected: domain.PppoeUser{
                                ProfileId: 1,
                                Username:  "testuser@example.com",
                                Password:  "password123",
                        },
                },
                {
                        name: "with VLAN settings",
                        request: PppoeUserRequest{
                                ProfileId: float64(1),
                                Username:  "testuser",
                                Vlanid1:   100,
                                Vlanid2:   200,
                        },
                        expected: domain.PppoeUser{
                                ProfileId: 1,
                                Username:  "testuser",
                                Vlanid1:   100,
                                Vlanid2:   200,
                        },
                },
                {
                        name: "with IPv6 settings",
                        request: PppoeUserRequest{
                                ProfileId:             float64(1),
                                Username:              "testuser",
                                IPv6Addr:              "2001:db8::1",
                                DelegatedIPv6Prefix:   "2001:db8:1000::/48",
                        },
                        expected: domain.PppoeUser{
                                ProfileId:             1,
                                Username:              "testuser",
                                IPv6Addr:              "2001:db8::1",
                                DelegatedIPv6Prefix:   "2001:db8:1000::/48",
                        },
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        result := tt.request.toPppoeUser()
                        assert.Equal(t, tt.expected.ProfileId, result.ProfileId)
                        assert.Equal(t, tt.expected.Username, result.Username)
                        assert.Equal(t, tt.expected.Password, result.Password)
                        assert.Equal(t, tt.expected.Vlanid1, result.Vlanid1)
                        assert.Equal(t, tt.expected.Vlanid2, result.Vlanid2)
                        assert.Equal(t, tt.expected.IPv6Addr, result.IPv6Addr)
                        assert.Equal(t, tt.expected.DelegatedIPv6Prefix, result.DelegatedIPv6Prefix)
                })
        }
}
