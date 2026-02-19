// internal/domain/pppoe_test.go
package domain

import (
        "encoding/json"
        "testing"
        "time"
)

// TestPppoeProfile_TableName tests the table name method.
func TestPppoeProfile_TableName(t *testing.T) {
        profile := PppoeProfile{}
        if profile.TableName() != "pppoe_profile" {
                t.Errorf("Expected table name 'pppoe_profile', got '%s'", profile.TableName())
        }
}

// TestPppoeUser_TableName tests the table name method.
func TestPppoeUser_TableName(t *testing.T) {
        user := PppoeUser{}
        if user.TableName() != "pppoe_user" {
                t.Errorf("Expected table name 'pppoe_user', got '%s'", user.TableName())
        }
}

// TestPppoeProfile_MarshalJSON tests JSON marshaling of PppoeProfile.
func TestPppoeProfile_MarshalJSON(t *testing.T) {
        profile := PppoeProfile{
                ID:              1,
                Name:            "Test PPPoE Profile",
                Status:          "enabled",
                AddrPool:        "pppoe-pool",
                IPv6PrefixPool:  "ipv6-prefix-pool",
                SessionTimeout:  86400,
                IdleTimeout:     300,
                InterimInterval: 600,
                UpRate:          10240,
                DownRate:        20480,
                Vlanid1:         100,
                Vlanid2:         200,
                ActiveNum:       1,
                Priority:        5,
        }

        data, err := json.Marshal(profile)
        if err != nil {
                t.Fatalf("Failed to marshal PppoeProfile: %v", err)
        }

        var result map[string]interface{}
        if err := json.Unmarshal(data, &result); err != nil {
                t.Fatalf("Failed to unmarshal JSON: %v", err)
        }

        if result["name"] != "Test PPPoE Profile" {
                t.Errorf("Expected name 'Test PPPoE Profile', got '%v'", result["name"])
        }
        if result["addr_pool"] != "pppoe-pool" {
                t.Errorf("Expected addr_pool 'pppoe-pool', got '%v'", result["addr_pool"])
        }
        if result["interim_interval"] != float64(600) {
                t.Errorf("Expected interim_interval 600, got '%v'", result["interim_interval"])
        }
}

// TestPppoeUser_MarshalJSON tests JSON marshaling of PppoeUser.
func TestPppoeUser_MarshalJSON(t *testing.T) {
        expireTime := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
        user := PppoeUser{
                ID:         1,
                Username:   "pppoeuser@example.com",
                ProfileId:  1,
                IpAddr:     "192.168.1.100",
                Vlanid1:    100,
                ExpireTime: expireTime,
                Status:     "enabled",
        }

        data, err := json.Marshal(user)
        if err != nil {
                t.Fatalf("Failed to marshal PppoeUser: %v", err)
        }

        var result map[string]interface{}
        if err := json.Unmarshal(data, &result); err != nil {
                t.Fatalf("Failed to unmarshal JSON: %v", err)
        }

        if result["username"] != "pppoeuser@example.com" {
                t.Errorf("Expected username 'pppoeuser@example.com', got '%v'", result["username"])
        }
        if result["ip_addr"] != "192.168.1.100" {
                t.Errorf("Expected ip_addr '192.168.1.100', got '%v'", result["ip_addr"])
        }
}

// TestPppoeProfile_DefaultValues tests default value assignments.
func TestPppoeProfile_DefaultValues(t *testing.T) {
        profile := PppoeProfile{}

        // Test that zero values are properly handled
        if profile.InterimInterval != 0 {
                t.Errorf("Expected default InterimInterval 0, got %d", profile.InterimInterval)
        }
        if profile.ActiveNum != 0 {
                t.Errorf("Expected default ActiveNum 0, got %d", profile.ActiveNum)
        }
}

// TestPppoeUser_VLANAssignment tests VLAN assignment logic.
func TestPppoeUser_VLANAssignment(t *testing.T) {
        user := PppoeUser{
                Vlanid1: 100,
                Vlanid2: 200,
        }

        if user.Vlanid1 != 100 {
                t.Errorf("Expected Vlanid1 100, got %d", user.Vlanid1)
        }
        if user.Vlanid2 != 200 {
                t.Errorf("Expected Vlanid2 200, got %d", user.Vlanid2)
        }
}
