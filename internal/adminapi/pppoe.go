// internal/adminapi/pppoe.go
package adminapi

import (
        "errors"
        "net/http"
        "strconv"
        "strings"
        "time"

        "github.com/labstack/echo/v4"
        "gorm.io/gorm"

        "github.com/talkincode/toughradius/v9/internal/domain"
        "github.com/talkincode/toughradius/v9/internal/webserver"
        "github.com/talkincode/toughradius/v9/pkg/common"
)

// PppoeProfileRequest represents the request body for creating a PPPoE profile.
type PppoeProfileRequest struct {
        Name           string      `json:"name" validate:"required,min=1,max=100"`
        NodeId         interface{} `json:"node_id"`
        Status         interface{} `json:"status"`
        AddrPool       string      `json:"addr_pool" validate:"omitempty,max=50"`
        IPv6PrefixPool string      `json:"ipv6_prefix_pool" validate:"omitempty,max=50"`
        IPv6AddrPool   string      `json:"ipv6_addr_pool" validate:"omitempty,max=50"`
        AcName         string      `json:"ac_name" validate:"omitempty,max=50"`
        ServiceName    string      `json:"service_name" validate:"omitempty,max=50"`
        SessionTimeout int         `json:"session_timeout" validate:"gte=0"`
        IdleTimeout    int         `json:"idle_timeout" validate:"gte=0"`
        InterimInterval int        `json:"interim_interval" validate:"gte=0"`
        UpRate         int         `json:"up_rate" validate:"gte=0"`
        DownRate       int         `json:"down_rate" validate:"gte=0"`
        UpBurstRate    int         `json:"up_burst_rate" validate:"gte=0"`
        DownBurstRate  int         `json:"down_burst_rate" validate:"gte=0"`
        UpBurstSize    int         `json:"up_burst_size" validate:"gte=0"`
        DownBurstSize  int         `json:"down_burst_size" validate:"gte=0"`
        Vlanid1        int         `json:"vlanid1" validate:"gte=0,lte=4096"`
        Vlanid2        int         `json:"vlanid2" validate:"gte=0,lte=4096"`
        PvcVPI         int         `json:"pvc_vpi" validate:"gte=0"`
        PvcVCI         int         `json:"pvc_vci" validate:"gte=0"`
        Domain         string      `json:"domain" validate:"omitempty,max=50"`
        BindMac        interface{} `json:"bind_mac"`
        BindVlan       interface{} `json:"bind_vlan"`
        ActiveNum      int         `json:"active_num" validate:"gte=0,lte=100"`
        Priority       int         `json:"priority" validate:"gte=0,lte=7"`
        Remark         string      `json:"remark" validate:"omitempty,max=500"`
}

// toPppoeProfile converts PppoeProfileRequest to PppoeProfile.
func (req *PppoeProfileRequest) toPppoeProfile() *domain.PppoeProfile {
        profile := &domain.PppoeProfile{
                Name:           strings.TrimSpace(req.Name),
                AddrPool:       req.AddrPool,
                IPv6PrefixPool: req.IPv6PrefixPool,
                IPv6AddrPool:   req.IPv6AddrPool,
                AcName:         req.AcName,
                ServiceName:    req.ServiceName,
                SessionTimeout: req.SessionTimeout,
                IdleTimeout:    req.IdleTimeout,
                InterimInterval: req.InterimInterval,
                UpRate:         req.UpRate,
                DownRate:       req.DownRate,
                UpBurstRate:    req.UpBurstRate,
                DownBurstRate:  req.DownBurstRate,
                UpBurstSize:    req.UpBurstSize,
                DownBurstSize:  req.DownBurstSize,
                Vlanid1:        req.Vlanid1,
                Vlanid2:        req.Vlanid2,
                PvcVPI:         req.PvcVPI,
                PvcVCI:         req.PvcVCI,
                Domain:         req.Domain,
                ActiveNum:      req.ActiveNum,
                Priority:       req.Priority,
                Remark:         req.Remark,
        }

        // Handle node_id
        switch v := req.NodeId.(type) {
        case float64:
                profile.NodeId = int64(v)
        case string:
                if v != "" {
                        id, _ := strconv.ParseInt(v, 10, 64)
                        profile.NodeId = id
                }
        }

        // Handle status
        switch v := req.Status.(type) {
        case bool:
                if v {
                        profile.Status = common.ENABLED
                } else {
                        profile.Status = common.DISABLED
                }
        case string:
                profile.Status = strings.ToLower(v)
        }

        // Handle bind_mac
        switch v := req.BindMac.(type) {
        case bool:
                if v {
                        profile.BindMac = 1
                } else {
                        profile.BindMac = 0
                }
        case float64:
                profile.BindMac = int(v)
        }

        // Handle bind_vlan
        switch v := req.BindVlan.(type) {
        case bool:
                if v {
                        profile.BindVlan = 1
                } else {
                        profile.BindVlan = 0
                }
        case float64:
                profile.BindVlan = int(v)
        }

        return profile
}

// ListPppoeProfiles retrieves the PPPoE profile list.
// @Summary Get PPPoE profile list
// @Tags PPPoE
// @Param page query int false "Page number"
// @Param perPage query int false "Items per page"
// @Success 200 {object} ListResponse
// @Router /api/v1/pppoe-profiles [get]
func ListPppoeProfiles(c echo.Context) error {
        db := GetDB(c)

        page, perPage := parsePagination(c)
        sortField := c.QueryParam("sort")
        order := c.QueryParam("order")
        if sortField == "" {
                sortField = "id"
        }
        if order != "ASC" && order != "DESC" {
                order = "DESC"
        }

        var total int64
        var profiles []domain.PppoeProfile

        query := db.Model(&domain.PppoeProfile{})

        // Filter by name
        if name := strings.TrimSpace(c.QueryParam("name")); name != "" {
                if strings.EqualFold(db.Name(), "postgres") {
                        query = query.Where("name ILIKE ?", "%"+name+"%")
                } else {
                        query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(name)+"%")
                }
        }

        // Filter by status
        if status := strings.TrimSpace(c.QueryParam("status")); status != "" {
                query = query.Where("status = ?", status)
        }

        // Filter by addr_pool
        if addrPool := strings.TrimSpace(c.QueryParam("addr_pool")); addrPool != "" {
                query = query.Where("addr_pool = ?", addrPool)
        }

        // Filter by domain
        if domain := strings.TrimSpace(c.QueryParam("domain")); domain != "" {
                query = query.Where("domain = ?", domain)
        }

        query.Count(&total)

        offset := (page - 1) * perPage
        query.Order(sortField + " " + order).Limit(perPage).Offset(offset).Find(&profiles)

        return paged(c, profiles, total, page, perPage)
}

// GetPppoeProfile retrieves a single PPPoE profile.
// @Summary Get PPPoE profile detail
// @Tags PPPoE
// @Param id path int true "Profile ID"
// @Success 200 {object} domain.PppoeProfile
// @Router /api/v1/pppoe-profiles/{id} [get]
func GetPppoeProfile(c echo.Context) error {
        id, err := strconv.ParseInt(c.Param("id"), 10, 64)
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_ID", "Invalid profile ID", nil)
        }

        var profile domain.PppoeProfile
        if err := GetDB(c).First(&profile, id).Error; err != nil {
                return fail(c, http.StatusNotFound, "NOT_FOUND", "PPPoE profile not found", nil)
        }

        return ok(c, profile)
}

// CreatePppoeProfile creates a new PPPoE profile.
// @Summary Create PPPoE profile
// @Tags PPPoE
// @Param profile body PppoeProfileRequest true "Profile information"
// @Success 201 {object} domain.PppoeProfile
// @Router /api/v1/pppoe-profiles [post]
func CreatePppoeProfile(c echo.Context) error {
        var req PppoeProfileRequest
        if err := c.Bind(&req); err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_REQUEST", "Unable to parse request parameters", err.Error())
        }

        if err := c.Validate(&req); err != nil {
                return err
        }

        profile := req.toPppoeProfile()

        // Check name uniqueness
        var count int64
        GetDB(c).Model(&domain.PppoeProfile{}).Where("name = ?", profile.Name).Count(&count)
        if count > 0 {
                return fail(c, http.StatusConflict, "NAME_EXISTS", "Profile name already exists", nil)
        }

        // Set defaults
        if profile.Status == "" {
                profile.Status = common.ENABLED
        }
        if profile.InterimInterval == 0 {
                profile.InterimInterval = 600 // Default 10 minutes
        }
        profile.CreatedAt = time.Now()
        profile.UpdatedAt = time.Now()

        if err := GetDB(c).Create(profile).Error; err != nil {
                return fail(c, http.StatusInternalServerError, "CREATE_FAILED", "Failed to create profile", err.Error())
        }

        return ok(c, profile)
}

// UpdatePppoeProfile updates a PPPoE profile.
// @Summary Update PPPoE profile
// @Tags PPPoE
// @Param id path int true "Profile ID"
// @Param profile body PppoeProfileRequest true "Profile information"
// @Success 200 {object} domain.PppoeProfile
// @Router /api/v1/pppoe-profiles/{id} [put]
func UpdatePppoeProfile(c echo.Context) error {
        id, err := strconv.ParseInt(c.Param("id"), 10, 64)
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_ID", "Invalid profile ID", nil)
        }

        var profile domain.PppoeProfile
        if err := GetDB(c).First(&profile, id).Error; err != nil {
                return fail(c, http.StatusNotFound, "NOT_FOUND", "PPPoE profile not found", nil)
        }

        var req PppoeProfileRequest
        if err := c.Bind(&req); err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_REQUEST", "Unable to parse request parameters", err.Error())
        }

        if err := c.Validate(&req); err != nil {
                return err
        }

        updateData := req.toPppoeProfile()

        // Validate name uniqueness
        if updateData.Name != "" && updateData.Name != profile.Name {
                var count int64
                GetDB(c).Model(&domain.PppoeProfile{}).Where("name = ? AND id != ?", updateData.Name, id).Count(&count)
                if count > 0 {
                        return fail(c, http.StatusConflict, "NAME_EXISTS", "Profile name already exists", nil)
                }
        }

        // Build updates map
        updates := map[string]interface{}{
                "updated_at": time.Now(),
        }
        if updateData.Name != "" {
                updates["name"] = updateData.Name
        }
        if updateData.Status != "" {
                updates["status"] = updateData.Status
        }
        if updateData.AddrPool != "" {
                updates["addr_pool"] = updateData.AddrPool
        }
        if updateData.IPv6PrefixPool != "" {
                updates["ipv6_prefix_pool"] = updateData.IPv6PrefixPool
        }
        if updateData.IPv6AddrPool != "" {
                updates["ipv6_addr_pool"] = updateData.IPv6AddrPool
        }
        if updateData.AcName != "" {
                updates["ac_name"] = updateData.AcName
        }
        if updateData.ServiceName != "" {
                updates["service_name"] = updateData.ServiceName
        }
        updates["session_timeout"] = updateData.SessionTimeout
        updates["idle_timeout"] = updateData.IdleTimeout
        updates["interim_interval"] = updateData.InterimInterval
        updates["up_rate"] = updateData.UpRate
        updates["down_rate"] = updateData.DownRate
        updates["up_burst_rate"] = updateData.UpBurstRate
        updates["down_burst_rate"] = updateData.DownBurstRate
        updates["up_burst_size"] = updateData.UpBurstSize
        updates["down_burst_size"] = updateData.DownBurstSize
        updates["vlanid1"] = updateData.Vlanid1
        updates["vlanid2"] = updateData.Vlanid2
        updates["pvc_vpi"] = updateData.PvcVPI
        updates["pvc_vci"] = updateData.PvcVCI
        if updateData.Domain != "" {
                updates["domain"] = updateData.Domain
        }
        updates["bind_mac"] = updateData.BindMac
        updates["bind_vlan"] = updateData.BindVlan
        updates["active_num"] = updateData.ActiveNum
        updates["priority"] = updateData.Priority
        if updateData.Remark != "" {
                updates["remark"] = updateData.Remark
        }
        if updateData.NodeId > 0 {
                updates["node_id"] = updateData.NodeId
        }

        if err := GetDB(c).Model(&profile).Updates(updates).Error; err != nil {
                return fail(c, http.StatusInternalServerError, "UPDATE_FAILED", "Failed to update profile", err.Error())
        }

        GetDB(c).First(&profile, id)
        return ok(c, profile)
}

// DeletePppoeProfile deletes a PPPoE profile.
// @Summary Delete PPPoE profile
// @Tags PPPoE
// @Param id path int true "Profile ID"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/pppoe-profiles/{id} [delete]
func DeletePppoeProfile(c echo.Context) error {
        id, err := strconv.ParseInt(c.Param("id"), 10, 64)
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_ID", "Invalid profile ID", nil)
        }

        // Check for users using this profile
        var userCount int64
        GetDB(c).Model(&domain.PppoeUser{}).Where("profile_id = ?", id).Count(&userCount)
        if userCount > 0 {
                return fail(c, http.StatusConflict, "IN_USE", "Profile is in use and cannot be deleted", map[string]interface{}{
                        "user_count": userCount,
                })
        }

        if err := GetDB(c).Delete(&domain.PppoeProfile{}, id).Error; err != nil {
                return fail(c, http.StatusInternalServerError, "DELETE_FAILED", "Failed to delete profile", err.Error())
        }

        return ok(c, map[string]interface{}{
                "message": "Deletion successful",
        })
}

// PppoeUserRequest represents the request body for creating a PPPoE user.
type PppoeUserRequest struct {
        NodeId                interface{} `json:"node_id"`
        ProfileId             interface{} `json:"profile_id" validate:"required"`
        Username              string      `json:"username" validate:"required,min=1,max=100"`
        Password              string      `json:"password" validate:"omitempty,min=6,max=128"`
        Realname              string      `json:"realname" validate:"omitempty,max=100"`
        Mobile                string      `json:"mobile" validate:"omitempty,max=20"`
        Email                 string      `json:"email" validate:"omitempty,email,max=100"`
        Address               string      `json:"address" validate:"omitempty,max=255"`
        MacAddr               string      `json:"mac_addr" validate:"omitempty,mac"`
        IpAddr                string      `json:"ip_addr" validate:"omitempty,ipv4"`
        IPv6Addr              string      `json:"ipv6_addr" validate:"omitempty"`
        DelegatedIPv6Prefix   string      `json:"delegated_ipv6_prefix" validate:"omitempty"`
        Vlanid1               int         `json:"vlanid1" validate:"gte=0,lte=4096"`
        Vlanid2               int         `json:"vlanid2" validate:"gte=0,lte=4096"`
        Domain                string      `json:"domain" validate:"omitempty,max=50"`
        ExpireTime            string      `json:"expire_time"`
        Status                interface{} `json:"status"`
        Remark                string      `json:"remark" validate:"omitempty,max=500"`
}

// toPppoeUser converts PppoeUserRequest to PppoeUser.
func (req *PppoeUserRequest) toPppoeUser() *domain.PppoeUser {
        user := &domain.PppoeUser{
                Username:            strings.TrimSpace(req.Username),
                Password:            req.Password,
                Realname:            req.Realname,
                Mobile:              req.Mobile,
                Email:               req.Email,
                Address:             req.Address,
                MacAddr:             req.MacAddr,
                IpAddr:              req.IpAddr,
                IPv6Addr:            req.IPv6Addr,
                DelegatedIPv6Prefix: req.DelegatedIPv6Prefix,
                Vlanid1:             req.Vlanid1,
                Vlanid2:             req.Vlanid2,
                Domain:              req.Domain,
                Remark:              req.Remark,
        }

        // Handle profile_id
        switch v := req.ProfileId.(type) {
        case float64:
                user.ProfileId = int64(v)
        case string:
                if v != "" {
                        id, _ := strconv.ParseInt(v, 10, 64)
                        user.ProfileId = id
                }
        }

        // Handle node_id
        switch v := req.NodeId.(type) {
        case float64:
                user.NodeId = int64(v)
        case string:
                if v != "" {
                        id, _ := strconv.ParseInt(v, 10, 64)
                        user.NodeId = id
                }
        }

        // Handle status
        switch v := req.Status.(type) {
        case bool:
                if v {
                        user.Status = common.ENABLED
                } else {
                        user.Status = common.DISABLED
                }
        case string:
                user.Status = strings.ToLower(v)
        }

        return user
}

// ListPppoeUsers retrieves the PPPoE user list.
// @Summary Get PPPoE user list
// @Tags PPPoE
// @Param page query int false "Page number"
// @Param perPage query int false "Items per page"
// @Success 200 {object} ListResponse
// @Router /api/v1/pppoe-users [get]
func ListPppoeUsers(c echo.Context) error {
        db := GetDB(c)

        page, perPage := parsePagination(c)
        sortField := c.QueryParam("sort")
        order := c.QueryParam("order")
        if sortField == "" {
                sortField = "id"
        }
        if order != "ASC" && order != "DESC" {
                order = "DESC"
        }

        var total int64
        var users []domain.PppoeUser

        query := db.Model(&domain.PppoeUser{}).
                Select("pppoe_user.*, COALESCE(ro.count, 0) AS online_count").
                Joins("LEFT JOIN (SELECT username, COUNT(1) AS count FROM radius_online GROUP BY username) ro ON pppoe_user.username = ro.username")

        // Filter by username
        if username := strings.TrimSpace(c.QueryParam("username")); username != "" {
                if strings.EqualFold(db.Name(), "postgres") {
                        query = query.Where("pppoe_user.username ILIKE ?", "%"+username+"%")
                } else {
                        query = query.Where("LOWER(pppoe_user.username) LIKE ?", "%"+strings.ToLower(username)+"%")
                }
        }

        // Filter by status
        if status := strings.TrimSpace(c.QueryParam("status")); status != "" {
                query = query.Where("pppoe_user.status = ?", status)
        }

        // Filter by profile_id
        if profileId := c.QueryParam("profile_id"); profileId != "" {
                if id, err := strconv.ParseInt(profileId, 10, 64); err == nil {
                        query = query.Where("pppoe_user.profile_id = ?", id)
                }
        }

        // Filter by mac_addr
        if macAddr := strings.TrimSpace(c.QueryParam("mac_addr")); macAddr != "" {
                query = query.Where("pppoe_user.mac_addr = ?", macAddr)
        }

        // Filter by ip_addr
        if ipAddr := strings.TrimSpace(c.QueryParam("ip_addr")); ipAddr != "" {
                query = query.Where("pppoe_user.ip_addr = ?", ipAddr)
        }

        query.Count(&total)

        offset := (page - 1) * perPage
        query.Order("pppoe_user." + sortField + " " + order).Limit(perPage).Offset(offset).Find(&users)

        // Clear passwords
        for i := range users {
                users[i].Password = ""
        }

        return paged(c, users, total, page, perPage)
}

// GetPppoeUser retrieves a single PPPoE user.
// @Summary Get PPPoE user detail
// @Tags PPPoE
// @Param id path int true "User ID"
// @Success 200 {object} domain.PppoeUser
// @Router /api/v1/pppoe-users/{id} [get]
func GetPppoeUser(c echo.Context) error {
        id, err := strconv.ParseInt(c.Param("id"), 10, 64)
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", nil)
        }

        var user domain.PppoeUser
        if err := GetDB(c).First(&user, id).Error; err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return fail(c, http.StatusNotFound, "NOT_FOUND", "PPPoE user not found", nil)
                }
                return fail(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to query user", err.Error())
        }

        user.Password = ""
        return ok(c, user)
}

// CreatePppoeUser creates a new PPPoE user.
// @Summary Create PPPoE user
// @Tags PPPoE
// @Param user body PppoeUserRequest true "User information"
// @Success 201 {object} domain.PppoeUser
// @Router /api/v1/pppoe-users [post]
func CreatePppoeUser(c echo.Context) error {
        var req PppoeUserRequest
        if err := c.Bind(&req); err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_REQUEST", "Unable to parse request parameters", err.Error())
        }

        if err := c.Validate(&req); err != nil {
                return err
        }

        user := req.toPppoeUser()

        // Validate profile exists
        var profile domain.PppoeProfile
        if err := GetDB(c).Where("id = ?", user.ProfileId).First(&profile).Error; err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return fail(c, http.StatusBadRequest, "PROFILE_NOT_FOUND", "Associated PPPoE profile not found", nil)
                }
                return fail(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to query profile", err.Error())
        }

        // Check username uniqueness
        var exists int64
        GetDB(c).Model(&domain.PppoeUser{}).Where("username = ?", user.Username).Count(&exists)
        if exists > 0 {
                return fail(c, http.StatusConflict, "USERNAME_EXISTS", "Username already exists", nil)
        }

        if user.Password == "" {
                return fail(c, http.StatusBadRequest, "MISSING_PASSWORD", "Password is required", nil)
        }

        // Parse expiration time
        expire, err := parseTimeInput(req.ExpireTime, time.Now().AddDate(1, 0, 0))
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_EXPIRE_TIME", "Invalid expire time format", nil)
        }

        // Set defaults and inherit from profile
        user.ID = common.UUIDint64()
        user.ExpireTime = expire
        if user.Status == "" {
                user.Status = common.ENABLED
        }
        if user.Domain == "" {
                user.Domain = profile.Domain
        }
        if user.Vlanid1 == 0 && profile.Vlanid1 > 0 {
                user.Vlanid1 = profile.Vlanid1
        }
        if user.Vlanid2 == 0 && profile.Vlanid2 > 0 {
                user.Vlanid2 = profile.Vlanid2
        }
        user.CreatedAt = time.Now()
        user.UpdatedAt = time.Now()

        if err := GetDB(c).Create(user).Error; err != nil {
                return fail(c, http.StatusInternalServerError, "CREATE_FAILED", "Failed to create user", err.Error())
        }

        user.Password = ""
        return ok(c, user)
}

// UpdatePppoeUser updates a PPPoE user.
// @Summary Update PPPoE user
// @Tags PPPoE
// @Param id path int true "User ID"
// @Param user body PppoeUserRequest true "User information"
// @Success 200 {object} domain.PppoeUser
// @Router /api/v1/pppoe-users/{id} [put]
func UpdatePppoeUser(c echo.Context) error {
        id, err := strconv.ParseInt(c.Param("id"), 10, 64)
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", nil)
        }

        var user domain.PppoeUser
        if err := GetDB(c).First(&user, id).Error; err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return fail(c, http.StatusNotFound, "NOT_FOUND", "PPPoE user not found", nil)
                }
                return fail(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to query user", err.Error())
        }

        var req PppoeUserRequest
        if err := c.Bind(&req); err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_REQUEST", "Unable to parse request parameters", err.Error())
        }

        if err := c.Validate(&req); err != nil {
                return err
        }

        updateData := req.toPppoeUser()

        // Validate username uniqueness
        if updateData.Username != "" && updateData.Username != user.Username {
                var count int64
                GetDB(c).Model(&domain.PppoeUser{}).Where("username = ? AND id != ?", updateData.Username, id).Count(&count)
                if count > 0 {
                        return fail(c, http.StatusConflict, "USERNAME_EXISTS", "Username already exists", nil)
                }
        }

        // Build updates map
        updates := map[string]interface{}{
                "updated_at": time.Now(),
        }
        if updateData.Username != "" {
                updates["username"] = updateData.Username
        }
        if req.Password != "" {
                updates["password"] = req.Password
        }
        if updateData.Realname != "" {
                updates["realname"] = updateData.Realname
        }
        if updateData.Mobile != "" {
                updates["mobile"] = updateData.Mobile
        }
        if updateData.Email != "" {
                updates["email"] = updateData.Email
        }
        if updateData.Address != "" {
                updates["address"] = updateData.Address
        }
        if updateData.MacAddr != "" {
                updates["mac_addr"] = updateData.MacAddr
        }
        if updateData.IpAddr != "" {
                updates["ip_addr"] = updateData.IpAddr
        }
        if updateData.IPv6Addr != "" {
                updates["ipv6_addr"] = updateData.IPv6Addr
        }
        if updateData.DelegatedIPv6Prefix != "" {
                updates["delegated_ipv6_prefix"] = updateData.DelegatedIPv6Prefix
        }
        if updateData.Vlanid1 >= 0 {
                updates["vlanid1"] = updateData.Vlanid1
        }
        if updateData.Vlanid2 >= 0 {
                updates["vlanid2"] = updateData.Vlanid2
        }
        if updateData.Domain != "" {
                updates["domain"] = updateData.Domain
        }
        if updateData.Status != "" {
                updates["status"] = updateData.Status
        }
        if updateData.ProfileId > 0 {
                updates["profile_id"] = updateData.ProfileId
        }
        if updateData.NodeId > 0 {
                updates["node_id"] = updateData.NodeId
        }
        if req.ExpireTime != "" {
                expire, err := parseTimeInput(req.ExpireTime, user.ExpireTime)
                if err != nil {
                        return fail(c, http.StatusBadRequest, "INVALID_EXPIRE_TIME", "Invalid expire time format", nil)
                }
                updates["expire_time"] = expire
        }
        if updateData.Remark != "" {
                updates["remark"] = updateData.Remark
        }

        if err := GetDB(c).Model(&user).Updates(updates).Error; err != nil {
                return fail(c, http.StatusInternalServerError, "UPDATE_FAILED", "Failed to update user", err.Error())
        }

        GetDB(c).First(&user, id)
        user.Password = ""
        return ok(c, user)
}

// DeletePppoeUser deletes a PPPoE user.
// @Summary Delete PPPoE user
// @Tags PPPoE
// @Param id path int true "User ID"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/pppoe-users/{id} [delete]
func DeletePppoeUser(c echo.Context) error {
        id, err := strconv.ParseInt(c.Param("id"), 10, 64)
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", nil)
        }

        if err := GetDB(c).Delete(&domain.PppoeUser{}, id).Error; err != nil {
                return fail(c, http.StatusInternalServerError, "DELETE_FAILED", "Failed to delete user", err.Error())
        }

        return ok(c, map[string]interface{}{
                "id": id,
        })
}

// registerPppoeRoutes registers PPPoE-related routes.
func registerPppoeRoutes() {
        webserver.ApiGET("/pppoe-profiles", ListPppoeProfiles)
        webserver.ApiGET("/pppoe-profiles/:id", GetPppoeProfile)
        webserver.ApiPOST("/pppoe-profiles", CreatePppoeProfile)
        webserver.ApiPUT("/pppoe-profiles/:id", UpdatePppoeProfile)
        webserver.ApiDELETE("/pppoe-profiles/:id", DeletePppoeProfile)
        webserver.ApiGET("/pppoe-users", ListPppoeUsers)
        webserver.ApiGET("/pppoe-users/:id", GetPppoeUser)
        webserver.ApiPOST("/pppoe-users", CreatePppoeUser)
        webserver.ApiPUT("/pppoe-users/:id", UpdatePppoeUser)
        webserver.ApiDELETE("/pppoe-users/:id", DeletePppoeUser)
}
