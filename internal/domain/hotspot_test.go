// internal/domain/hotspot_test.go
package domain

import (
        "encoding/json"
        "testing"
        "time"
)

// TestHotspotProfile_TableName tests the table name method.
func TestHotspotProfile_TableName(t *testing.T) {
        profile := HotspotProfile{}
        if profile.TableName() != "hotspot_profile" {
                t.Errorf("Expected table name 'hotspot_profile', got '%s'", profile.TableName())
        }
}

// TestHotspotUser_TableName tests the table name method.
func TestHotspotUser_TableName(t *testing.T) {
        user := HotspotUser{}
        if user.TableName() != "hotspot_user" {
                t.Errorf("Expected table name 'hotspot_user', got '%s'", user.TableName())
        }
}

// TestHotspotProfile_MarshalJSON tests JSON marshaling of HotspotProfile.
func TestHotspotProfile_MarshalJSON(t *testing.T) {
        profile := HotspotProfile{
                ID:             1,
                Name:           "Test Hotspot Profile",
                Status:         "enabled",
                AuthMode:       HotspotAuthModeUserPass,
                SessionTimeout: 60,
                IdleTimeout:    10,
                UpRate:         1024,
                DownRate:       2048,
                MaxDevices:     2,
        }

        data, err := json.Marshal(profile)
        if err != nil {
                t.Fatalf("Failed to marshal HotspotProfile: %v", err)
        }

        var result map[string]interface{}
        if err := json.Unmarshal(data, &result); err != nil {
                t.Fatalf("Failed to unmarshal JSON: %v", err)
        }

        if result["name"] != "Test Hotspot Profile" {
                t.Errorf("Expected name 'Test Hotspot Profile', got '%v'", result["name"])
        }
        if result["auth_mode"] != HotspotAuthModeUserPass {
                t.Errorf("Expected auth_mode '%s', got '%v'", HotspotAuthModeUserPass, result["auth_mode"])
        }
}

// TestHotspotUser_MarshalJSON tests JSON marshaling of HotspotUser.
func TestHotspotUser_MarshalJSON(t *testing.T) {
        expireTime := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
        user := HotspotUser{
                ID:         1,
                Username:   "testuser",
                ProfileId:  1,
                ExpireTime: expireTime,
                Status:     "enabled",
        }

        data, err := json.Marshal(user)
        if err != nil {
                t.Fatalf("Failed to marshal HotspotUser: %v", err)
        }

        var result map[string]interface{}
        if err := json.Unmarshal(data, &result); err != nil {
                t.Fatalf("Failed to unmarshal JSON: %v", err)
        }

        if result["username"] != "testuser" {
                t.Errorf("Expected username 'testuser', got '%v'", result["username"])
        }
}

// TestHotspotAuthModeConstants tests the authentication mode constants.
func TestHotspotAuthModeConstants(t *testing.T) {
        tests := []struct {
                name     string
                constant string
                expected string
        }{
                {"UserPass", HotspotAuthModeUserPass, "userpass"},
                {"MAC", HotspotAuthModeMAC, "mac"},
                {"MACUserPass", HotspotAuthModeMACUserPass, "mac-userpass"},
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        if tt.constant != tt.expected {
                                t.Errorf("Expected constant '%s', got '%s'", tt.expected, tt.constant)
                        }
                })
        }
}
