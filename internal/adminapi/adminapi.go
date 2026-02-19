package adminapi

import (
        "github.com/talkincode/toughradius/v9/internal/app"
)

// Init registers all admin API routes.
// This function is called during application startup to register all
// administrative API endpoints for the ToughRADIUS web interface.
//
// Registered modules:
//   - Authentication: Login, logout, session management
//   - Users: RADIUS user CRUD operations
//   - Profiles: RADIUS billing profile management
//   - Accounting: RADIUS accounting records
//   - Sessions: Online session management
//   - NAS: Network Access Server management
//   - Settings: System configuration
//   - Nodes: Network node management
//   - Operators: Admin operator management
//   - Vouchers: Prepaid voucher management
//   - Hotspot: Hotspot profile and user management
//   - PPPoE: PPPoE profile and user management
func Init(appCtx app.AppContext) {
        registerAuthRoutes()
        registerUserRoutes()
        registerDashboardRoutes()
        registerProfileRoutes()
        registerAccountingRoutes()
        registerSessionRoutes()
        registerNASRoutes()
        registerSettingsRoutes()
        registerNodesRoutes()
        registerOperatorsRoutes()
        registerVoucherRoutes()
        registerHotspotRoutes()
        registerPppoeRoutes()
}
