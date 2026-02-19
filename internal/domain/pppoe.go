// internal/domain/pppoe.go
package domain

import "time"

// PppoeProfile represents a billing profile specifically for PPPoE authentication.
// PPPoE profiles are optimized for broadband PPPoE connections with features like
// IP pool assignment, PVC/VLAN configuration, and QoS settings.
//
// Database table: pppoe_profile
// GORM features: Auto-migration, soft delete (DeletedAt), timestamps
//
// Key differences from RadiusProfile:
//  - Support for PPPoE-specific attributes (AC name, service name)
//  - PVC/VLAN configuration for DSL networks
//  - IP pool assignment for PPP connections
//  - Interim accounting interval configuration
//  - QoS and priority settings
//
// Lifecycle:
//  1. Created via Admin API POST /api/v1/pppoe-profiles
//  2. Associated with PPPoE users
//  3. Used during PPPoE RADIUS authentication
type PppoeProfile struct {
        // ID is the auto-incrementing primary key.
        ID int64 `json:"id,string" gorm:"primaryKey" form:"id"`

        // NodeId references the network node this profile belongs to.
        NodeId int64 `json:"node_id,string" form:"node_id"`

        // Name is the display name for this PPPoE profile.
        // Must be unique across the system.
        Name string `json:"name" gorm:"uniqueIndex;size:100" form:"name"`

        // Status indicates the profile status.
        // Possible values: "enabled", "disabled"
        // Default: "enabled"
        Status string `json:"status" gorm:"default:'enabled';size:20;index" form:"status"`

        // AddrPool is the IP address pool name for PPPoE users.
        // This is used to assign IP addresses to PPP interfaces.
        AddrPool string `json:"addr_pool" form:"addr_pool"`

        // IPv6PrefixPool is the IPv6 prefix pool for delegated prefix assignment.
        IPv6PrefixPool string `json:"ipv6_prefix_pool" form:"ipv6_prefix_pool"`

        // IPv6AddrPool is the IPv6 address pool for interface ID assignment.
        IPv6AddrPool string `json:"ipv6_addr_pool" form:"ipv6_addr_pool"`

        // AcName is the Access Concentrator name filter.
        // If set, only PPPoE clients matching this AC name are allowed.
        AcName string `json:"ac_name" form:"ac_name"`

        // ServiceName is the PPPoE service name filter.
        // If set, only PPPoE clients requesting this service are allowed.
        ServiceName string `json:"service_name" form:"service_name"`

        // SessionTimeout is the maximum session duration in seconds.
        // 0 means no limit.
        SessionTimeout int `json:"session_timeout" form:"session_timeout"`

        // IdleTimeout is the idle timeout in seconds.
        // Session is terminated after this period of inactivity.
        IdleTimeout int `json:"idle_timeout" form:"idle_timeout"`

        // InterimInterval is the interim accounting interval in seconds.
        // Default: 600 (10 minutes)
        InterimInterval int `json:"interim_interval" form:"interim_interval"`

        // UpRate is the upload bandwidth limit in Kbps.
        UpRate int `json:"up_rate" form:"up_rate"`

        // DownRate is the download bandwidth limit in Kbps.
        DownRate int `json:"down_rate" form:"down_rate"`

        // UpBurstRate is the upload burst rate in Kbps.
        // For rate limiting with burst support.
        UpBurstRate int `json:"up_burst_rate" form:"up_burst_rate"`

        // DownBurstRate is the download burst rate in Kbps.
        DownBurstRate int `json:"down_burst_rate" form:"down_burst_rate"`

        // UpBurstSize is the upload burst size in KB.
        UpBurstSize int `json:"up_burst_size" form:"up_burst_size"`

        // DownBurstSize is the download burst size in KB.
        DownBurstSize int `json:"down_burst_size" form:"down_burst_size"`

        // Vlanid1 is the inner VLAN ID (C-VLAN).
        Vlanid1 int `json:"vlanid1" form:"vlanid1"`

        // Vlanid2 is the outer VLAN ID (S-VLAN).
        Vlanid2 int `json:"vlanid2" form:"vlanid2"`

        // PvcVPI is the ATM VPI (Virtual Path Identifier).
        // Used for DSL networks with ATM transport.
        PvcVPI int `json:"pvc_vpi" form:"pvc_vpi"`

        // PvcVCI is the ATM VCI (Virtual Channel Identifier).
        // Used for DSL networks with ATM transport.
        PvcVCI int `json:"pvc_vci" form:"pvc_vci"`

        // Domain is the domain for vendor-specific features.
        Domain string `json:"domain" form:"domain"`

        // BindMac enables MAC address binding.
        BindMac int `json:"bind_mac" form:"bind_mac"`

        // BindVlan enables VLAN binding.
        BindVlan int `json:"bind_vlan" form:"bind_vlan"`

        // ActiveNum is the maximum number of concurrent sessions.
        ActiveNum int `json:"active_num" form:"active_num"`

        // Priority is the QoS priority level (0-7).
        // Higher values indicate higher priority.
        Priority int `json:"priority" form:"priority"`

        // Remark is an optional description for this profile.
        Remark string `json:"remark" form:"remark"`

        // CreatedAt is automatically set by GORM on INSERT.
        CreatedAt time.Time `json:"created_at"`

        // UpdatedAt is automatically updated by GORM on UPDATE.
        UpdatedAt time.Time `json:"updated_at"`

        // DeletedAt enables GORM soft delete. Non-null means record is deleted.
        DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

// TableName returns the database table name for PppoeProfile.
func (PppoeProfile) TableName() string {
        return "pppoe_profile"
}

// PppoeUser represents a user account for PPPoE authentication.
// PPPoE users connect via PPPoE clients (routers, DSL modems, etc.).
//
// Database table: pppoe_user
// GORM features: Auto-migration, soft delete (DeletedAt), timestamps
type PppoeUser struct {
        // ID is the auto-incrementing primary key.
        ID int64 `json:"id,string" gorm:"primaryKey" form:"id"`

        // NodeId references the network node.
        NodeId int64 `json:"node_id,string" form:"node_id"`

        // ProfileId references the PppoeProfile.
        ProfileId int64 `json:"profile_id,string" gorm:"index" form:"profile_id"`

        // Username is the PPPoE login name.
        // Format: typically user@realm or plain username.
        Username string `json:"username" gorm:"uniqueIndex;size:100" form:"username"`

        // Password is the PPPoE password.
        Password string `json:"password" form:"password"`

        // Realname is the user's real name.
        Realname string `json:"realname" form:"realname"`

        // Mobile is the user's mobile number.
        Mobile string `json:"mobile" form:"mobile"`

        // Email is the user's email address.
        Email string `json:"email" form:"email"`

        // Address is the user's physical address.
        Address string `json:"address" form:"address"`

        // MacAddr is the bound MAC address.
        MacAddr string `json:"mac_addr" form:"mac_addr"`

        // IpAddr is the static IPv4 address.
        // If set, this IP is assigned to the PPP interface.
        IpAddr string `json:"ip_addr" form:"ip_addr"`

        // IPv6Addr is the static IPv6 address.
        IPv6Addr string `json:"ipv6_addr" form:"ipv6_addr"`

        // DelegatedIPv6Prefix is the delegated IPv6 prefix.
        DelegatedIPv6Prefix string `json:"delegated_ipv6_prefix" form:"delegated_ipv6_prefix"`

        // Vlanid1 is the inner VLAN ID (C-VLAN).
        Vlanid1 int `json:"vlanid1" form:"vlanid1"`

        // Vlanid2 is the outer VLAN ID (S-VLAN).
        Vlanid2 int `json:"vlanid2" form:"vlanid2"`

        // Domain is the user-specific domain.
        Domain string `json:"domain" form:"domain"`

        // Status indicates the user status.
        // Possible values: "enabled", "disabled", "expired"
        Status string `json:"status" gorm:"default:'enabled';size:20;index" form:"status"`

        // ExpireTime is the account expiration timestamp.
        ExpireTime time.Time `json:"expire_time" gorm:"index" form:"expire_time"`

        // OnlineCount is the current number of online sessions.
        // This is a computed field, not stored in database.
        OnlineCount int `json:"online_count" gorm:"-:migration;<-:false"`

        // TotalSessionTime is the total session time in seconds.
        TotalSessionTime int64 `json:"total_session_time" form:"total_session_time"`

        // TotalInputBytes is the total input bytes.
        TotalInputBytes int64 `json:"total_input_bytes" form:"total_input_bytes"`

        // TotalOutputBytes is the total output bytes.
        TotalOutputBytes int64 `json:"total_output_bytes" form:"total_output_bytes"`

        // LastOnline is the timestamp of last online session.
        LastOnline time.Time `json:"last_online" form:"last_online"`

        // Remark is an optional note.
        Remark string `json:"remark" form:"remark"`

        // CreatedAt is automatically set by GORM on INSERT.
        CreatedAt time.Time `json:"created_at" gorm:"index"`

        // UpdatedAt is automatically updated by GORM on UPDATE.
        UpdatedAt time.Time `json:"updated_at"`

        // DeletedAt enables GORM soft delete.
        DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

// TableName returns the database table name for PppoeUser.
func (PppoeUser) TableName() string {
        return "pppoe_user"
}
