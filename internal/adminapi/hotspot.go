// internal/adminapi/hotspot.go
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

// HotspotProfileRequest represents the request body for creating a hotspot profile.
type HotspotProfileRequest struct {
        Name           string      `json:"name" validate:"required,min=1,max=100"`
        NodeId         interface{} `json:"node_id"`
        Status         interface{} `json:"status"`
        AuthMode       string      `json:"auth_mode" validate:"omitempty,oneof=userpass mac mac-userpass"`
        SessionTimeout int         `json:"session_timeout" validate:"gte=0"`
        IdleTimeout    int         `json:"idle_timeout" validate:"gte=0"`
        DailyLimit     int         `json:"daily_limit" validate:"gte=0"`
        MonthlyLimit   int         `json:"monthly_limit" validate:"gte=0"`
        UpRate         int         `json:"up_rate" validate:"gte=0"`
        DownRate       int         `json:"down_rate" validate:"gte=0"`
        UpLimit        int64       `json:"up_limit" validate:"gte=0"`
        DownLimit      int64       `json:"down_limit" validate:"gte=0"`
        TotalLimit     int64       `json:"total_limit" validate:"gte=0"`
        AddrPool       string      `json:"addr_pool" validate:"omitempty,max=50"`
        Domain         string      `json:"domain" validate:"omitempty,max=50"`
        WelcomeUrl     string      `json:"welcome_url" validate:"omitempty,url,max=255"`
        LogoutUrl      string      `json:"logout_url" validate:"omitempty,url,max=255"`
        BindMac        interface{} `json:"bind_mac"`
        MaxDevices     int         `json:"max_devices" validate:"gte=0,lte=100"`
        Remark         string      `json:"remark" validate:"omitempty,max=500"`
}

// toHotspotProfile converts HotspotProfileRequest to HotspotProfile.
func (req *HotspotProfileRequest) toHotspotProfile() *domain.HotspotProfile {
        profile := &domain.HotspotProfile{
                Name:           strings.TrimSpace(req.Name),
                AuthMode:       req.AuthMode,
                SessionTimeout: req.SessionTimeout,
                IdleTimeout:    req.IdleTimeout,
                DailyLimit:     req.DailyLimit,
                MonthlyLimit:   req.MonthlyLimit,
                UpRate:         req.UpRate,
                DownRate:       req.DownRate,
                UpLimit:        req.UpLimit,
                DownLimit:      req.DownLimit,
                TotalLimit:     req.TotalLimit,
                AddrPool:       req.AddrPool,
                Domain:         req.Domain,
                WelcomeUrl:     req.WelcomeUrl,
                LogoutUrl:      req.LogoutUrl,
                MaxDevices:     req.MaxDevices,
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

        return profile
}

// ListHotspotProfiles retrieves the hotspot profile list.
// @Summary Get hotspot profile list
// @Tags Hotspot
// @Param page query int false "Page number"
// @Param perPage query int false "Items per page"
// @Success 200 {object} ListResponse
// @Router /api/v1/hotspot-profiles [get]
func ListHotspotProfiles(c echo.Context) error {
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
        var profiles []domain.HotspotProfile

        query := db.Model(&domain.HotspotProfile{})

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

        // Filter by auth_mode
        if authMode := strings.TrimSpace(c.QueryParam("auth_mode")); authMode != "" {
                query = query.Where("auth_mode = ?", authMode)
        }

        query.Count(&total)

        offset := (page - 1) * perPage
        query.Order(sortField + " " + order).Limit(perPage).Offset(offset).Find(&profiles)

        return paged(c, profiles, total, page, perPage)
}

// GetHotspotProfile retrieves a single hotspot profile.
// @Summary Get hotspot profile detail
// @Tags Hotspot
// @Param id path int true "Profile ID"
// @Success 200 {object} domain.HotspotProfile
// @Router /api/v1/hotspot-profiles/{id} [get]
func GetHotspotProfile(c echo.Context) error {
        id, err := strconv.ParseInt(c.Param("id"), 10, 64)
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_ID", "Invalid profile ID", nil)
        }

        var profile domain.HotspotProfile
        if err := GetDB(c).First(&profile, id).Error; err != nil {
                return fail(c, http.StatusNotFound, "NOT_FOUND", "Hotspot profile not found", nil)
        }

        return ok(c, profile)
}

// CreateHotspotProfile creates a new hotspot profile.
// @Summary Create hotspot profile
// @Tags Hotspot
// @Param profile body HotspotProfileRequest true "Profile information"
// @Success 201 {object} domain.HotspotProfile
// @Router /api/v1/hotspot-profiles [post]
func CreateHotspotProfile(c echo.Context) error {
        var req HotspotProfileRequest
        if err := c.Bind(&req); err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_REQUEST", "Unable to parse request parameters", err.Error())
        }

        if err := c.Validate(&req); err != nil {
                return err
        }

        profile := req.toHotspotProfile()

        // Check name uniqueness
        var count int64
        GetDB(c).Model(&domain.HotspotProfile{}).Where("name = ?", profile.Name).Count(&count)
        if count > 0 {
                return fail(c, http.StatusConflict, "NAME_EXISTS", "Profile name already exists", nil)
        }

        // Set defaults
        if profile.Status == "" {
                profile.Status = common.ENABLED
        }
        if profile.AuthMode == "" {
                profile.AuthMode = domain.HotspotAuthModeUserPass
        }
        profile.CreatedAt = time.Now()
        profile.UpdatedAt = time.Now()

        if err := GetDB(c).Create(profile).Error; err != nil {
                return fail(c, http.StatusInternalServerError, "CREATE_FAILED", "Failed to create profile", err.Error())
        }

        return ok(c, profile)
}

// UpdateHotspotProfile updates a hotspot profile.
// @Summary Update hotspot profile
// @Tags Hotspot
// @Param id path int true "Profile ID"
// @Param profile body HotspotProfileRequest true "Profile information"
// @Success 200 {object} domain.HotspotProfile
// @Router /api/v1/hotspot-profiles/{id} [put]
func UpdateHotspotProfile(c echo.Context) error {
        id, err := strconv.ParseInt(c.Param("id"), 10, 64)
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_ID", "Invalid profile ID", nil)
        }

        var profile domain.HotspotProfile
        if err := GetDB(c).First(&profile, id).Error; err != nil {
                return fail(c, http.StatusNotFound, "NOT_FOUND", "Hotspot profile not found", nil)
        }

        var req HotspotProfileRequest
        if err := c.Bind(&req); err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_REQUEST", "Unable to parse request parameters", err.Error())
        }

        if err := c.Validate(&req); err != nil {
                return err
        }

        updateData := req.toHotspotProfile()

        // Validate name uniqueness
        if updateData.Name != "" && updateData.Name != profile.Name {
                var count int64
                GetDB(c).Model(&domain.HotspotProfile{}).Where("name = ? AND id != ?", updateData.Name, id).Count(&count)
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
        if updateData.AuthMode != "" {
                updates["auth_mode"] = updateData.AuthMode
        }
        if updateData.SessionTimeout >= 0 {
                updates["session_timeout"] = updateData.SessionTimeout
        }
        if updateData.IdleTimeout >= 0 {
                updates["idle_timeout"] = updateData.IdleTimeout
        }
        if updateData.DailyLimit >= 0 {
                updates["daily_limit"] = updateData.DailyLimit
        }
        if updateData.MonthlyLimit >= 0 {
                updates["monthly_limit"] = updateData.MonthlyLimit
        }
        if updateData.UpRate >= 0 {
                updates["up_rate"] = updateData.UpRate
        }
        if updateData.DownRate >= 0 {
                updates["down_rate"] = updateData.DownRate
        }
        if updateData.UpLimit >= 0 {
                updates["up_limit"] = updateData.UpLimit
        }
        if updateData.DownLimit >= 0 {
                updates["down_limit"] = updateData.DownLimit
        }
        if updateData.TotalLimit >= 0 {
                updates["total_limit"] = updateData.TotalLimit
        }
        if updateData.AddrPool != "" {
                updates["addr_pool"] = updateData.AddrPool
        }
        if updateData.Domain != "" {
                updates["domain"] = updateData.Domain
        }
        if updateData.WelcomeUrl != "" {
                updates["welcome_url"] = updateData.WelcomeUrl
        }
        if updateData.LogoutUrl != "" {
                updates["logout_url"] = updateData.LogoutUrl
        }
        updates["bind_mac"] = updateData.BindMac
        updates["max_devices"] = updateData.MaxDevices
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

// DeleteHotspotProfile deletes a hotspot profile.
// @Summary Delete hotspot profile
// @Tags Hotspot
// @Param id path int true "Profile ID"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/hotspot-profiles/{id} [delete]
func DeleteHotspotProfile(c echo.Context) error {
        id, err := strconv.ParseInt(c.Param("id"), 10, 64)
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_ID", "Invalid profile ID", nil)
        }

        // Check for users using this profile
        var userCount int64
        GetDB(c).Model(&domain.HotspotUser{}).Where("profile_id = ?", id).Count(&userCount)
        if userCount > 0 {
                return fail(c, http.StatusConflict, "IN_USE", "Profile is in use and cannot be deleted", map[string]interface{}{
                        "user_count": userCount,
                })
        }

        if err := GetDB(c).Delete(&domain.HotspotProfile{}, id).Error; err != nil {
                return fail(c, http.StatusInternalServerError, "DELETE_FAILED", "Failed to delete profile", err.Error())
        }

        return ok(c, map[string]interface{}{
                "message": "Deletion successful",
        })
}

// HotspotUserRequest represents the request body for creating a hotspot user.
type HotspotUserRequest struct {
        NodeId     interface{} `json:"node_id"`
        ProfileId  interface{} `json:"profile_id" validate:"required"`
        Username   string      `json:"username" validate:"required,min=1,max=50"`
        Password   string      `json:"password" validate:"omitempty,min=6,max=128"`
        Realname   string      `json:"realname" validate:"omitempty,max=100"`
        Mobile     string      `json:"mobile" validate:"omitempty,max=20"`
        Email      string      `json:"email" validate:"omitempty,email,max=100"`
        MacAddr    string      `json:"mac_addr" validate:"omitempty,mac"`
        IpAddr     string      `json:"ip_addr" validate:"omitempty,ipv4"`
        ExpireTime string      `json:"expire_time"`
        Status     interface{} `json:"status"`
        Remark     string      `json:"remark" validate:"omitempty,max=500"`
}

// toHotspotUser converts HotspotUserRequest to HotspotUser.
func (req *HotspotUserRequest) toHotspotUser() *domain.HotspotUser {
        user := &domain.HotspotUser{
                Username: strings.TrimSpace(req.Username),
                Password: req.Password,
                Realname: req.Realname,
                Mobile:   req.Mobile,
                Email:    req.Email,
                MacAddr:  req.MacAddr,
                IpAddr:   req.IpAddr,
                Remark:   req.Remark,
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

// ListHotspotUsers retrieves the hotspot user list.
// @Summary Get hotspot user list
// @Tags Hotspot
// @Param page query int false "Page number"
// @Param perPage query int false "Items per page"
// @Success 200 {object} ListResponse
// @Router /api/v1/hotspot-users [get]
func ListHotspotUsers(c echo.Context) error {
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
        var users []domain.HotspotUser

        query := db.Model(&domain.HotspotUser{}).
                Select("hotspot_user.*, COALESCE(ro.count, 0) AS online_count").
                Joins("LEFT JOIN (SELECT username, COUNT(1) AS count FROM radius_online GROUP BY username) ro ON hotspot_user.username = ro.username")

        // Filter by username
        if username := strings.TrimSpace(c.QueryParam("username")); username != "" {
                if strings.EqualFold(db.Name(), "postgres") {
                        query = query.Where("hotspot_user.username ILIKE ?", "%"+username+"%")
                } else {
                        query = query.Where("LOWER(hotspot_user.username) LIKE ?", "%"+strings.ToLower(username)+"%")
                }
        }

        // Filter by status
        if status := strings.TrimSpace(c.QueryParam("status")); status != "" {
                query = query.Where("hotspot_user.status = ?", status)
        }

        // Filter by profile_id
        if profileId := c.QueryParam("profile_id"); profileId != "" {
                if id, err := strconv.ParseInt(profileId, 10, 64); err == nil {
                        query = query.Where("hotspot_user.profile_id = ?", id)
                }
        }

        // Filter by mac_addr
        if macAddr := strings.TrimSpace(c.QueryParam("mac_addr")); macAddr != "" {
                query = query.Where("hotspot_user.mac_addr = ?", macAddr)
        }

        query.Count(&total)

        offset := (page - 1) * perPage
        query.Order("hotspot_user." + sortField + " " + order).Limit(perPage).Offset(offset).Find(&users)

        // Clear passwords
        for i := range users {
                users[i].Password = ""
        }

        return paged(c, users, total, page, perPage)
}

// GetHotspotUser retrieves a single hotspot user.
// @Summary Get hotspot user detail
// @Tags Hotspot
// @Param id path int true "User ID"
// @Success 200 {object} domain.HotspotUser
// @Router /api/v1/hotspot-users/{id} [get]
func GetHotspotUser(c echo.Context) error {
        id, err := strconv.ParseInt(c.Param("id"), 10, 64)
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", nil)
        }

        var user domain.HotspotUser
        if err := GetDB(c).First(&user, id).Error; err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return fail(c, http.StatusNotFound, "NOT_FOUND", "Hotspot user not found", nil)
                }
                return fail(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to query user", err.Error())
        }

        user.Password = ""
        return ok(c, user)
}

// CreateHotspotUser creates a new hotspot user.
// @Summary Create hotspot user
// @Tags Hotspot
// @Param user body HotspotUserRequest true "User information"
// @Success 201 {object} domain.HotspotUser
// @Router /api/v1/hotspot-users [post]
func CreateHotspotUser(c echo.Context) error {
        var req HotspotUserRequest
        if err := c.Bind(&req); err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_REQUEST", "Unable to parse request parameters", err.Error())
        }

        if err := c.Validate(&req); err != nil {
                return err
        }

        user := req.toHotspotUser()

        // Validate profile exists
        var profile domain.HotspotProfile
        if err := GetDB(c).Where("id = ?", user.ProfileId).First(&profile).Error; err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return fail(c, http.StatusBadRequest, "PROFILE_NOT_FOUND", "Associated hotspot profile not found", nil)
                }
                return fail(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to query profile", err.Error())
        }

        // Check username uniqueness
        var exists int64
        GetDB(c).Model(&domain.HotspotUser{}).Where("username = ?", user.Username).Count(&exists)
        if exists > 0 {
                return fail(c, http.StatusConflict, "USERNAME_EXISTS", "Username already exists", nil)
        }

        // Set password for MAC auth if not provided
        if user.Password == "" && profile.AuthMode == domain.HotspotAuthModeMAC {
                user.Password = user.Username // MAC address as password
        } else if user.Password == "" {
                return fail(c, http.StatusBadRequest, "MISSING_PASSWORD", "Password is required", nil)
        }

        // Parse expiration time
        expire, err := parseTimeInput(req.ExpireTime, time.Now().AddDate(1, 0, 0))
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_EXPIRE_TIME", "Invalid expire time format", nil)
        }

        // Set defaults
        user.ID = common.UUIDint64()
        user.ExpireTime = expire
        if user.Status == "" {
                user.Status = common.ENABLED
        }
        user.CreatedAt = time.Now()
        user.UpdatedAt = time.Now()

        if err := GetDB(c).Create(user).Error; err != nil {
                return fail(c, http.StatusInternalServerError, "CREATE_FAILED", "Failed to create user", err.Error())
        }

        user.Password = ""
        return ok(c, user)
}

// UpdateHotspotUser updates a hotspot user.
// @Summary Update hotspot user
// @Tags Hotspot
// @Param id path int true "User ID"
// @Param user body HotspotUserRequest true "User information"
// @Success 200 {object} domain.HotspotUser
// @Router /api/v1/hotspot-users/{id} [put]
func UpdateHotspotUser(c echo.Context) error {
        id, err := strconv.ParseInt(c.Param("id"), 10, 64)
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", nil)
        }

        var user domain.HotspotUser
        if err := GetDB(c).First(&user, id).Error; err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return fail(c, http.StatusNotFound, "NOT_FOUND", "Hotspot user not found", nil)
                }
                return fail(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to query user", err.Error())
        }

        var req HotspotUserRequest
        if err := c.Bind(&req); err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_REQUEST", "Unable to parse request parameters", err.Error())
        }

        if err := c.Validate(&req); err != nil {
                return err
        }

        updateData := req.toHotspotUser()

        // Validate username uniqueness
        if updateData.Username != "" && updateData.Username != user.Username {
                var count int64
                GetDB(c).Model(&domain.HotspotUser{}).Where("username = ? AND id != ?", updateData.Username, id).Count(&count)
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
        if updateData.MacAddr != "" {
                updates["mac_addr"] = updateData.MacAddr
        }
        if updateData.IpAddr != "" {
                updates["ip_addr"] = updateData.IpAddr
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

// DeleteHotspotUser deletes a hotspot user.
// @Summary Delete hotspot user
// @Tags Hotspot
// @Param id path int true "User ID"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/hotspot-users/{id} [delete]
func DeleteHotspotUser(c echo.Context) error {
        id, err := strconv.ParseInt(c.Param("id"), 10, 64)
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", nil)
        }

        if err := GetDB(c).Delete(&domain.HotspotUser{}, id).Error; err != nil {
                return fail(c, http.StatusInternalServerError, "DELETE_FAILED", "Failed to delete user", err.Error())
        }

        return ok(c, map[string]interface{}{
                "id": id,
        })
}

// registerHotspotRoutes registers hotspot-related routes.
func registerHotspotRoutes() {
        webserver.ApiGET("/hotspot-profiles", ListHotspotProfiles)
        webserver.ApiGET("/hotspot-profiles/:id", GetHotspotProfile)
        webserver.ApiPOST("/hotspot-profiles", CreateHotspotProfile)
        webserver.ApiPUT("/hotspot-profiles/:id", UpdateHotspotProfile)
        webserver.ApiDELETE("/hotspot-profiles/:id", DeleteHotspotProfile)
        webserver.ApiGET("/hotspot-users", ListHotspotUsers)
        webserver.ApiGET("/hotspot-users/:id", GetHotspotUser)
        webserver.ApiPOST("/hotspot-users", CreateHotspotUser)
        webserver.ApiPUT("/hotspot-users/:id", UpdateHotspotUser)
        webserver.ApiDELETE("/hotspot-users/:id", DeleteHotspotUser)
}
