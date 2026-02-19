package domain

var Tables = []interface{}{
        // System
        &SysConfig{},
        &SysOpr{},
        &SysOprLog{},
        // Network
        &NetNode{},
        &NetNas{},
        // Radius
        &RadiusAccounting{},
        &RadiusOnline{},
        &RadiusProfile{},
        &RadiusUser{},
        // Voucher
        &VoucherBatch{},
        &Voucher{},
        // Hotspot
        &HotspotProfile{},
        &HotspotUser{},
        // PPPoE
        &PppoeProfile{},
        &PppoeUser{},
}
