// internal/adminapi/hotspot_test.go
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

// TestHotspotProfileRequest_Validation tests HotspotProfileRequest validation.
func TestHotspotProfileRequest_Validation(t *testing.T) {
        tests := []struct {
                name    string
                request HotspotProfileRequest
                wantErr bool
        }{
                {
                        name: "valid request",
                        request: HotspotProfileRequest{
                                Name:           "Test Hotspot",
                                AuthMode:       domain.HotspotAuthModeUserPass,
                                SessionTimeout: 60,
                                UpRate:         1024,
                                DownRate:       2048,
                        },
                        wantErr: false,
                },
                {
                        name: "missing name",
                        request: HotspotProfileRequest{
                                AuthMode: domain.HotspotAuthModeUserPass,
                        },
                        wantErr: true,
                },
                {
                        name: "invalid auth mode",
                        request: HotspotProfileRequest{
                                Name:     "Test",
                                AuthMode: "invalid",
                        },
                        wantErr: true,
                },
                {
                        name: "valid MAC auth mode",
                        request: HotspotProfileRequest{
                                Name:     "Test MAC Auth",
                                AuthMode: domain.HotspotAuthModeMAC,
                        },
                        wantErr: false,
                },
                {
                        name: "invalid max devices",
                        request: HotspotProfileRequest{
                                Name:       "Test",
                                MaxDevices: 200, // Exceeds limit
                        },
                        wantErr: true,
                },
                {
                        name: "invalid URL",
                        request: HotspotProfileRequest{
                                Name:       "Test",
                                WelcomeUrl: "not-a-url",
                        },
                        wantErr: true,
                },
                {
                        name: "valid URLs",
                        request: HotspotProfileRequest{
                                Name:       "Test",
                                WelcomeUrl: "https://example.com/welcome",
                                LogoutUrl:  "https://example.com/logout",
                        },
                        wantErr: false,
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        e := echo.New()
                        body, _ := json.Marshal(tt.request)
                        req := httptest.NewRequest(http.MethodPost, "/api/v1/hotspot-profiles", bytes.NewReader(body))
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

// TestHotspotProfileRequest_toHotspotProfile tests the conversion method.
func TestHotspotProfileRequest_toHotspotProfile(t *testing.T) {
        tests := []struct {
                name     string
                request  HotspotProfileRequest
                expected domain.HotspotProfile
        }{
                {
                        name: "basic conversion",
                        request: HotspotProfileRequest{
                                Name:           "Test Profile",
                                AuthMode:       domain.HotspotAuthModeUserPass,
                                SessionTimeout: 60,
                                UpRate:         1024,
                                DownRate:       2048,
                        },
                        expected: domain.HotspotProfile{
                                Name:           "Test Profile",
                                AuthMode:       domain.HotspotAuthModeUserPass,
                                SessionTimeout: 60,
                                UpRate:         1024,
                                DownRate:       2048,
                        },
                },
                {
                        name: "with node_id",
                        request: HotspotProfileRequest{
                                Name:   "Test",
                                NodeId: float64(123),
                        },
                        expected: domain.HotspotProfile{
                                Name:   "Test",
                                NodeId: 123,
                        },
                },
                {
                        name: "with bind_mac true",
                        request: HotspotProfileRequest{
                                Name:    "Test",
                                BindMac: true,
                        },
                        expected: domain.HotspotProfile{
                                Name:    "Test",
                                BindMac: 1,
                        },
                },
                {
                        name: "with bind_mac false",
                        request: HotspotProfileRequest{
                                Name:    "Test",
                                BindMac: false,
                        },
                        expected: domain.HotspotProfile{
                                Name:    "Test",
                                BindMac: 0,
                        },
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        result := tt.request.toHotspotProfile()
                        assert.Equal(t, tt.expected.Name, result.Name)
                        assert.Equal(t, tt.expected.AuthMode, result.AuthMode)
                        assert.Equal(t, tt.expected.SessionTimeout, result.SessionTimeout)
                        assert.Equal(t, tt.expected.UpRate, result.UpRate)
                        assert.Equal(t, tt.expected.DownRate, result.DownRate)
                        assert.Equal(t, tt.expected.NodeId, result.NodeId)
                        assert.Equal(t, tt.expected.BindMac, result.BindMac)
                })
        }
}

// TestHotspotUserRequest_Validation tests HotspotUserRequest validation.
func TestHotspotUserRequest_Validation(t *testing.T) {
        tests := []struct {
                name    string
                request HotspotUserRequest
                wantErr bool
        }{
                {
                        name: "valid request",
                        request: HotspotUserRequest{
                                ProfileId: float64(1),
                                Username:  "testuser",
                                Password:  "password123",
                        },
                        wantErr: false,
                },
                {
                        name: "missing profile_id",
                        request: HotspotUserRequest{
                                Username: "testuser",
                                Password: "password123",
                        },
                        wantErr: true,
                },
                {
                        name: "missing username",
                        request: HotspotUserRequest{
                                ProfileId: float64(1),
                                Password:  "password123",
                        },
                        wantErr: true,
                },
                {
                        name: "short username",
                        request: HotspotUserRequest{
                                ProfileId: float64(1),
                                Username:  "ab", // Less than 3 chars
                                Password:  "password123",
                        },
                        wantErr: true,
                },
                {
                        name: "invalid email",
                        request: HotspotUserRequest{
                                ProfileId: float64(1),
                                Username:  "testuser",
                                Email:     "not-an-email",
                        },
                        wantErr: true,
                },
                {
                        name: "invalid MAC address",
                        request: HotspotUserRequest{
                                ProfileId: float64(1),
                                Username:  "testuser",
                                MacAddr:   "invalid-mac",
                        },
                        wantErr: true,
                },
                {
                        name: "valid MAC address",
                        request: HotspotUserRequest{
                                ProfileId: float64(1),
                                Username:  "testuser",
                                MacAddr:   "AA:BB:CC:DD:EE:FF",
                        },
                        wantErr: false,
                },
                {
                        name: "invalid IP address",
                        request: HotspotUserRequest{
                                ProfileId: float64(1),
                                Username:  "testuser",
                                IpAddr:    "invalid-ip",
                        },
                        wantErr: true,
                },
                {
                        name: "valid IP address",
                        request: HotspotUserRequest{
                                ProfileId: float64(1),
                                Username:  "testuser",
                                IpAddr:    "192.168.1.100",
                        },
                        wantErr: false,
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        e := echo.New()
                        body, _ := json.Marshal(tt.request)
                        req := httptest.NewRequest(http.MethodPost, "/api/v1/hotspot-users", bytes.NewReader(body))
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
