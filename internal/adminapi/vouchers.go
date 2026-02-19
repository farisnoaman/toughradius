// internal/adminapi/vouchers.go
package adminapi

import (
        "crypto/rand"
        "encoding/hex"
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

// VoucherBatchRequest represents the request body for creating a voucher batch.
type VoucherBatchRequest struct {
        Name       string      `json:"name" validate:"required,min=1,max=100"`
        NodeId     interface{} `json:"node_id"`
        ProfileId  interface{} `json:"profile_id" validate:"required"`
        TotalCount int         `json:"total_count" validate:"required,gte=1,lte=10000"`
        ExpireTime string      `json:"expire_time" validate:"required"`
        ValidDays  int         `json:"valid_days" validate:"gte=0,lte=3650"`
        Prefix     string      `json:"prefix" validate:"omitempty,max=10"`
        CodeLength int         `json:"code_length" validate:"gte=6,lte=32"`
        Status     interface{} `json:"status"`
        Remark     string      `json:"remark" validate:"omitempty,max=500"`
}

// toVoucherBatch converts VoucherBatchRequest to VoucherBatch.
func (req *VoucherBatchRequest) toVoucherBatch() *domain.VoucherBatch {
        batch := &domain.VoucherBatch{
                Name:       strings.TrimSpace(req.Name),
                TotalCount: req.TotalCount,
                ValidDays:  req.ValidDays,
                Prefix:     req.Prefix,
                CodeLength: req.CodeLength,
                Remark:     req.Remark,
        }

        // Handle profile_id
        switch v := req.ProfileId.(type) {
        case float64:
                batch.ProfileId = int64(v)
        case string:
                if v != "" {
                        id, _ := strconv.ParseInt(v, 10, 64)
                        batch.ProfileId = id
                }
        }

        // Handle node_id
        switch v := req.NodeId.(type) {
        case float64:
                batch.NodeId = int64(v)
        case string:
                if v != "" {
                        id, _ := strconv.ParseInt(v, 10, 64)
                        batch.NodeId = id
                }
        }

        // Handle status
        switch v := req.Status.(type) {
        case bool:
                if v {
                        batch.Status = domain.VoucherBatchStatusEnabled
                } else {
                        batch.Status = domain.VoucherBatchStatusDisabled
                }
        case string:
                batch.Status = strings.ToLower(v)
        }

        // Set default code length
        if batch.CodeLength < 6 {
                batch.CodeLength = 10
        }

        return batch
}

// VoucherBatchUpdateRequest represents the request body for updating a voucher batch.
type VoucherBatchUpdateRequest struct {
        Name       string      `json:"name" validate:"omitempty,min=1,max=100"`
        NodeId     interface{} `json:"node_id"`
        ProfileId  interface{} `json:"profile_id"`
        ExpireTime string      `json:"expire_time"`
        ValidDays  int         `json:"valid_days" validate:"gte=0,lte=3650"`
        Status     interface{} `json:"status"`
        Remark     string      `json:"remark" validate:"omitempty,max=500"`
}

// toVoucherBatch converts VoucherBatchUpdateRequest to VoucherBatch.
func (req *VoucherBatchUpdateRequest) toVoucherBatch() *domain.VoucherBatch {
        batch := &domain.VoucherBatch{
                Name:      strings.TrimSpace(req.Name),
                ValidDays: req.ValidDays,
                Remark:    req.Remark,
        }

        // Handle profile_id
        switch v := req.ProfileId.(type) {
        case float64:
                batch.ProfileId = int64(v)
        case string:
                if v != "" {
                        id, _ := strconv.ParseInt(v, 10, 64)
                        batch.ProfileId = id
                }
        }

        // Handle node_id
        switch v := req.NodeId.(type) {
        case float64:
                batch.NodeId = int64(v)
        case string:
                if v != "" {
                        id, _ := strconv.ParseInt(v, 10, 64)
                        batch.NodeId = id
                }
        }

        // Handle status
        switch v := req.Status.(type) {
        case bool:
                if v {
                        batch.Status = domain.VoucherBatchStatusEnabled
                } else {
                        batch.Status = domain.VoucherBatchStatusDisabled
                }
        case string:
                batch.Status = strings.ToLower(v)
        }

        return batch
}

// ListVoucherBatches retrieves the voucher batch list.
// @Summary Get voucher batch list
// @Tags Voucher
// @Param page query int false "Page number"
// @Param perPage query int false "Items per page"
// @Param sort query string false "Sort field"
// @Param order query string false "Sort direction"
// @Success 200 {object} ListResponse
// @Router /api/v1/voucher-batches [get]
func ListVoucherBatches(c echo.Context) error {
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
        var batches []domain.VoucherBatch

        query := db.Model(&domain.VoucherBatch{})

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

        query.Count(&total)

        offset := (page - 1) * perPage
        query.Order(sortField + " " + order).Limit(perPage).Offset(offset).Find(&batches)

        return paged(c, batches, total, page, perPage)
}

// GetVoucherBatch retrieves a single voucher batch.
// @Summary Get voucher batch detail
// @Tags Voucher
// @Param id path int true "Batch ID"
// @Success 200 {object} domain.VoucherBatch
// @Router /api/v1/voucher-batches/{id} [get]
func GetVoucherBatch(c echo.Context) error {
        id, err := strconv.ParseInt(c.Param("id"), 10, 64)
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_ID", "Invalid batch ID", nil)
        }

        var batch domain.VoucherBatch
        if err := GetDB(c).First(&batch, id).Error; err != nil {
                return fail(c, http.StatusNotFound, "NOT_FOUND", "Voucher batch not found", nil)
        }

        return ok(c, batch)
}

// CreateVoucherBatch creates a new voucher batch and generates vouchers.
// @Summary Create voucher batch
// @Tags Voucher
// @Param batch body VoucherBatchRequest true "Batch information"
// @Success 201 {object} domain.VoucherBatch
// @Router /api/v1/voucher-batches [post]
func CreateVoucherBatch(c echo.Context) error {
        var req VoucherBatchRequest
        if err := c.Bind(&req); err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_REQUEST", "Unable to parse request parameters", err.Error())
        }

        if err := c.Validate(&req); err != nil {
                return err
        }

        batch := req.toVoucherBatch()

        // Validate profile exists
        var profile domain.RadiusProfile
        if err := GetDB(c).Where("id = ?", batch.ProfileId).First(&profile).Error; err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return fail(c, http.StatusBadRequest, "PROFILE_NOT_FOUND", "Associated billing profile not found", nil)
                }
                return fail(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to query profile", err.Error())
        }

        // Check name uniqueness
        var count int64
        GetDB(c).Model(&domain.VoucherBatch{}).Where("name = ?", batch.Name).Count(&count)
        if count > 0 {
                return fail(c, http.StatusConflict, "NAME_EXISTS", "Batch name already exists", nil)
        }

        // Parse expiration time
        expire, err := parseTimeInput(req.ExpireTime, time.Now().AddDate(1, 0, 0))
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_EXPIRE_TIME", "Invalid expire time format", nil)
        }
        batch.ExpireTime = expire

        // Set defaults
        if batch.Status == "" {
                batch.Status = domain.VoucherBatchStatusEnabled
        }
        batch.CreatedAt = time.Now()
        batch.UpdatedAt = time.Now()

        // Start transaction to create batch and vouchers
        tx := GetDB(c).Begin()
        defer func() {
                if r := recover(); r != nil {
                        tx.Rollback()
                }
        }()

        if err := tx.Create(batch).Error; err != nil {
                tx.Rollback()
                return fail(c, http.StatusInternalServerError, "CREATE_FAILED", "Failed to create voucher batch", err.Error())
        }

        // Generate vouchers
        vouchers := make([]domain.Voucher, 0, batch.TotalCount)
        for i := 0; i < batch.TotalCount; i++ {
                code, err := generateVoucherCode(batch.Prefix, batch.CodeLength)
                if err != nil {
                        tx.Rollback()
                        return fail(c, http.StatusInternalServerError, "GENERATE_FAILED", "Failed to generate voucher codes", err.Error())
                }
                vouchers = append(vouchers, domain.Voucher{
                        BatchId:    batch.ID,
                        Code:       code,
                        ProfileId:  batch.ProfileId,
                        Status:     domain.VoucherStatusAvailable,
                        ExpireTime: &batch.ExpireTime,
                        CreatedAt:  time.Now(),
                        UpdatedAt:  time.Now(),
                })
        }

        // Batch insert vouchers
        if err := tx.CreateInBatches(vouchers, 100).Error; err != nil {
                tx.Rollback()
                return fail(c, http.StatusInternalServerError, "CREATE_FAILED", "Failed to create vouchers", err.Error())
        }

        if err := tx.Commit().Error; err != nil {
                return fail(c, http.StatusInternalServerError, "COMMIT_FAILED", "Failed to commit transaction", err.Error())
        }

        return ok(c, batch)
}

// UpdateVoucherBatch updates a voucher batch.
// @Summary Update voucher batch
// @Tags Voucher
// @Param id path int true "Batch ID"
// @Param batch body VoucherBatchUpdateRequest true "Batch information"
// @Success 200 {object} domain.VoucherBatch
// @Router /api/v1/voucher-batches/{id} [put]
func UpdateVoucherBatch(c echo.Context) error {
        id, err := strconv.ParseInt(c.Param("id"), 10, 64)
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_ID", "Invalid batch ID", nil)
        }

        var batch domain.VoucherBatch
        if err := GetDB(c).First(&batch, id).Error; err != nil {
                return fail(c, http.StatusNotFound, "NOT_FOUND", "Voucher batch not found", nil)
        }

        var req VoucherBatchUpdateRequest
        if err := c.Bind(&req); err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_REQUEST", "Unable to parse request parameters", err.Error())
        }

        if err := c.Validate(&req); err != nil {
                return err
        }

        updateData := req.toVoucherBatch()

        // Validate name uniqueness
        if updateData.Name != "" && updateData.Name != batch.Name {
                var count int64
                GetDB(c).Model(&domain.VoucherBatch{}).Where("name = ? AND id != ?", updateData.Name, id).Count(&count)
                if count > 0 {
                        return fail(c, http.StatusConflict, "NAME_EXISTS", "Batch name already exists", nil)
                }
        }

        // Validate profile exists
        if updateData.ProfileId > 0 {
                var profile domain.RadiusProfile
                if err := GetDB(c).Where("id = ?", updateData.ProfileId).First(&profile).Error; err != nil {
                        if errors.Is(err, gorm.ErrRecordNotFound) {
                                return fail(c, http.StatusBadRequest, "PROFILE_NOT_FOUND", "Associated billing profile not found", nil)
                        }
                }
        }

        // Build updates map
        updates := map[string]interface{}{}
        if updateData.Name != "" {
                updates["name"] = updateData.Name
        }
        if updateData.NodeId > 0 {
                updates["node_id"] = updateData.NodeId
        }
        if updateData.ProfileId > 0 {
                updates["profile_id"] = updateData.ProfileId
        }
        if updateData.Status != "" {
                updates["status"] = updateData.Status
        }
        if updateData.ValidDays > 0 {
                updates["valid_days"] = updateData.ValidDays
        }
        if updateData.Remark != "" {
                updates["remark"] = updateData.Remark
        }
        if req.ExpireTime != "" {
                expire, err := parseTimeInput(req.ExpireTime, batch.ExpireTime)
                if err != nil {
                        return fail(c, http.StatusBadRequest, "INVALID_EXPIRE_TIME", "Invalid expire time format", nil)
                }
                updates["expire_time"] = expire
        }
        updates["updated_at"] = time.Now()

        if err := GetDB(c).Model(&batch).Updates(updates).Error; err != nil {
                return fail(c, http.StatusInternalServerError, "UPDATE_FAILED", "Failed to update batch", err.Error())
        }

        GetDB(c).First(&batch, id)
        return ok(c, batch)
}

// DeleteVoucherBatch deletes a voucher batch.
// @Summary Delete voucher batch
// @Tags Voucher
// @Param id path int true "Batch ID"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/voucher-batches/{id} [delete]
func DeleteVoucherBatch(c echo.Context) error {
        id, err := strconv.ParseInt(c.Param("id"), 10, 64)
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_ID", "Invalid batch ID", nil)
        }

        // Check for used vouchers
        var usedCount int64
        GetDB(c).Model(&domain.Voucher{}).Where("batch_id = ? AND status = ?", id, domain.VoucherStatusUsed).Count(&usedCount)
        if usedCount > 0 {
                return fail(c, http.StatusConflict, "IN_USE", "Cannot delete batch with used vouchers", map[string]interface{}{
                        "used_count": usedCount,
                })
        }

        // Delete vouchers and batch in transaction
        tx := GetDB(c).Begin()
        if err := tx.Where("batch_id = ?", id).Delete(&domain.Voucher{}).Error; err != nil {
                tx.Rollback()
                return fail(c, http.StatusInternalServerError, "DELETE_FAILED", "Failed to delete vouchers", err.Error())
        }
        if err := tx.Delete(&domain.VoucherBatch{}, id).Error; err != nil {
                tx.Rollback()
                return fail(c, http.StatusInternalServerError, "DELETE_FAILED", "Failed to delete batch", err.Error())
        }
        if err := tx.Commit().Error; err != nil {
                return fail(c, http.StatusInternalServerError, "COMMIT_FAILED", "Failed to commit transaction", err.Error())
        }

        return ok(c, map[string]interface{}{
                "message": "Deletion successful",
        })
}

// ListVouchers retrieves the voucher list for a batch.
// @Summary Get voucher list
// @Tags Voucher
// @Param batch_id query int false "Batch ID"
// @Param status query string false "Voucher status"
// @Param page query int false "Page number"
// @Param perPage query int false "Items per page"
// @Success 200 {object} ListResponse
// @Router /api/v1/vouchers [get]
func ListVouchers(c echo.Context) error {
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
        var vouchers []domain.Voucher

        query := db.Model(&domain.Voucher{})

        // Filter by batch_id
        if batchId := c.QueryParam("batch_id"); batchId != "" {
                if id, err := strconv.ParseInt(batchId, 10, 64); err == nil {
                        query = query.Where("batch_id = ?", id)
                }
        }

        // Filter by status
        if status := strings.TrimSpace(c.QueryParam("status")); status != "" {
                query = query.Where("status = ?", status)
        }

        // Filter by code
        if code := strings.TrimSpace(c.QueryParam("code")); code != "" {
                query = query.Where("code LIKE ?", "%"+code+"%")
        }

        query.Count(&total)

        offset := (page - 1) * perPage
        query.Order(sortField + " " + order).Limit(perPage).Offset(offset).Find(&vouchers)

        // Clear passwords in response
        for i := range vouchers {
                vouchers[i].Password = ""
        }

        return paged(c, vouchers, total, page, perPage)
}

// GetVoucher retrieves a single voucher.
// @Summary Get voucher detail
// @Tags Voucher
// @Param id path int true "Voucher ID"
// @Success 200 {object} domain.Voucher
// @Router /api/v1/vouchers/{id} [get]
func GetVoucher(c echo.Context) error {
        id, err := strconv.ParseInt(c.Param("id"), 10, 64)
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_ID", "Invalid voucher ID", nil)
        }

        var voucher domain.Voucher
        if err := GetDB(c).First(&voucher, id).Error; err != nil {
                return fail(c, http.StatusNotFound, "NOT_FOUND", "Voucher not found", nil)
        }

        voucher.Password = ""
        return ok(c, voucher)
}

// RedeemVoucherRequest represents the request body for redeeming a voucher.
type RedeemVoucherRequest struct {
        Code     string `json:"code" validate:"required"`
        Password string `json:"password"`
        Username string `json:"username" validate:"required,min=3,max=50"`
        Password2 string `json:"password2" validate:"required,min=6,max=128"`
        Realname string `json:"realname"`
        Mobile   string `json:"mobile"`
        Email    string `json:"email"`
}

// RedeemVoucher redeems a voucher to create a user account.
// @Summary Redeem voucher
// @Tags Voucher
// @Param request body RedeemVoucherRequest true "Redemption information"
// @Success 200 {object} domain.RadiusUser
// @Router /api/v1/vouchers/redeem [post]
func RedeemVoucher(c echo.Context) error {
        var req RedeemVoucherRequest
        if err := c.Bind(&req); err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_REQUEST", "Unable to parse request parameters", err.Error())
        }

        if err := c.Validate(&req); err != nil {
                return err
        }

        // Find voucher by code
        var voucher domain.Voucher
        if err := GetDB(c).Where("code = ?", req.Code).First(&voucher).Error; err != nil {
                if errors.Is(err, gorm.ErrRecordNotFound) {
                        return fail(c, http.StatusNotFound, "VOUCHER_NOT_FOUND", "Invalid voucher code", nil)
                }
                return fail(c, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to query voucher", err.Error())
        }

        // Validate voucher status
        if voucher.Status != domain.VoucherStatusAvailable {
                return fail(c, http.StatusBadRequest, "VOUCHER_NOT_AVAILABLE", "Voucher is not available for redemption", map[string]interface{}{
                        "status": voucher.Status,
                })
        }

        // Check if voucher has password and validate
        if voucher.Password != "" && voucher.Password != req.Password {
                return fail(c, http.StatusUnauthorized, "INVALID_PASSWORD", "Invalid voucher password", nil)
        }

        // Check expiration
        if voucher.ExpireTime != nil && voucher.ExpireTime.Before(time.Now()) {
                // Update voucher status to expired
                GetDB(c).Model(&voucher).Update("status", domain.VoucherStatusExpired)
                return fail(c, http.StatusBadRequest, "VOUCHER_EXPIRED", "Voucher has expired", nil)
        }

        // Check batch status
        var batch domain.VoucherBatch
        if err := GetDB(c).First(&batch, voucher.BatchId).Error; err == nil {
                if batch.Status != domain.VoucherBatchStatusEnabled {
                        return fail(c, http.StatusBadRequest, "BATCH_DISABLED", "Voucher batch is disabled", nil)
                }
        }

        // Check if username already exists
        var exists int64
        GetDB(c).Model(&domain.RadiusUser{}).Where("username = ?", req.Username).Count(&exists)
        if exists > 0 {
                return fail(c, http.StatusConflict, "USERNAME_EXISTS", "Username already exists", nil)
        }

        // Get profile
        var profile domain.RadiusProfile
        if err := GetDB(c).First(&profile, voucher.ProfileId).Error; err != nil {
                return fail(c, http.StatusBadRequest, "PROFILE_NOT_FOUND", "Associated billing profile not found", nil)
        }

        // Calculate expiration
        expireTime := time.Now().AddDate(0, 0, batch.ValidDays)
        if batch.ValidDays == 0 && voucher.ExpireTime != nil {
                expireTime = *voucher.ExpireTime
        }

        // Start transaction
        tx := GetDB(c).Begin()
        defer func() {
                if r := recover(); r != nil {
                        tx.Rollback()
                }
        }()

        // Create user
        user := &domain.RadiusUser{
                ID:         common.UUIDint64(),
                ProfileId:  voucher.ProfileId,
                Username:   strings.TrimSpace(req.Username),
                Password:   req.Password2,
                Realname:   req.Realname,
                Mobile:     req.Mobile,
                Email:      req.Email,
                AddrPool:   profile.AddrPool,
                ActiveNum:  profile.ActiveNum,
                UpRate:     profile.UpRate,
                DownRate:   profile.DownRate,
                Domain:     profile.Domain,
                BindMac:    profile.BindMac,
                BindVlan:   profile.BindVlan,
                Status:     common.ENABLED,
                ExpireTime: expireTime,
                CreatedAt:  time.Now(),
                UpdatedAt:  time.Now(),
        }

        if err := tx.Create(user).Error; err != nil {
                tx.Rollback()
                return fail(c, http.StatusInternalServerError, "CREATE_USER_FAILED", "Failed to create user", err.Error())
        }

        // Update voucher
        now := time.Now()
        if err := tx.Model(&voucher).Updates(map[string]interface{}{
                "status":      domain.VoucherStatusUsed,
                "user_id":     user.ID,
                "redeemed_at": now,
                "updated_at":  now,
        }).Error; err != nil {
                tx.Rollback()
                return fail(c, http.StatusInternalServerError, "UPDATE_VOUCHER_FAILED", "Failed to update voucher", err.Error())
        }

        // Update batch used count
        if err := tx.Model(&batch).UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error; err != nil {
                tx.Rollback()
                return fail(c, http.StatusInternalServerError, "UPDATE_BATCH_FAILED", "Failed to update batch", err.Error())
        }

        if err := tx.Commit().Error; err != nil {
                return fail(c, http.StatusInternalServerError, "COMMIT_FAILED", "Failed to commit transaction", err.Error())
        }

        user.Password = ""
        return ok(c, user)
}

// DisableVoucher disables a voucher.
// @Summary Disable voucher
// @Tags Voucher
// @Param id path int true "Voucher ID"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/vouchers/{id}/disable [post]
func DisableVoucher(c echo.Context) error {
        id, err := strconv.ParseInt(c.Param("id"), 10, 64)
        if err != nil {
                return fail(c, http.StatusBadRequest, "INVALID_ID", "Invalid voucher ID", nil)
        }

        var voucher domain.Voucher
        if err := GetDB(c).First(&voucher, id).Error; err != nil {
                return fail(c, http.StatusNotFound, "NOT_FOUND", "Voucher not found", nil)
        }

        if voucher.Status == domain.VoucherStatusUsed {
                return fail(c, http.StatusBadRequest, "VOUCHER_USED", "Cannot disable a used voucher", nil)
        }

        if err := GetDB(c).Model(&voucher).Updates(map[string]interface{}{
                "status":     domain.VoucherStatusDisabled,
                "updated_at": time.Now(),
        }).Error; err != nil {
                return fail(c, http.StatusInternalServerError, "UPDATE_FAILED", "Failed to disable voucher", err.Error())
        }

        return ok(c, map[string]interface{}{
                "message": "Voucher disabled successfully",
        })
}

// generateVoucherCode generates a unique voucher code with optional prefix.
func generateVoucherCode(prefix string, length int) (string, error) {
        bytes := make([]byte, length/2+1)
        if _, err := rand.Read(bytes); err != nil {
                return "", err
        }
        code := hex.EncodeToString(bytes)
        if len(code) > length {
                code = code[:length]
        }
        if prefix != "" {
                code = prefix + code
        }
        return strings.ToUpper(code), nil
}

// registerVoucherRoutes registers voucher-related routes.
func registerVoucherRoutes() {
        webserver.ApiGET("/voucher-batches", ListVoucherBatches)
        webserver.ApiGET("/voucher-batches/:id", GetVoucherBatch)
        webserver.ApiPOST("/voucher-batches", CreateVoucherBatch)
        webserver.ApiPUT("/voucher-batches/:id", UpdateVoucherBatch)
        webserver.ApiDELETE("/voucher-batches/:id", DeleteVoucherBatch)
        webserver.ApiGET("/vouchers", ListVouchers)
        webserver.ApiGET("/vouchers/:id", GetVoucher)
        webserver.ApiPOST("/vouchers/redeem", RedeemVoucher)
        webserver.ApiPOST("/vouchers/:id/disable", DisableVoucher)
}
