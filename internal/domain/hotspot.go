// internal/domain/hotspot.go
package domain

import "time"

// HotspotProfile represents a billing profile specifically for hotspot authentication.
// Hotspot profiles are optimized for wireless hotspot scenarios with features like
// MAC address authentication, session time limits, and landing page redirection.
//
// Database table: hotspot_profile
// GORM features: Auto-migration, soft delete (DeletedAt), timestamps
//
// Key differences from RadiusProfile:
//  - Support for MAC address authentication (MAC-Auth)
//  - Session timeout in minutes (typical for hotspot scenarios)
//  - Idle timeout for inactive sessions
//  - Landing page and welcome page URL support
//  - WISPr attributes for captive portal integration
//
// Lifecycle:
//  1. Created via Admin API POST /api/v1/hotspot-profiles
//  2. Associated with hotspot users or MAC addresses
//  3. Used during RADIUS authentication for hotspot NAS devices
type HotspotProfile struct {
        // ID is the auto-incrementing primary key.
        ID int64 `json:"id,string" gorm:"primaryKey" form:"id"`

        // NodeId references the network node this profile belongs to.
        NodeId int64 `json:"node_id,string" form:"node_id"`

        // Name is the display name for this hotspot profile.
        // Must be unique across the system.
        Name string `json:"name" gorm:"uniqueIndex;size:100" form:"name"`

        // Status indicates the profile status.
        // Possible values: "enabled", "disabled"
        // Default: "enabled"
        Status string `json:"status" gorm:"default:'enabled';size:20;index" form:"status"`

        // AuthMode specifies the authentication mode.
        // Possible values: "userpass", "mac", "mac-userpass"
        // - userpass: Username/password authentication
        // - mac: MAC address authentication (password = MAC address)
        // - mac-userpass: Both MAC and username/password required
        AuthMode string `json:"auth_mode" gorm:"default:'userpass';size:20" form:"auth_mode"`

        // SessionTimeout is the maximum session duration in minutes.
        // 0 means no limit.
        SessionTimeout int `json:"session_timeout" form:"session_timeout"`

        // IdleTimeout is the idle timeout in minutes.
        // Session is terminated after this period of inactivity.
        // 0 means no idle timeout.
        IdleTimeout int `json:"idle_timeout" form:"idle_timeout"`

        // DailyLimit is the daily usage limit in minutes.
        // 0 means no daily limit.
        DailyLimit int `json:"daily_limit" form:"daily_limit"`

        // MonthlyLimit is the monthly usage limit in minutes.
        // 0 means no monthly limit.
        MonthlyLimit int `json:"monthly_limit" form:"monthly_limit"`

        // UpRate is the upload bandwidth limit in Kbps.
        UpRate int `json:"up_rate" form:"up_rate"`

        // DownRate is the download bandwidth limit in Kbps.
        DownRate int `json:"down_rate" form:"down_rate"`

        // UpLimit is the total upload data limit in MB.
        // 0 means no limit.
        UpLimit int64 `json:"up_limit" form:"up_limit"`

        // DownLimit is the total download data limit in MB.
        // 0 means no limit.
        DownLimit int64 `json:"down_limit" form:"down_limit"`

        // TotalLimit is the total data transfer limit in MB.
        // 0 means no limit.
        TotalLimit int64 `json:"total_limit" form:"total_limit"`

        // AddrPool is the IP address pool name for hotspot users.
        AddrPool string `json:"addr_pool" form:"addr_pool"`

        // Domain is the domain for vendor-specific features.
        Domain string `json:"domain" form:"domain"`

        // WelcomeUrl is the URL to redirect users after successful authentication.
        // Used for captive portal welcome pages.
        WelcomeUrl string `json:"welcome_url" form:"welcome_url"`

        // LogoutUrl is the URL to redirect users after logout.
        LogoutUrl string `json:"logout_url" form:"logout_url"`

        // BindMac enables MAC address binding.
        // If enabled, the MAC address used during first authentication is bound.
        BindMac int `json:"bind_mac" form:"bind_mac"`

        // MaxDevices is the maximum number of concurrent devices allowed.
        // 0 means unlimited.
        MaxDevices int `json:"max_devices" form:"max_devices"`

        // Remark is an optional description for this profile.
        Remark string `json:"remark" form:"remark"`

        // CreatedAt is automatically set by GORM on INSERT.
        CreatedAt time.Time `json:"created_at"`

        // UpdatedAt is automatically updated by GORM on UPDATE.
        UpdatedAt time.Time `json:"updated_at"`

        // DeletedAt enables GORM soft delete. Non-null means record is deleted.
        DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

// TableName returns the database table name for HotspotProfile.
func (HotspotProfile) TableName() string {
        return "hotspot_profile"
}

// HotspotUser represents a user account for hotspot authentication.
// Hotspot users can be created from vouchers or registered directly.
//
// Database table: hotspot_user
// GORM features: Auto-migration, soft delete (DeletedAt), timestamps
type HotspotUser struct {
        // ID is the auto-incrementing primary key.
        ID int64 `json:"id,string" gorm:"primaryKey" form:"id"`

        // NodeId references the network node.
        NodeId int64 `json:"node_id,string" form:"node_id"`

        // ProfileId references the HotspotProfile.
        ProfileId int64 `json:"profile_id,string" gorm:"index" form:"profile_id"`

        // Username is the login name (or MAC address for MAC auth).
        Username string `json:"username" gorm:"uniqueIndex;size:50" form:"username"`

        // Password is the authentication password.
        Password string `json:"password" form:"password"`

        // Realname is the user's real name.
        Realname string `json:"realname" form:"realname"`

        // Mobile is the user's mobile number.
        Mobile string `json:"mobile" form:"mobile"`

        // Email is the user's email address.
        Email string `json:"email" form:"email"`

        // MacAddr is the bound MAC address (if MAC binding is enabled).
        MacAddr string `json:"mac_addr" form:"mac_addr"`

        // IpAddr is the static IP address (if assigned).
        IpAddr string `json:"ip_addr" form:"ip_addr"`

        // Status indicates the user status.
        // Possible values: "enabled", "disabled", "expired"
        Status string `json:"status" gorm:"default:'enabled';size:20;index" form:"status"`

        // ExpireTime is the account expiration timestamp.
        ExpireTime time.Time `json:"expire_time" gorm:"index" form:"expire_time"`

        // OnlineCount is the current number of online sessions.
        // This is a computed field, not stored in database.
        OnlineCount int `json:"online_count" gorm:"-:migration;<-:false"`

        // TotalSessionTime is the total session time in minutes.
        TotalSessionTime int64 `json:"total_session_time" form:"total_session_time"`

        // TotalInputBytes is the total input bytes.
        TotalInputBytes int64 `json:"total_input_bytes" form:"total_input_bytes"`

        // TotalOutputBytes is the total output bytes.
        TotalOutputBytes int64 `json:"total_output_bytes" form:"total_output_bytes"`

        // VoucherId references the Voucher if this user was created from a voucher.
        VoucherId int64 `json:"voucher_id,string" form:"voucher_id"`

        // Remark is an optional note.
        Remark string `json:"remark" form:"remark"`

        // LastOnline is the timestamp of last online session.
        LastOnline time.Time `json:"last_online" form:"last_online"`

        // CreatedAt is automatically set by GORM on INSERT.
        CreatedAt time.Time `json:"created_at" gorm:"index"`

        // UpdatedAt is automatically updated by GORM on UPDATE.
        UpdatedAt time.Time `json:"updated_at"`

        // DeletedAt enables GORM soft delete.
        DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

// TableName returns the database table name for HotspotUser.
func (HotspotUser) TableName() string {
        return "hotspot_user"
}

// Hotspot authentication mode constants
const (
        HotspotAuthModeUserPass    = "userpass"
        HotspotAuthModeMAC         = "mac"
        HotspotAuthModeMACUserPass = "mac-userpass"
)
