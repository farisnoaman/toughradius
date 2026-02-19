// internal/domain/voucher_marshal.go
package domain

import (
        "encoding/json"
        "time"

        "github.com/araddon/dateparse"
        "github.com/talkincode/toughradius/v9/pkg/timeutil"
)

// MarshalJSON implements custom JSON marshaling for VoucherBatch.
func (d VoucherBatch) MarshalJSON() ([]byte, error) {
        type Alias VoucherBatch
        return json.Marshal(&struct {
                Alias
                ExpireTime string `json:"expire_time"`
        }{
                Alias:      (Alias)(d),
                ExpireTime: d.ExpireTime.Format(timeutil.YYYYMMDDHHMMSS_LAYOUT),
        })
}

// UnmarshalJSON implements custom JSON unmarshaling for VoucherBatch.
func (d *VoucherBatch) UnmarshalJSON(b []byte) error {
        type Alias VoucherBatch
        var tmp = struct {
                *Alias
                ExpireTime string `json:"expire_time"`
        }{
                Alias: (*Alias)(d),
        }
        if err := json.Unmarshal(b, &tmp); err != nil {
                return err
        }
        if t, err := time.ParseInLocation(timeutil.YYYYMMDDHHMMSS_LAYOUT, tmp.ExpireTime, time.Local); err == nil {
                d.ExpireTime = t
        } else if tmp.ExpireTime != "" {
                d.ExpireTime, _ = dateparse.ParseAny(tmp.ExpireTime)
        }
        return nil
}

// MarshalJSON implements custom JSON marshaling for Voucher.
func (d Voucher) MarshalJSON() ([]byte, error) {
        type Alias Voucher
        return json.Marshal(&struct {
                Alias
                RedeemedAt string `json:"redeemed_at"`
                ExpireTime string `json:"expire_time"`
        }{
                Alias:      (Alias)(d),
                RedeemedAt: formatTimePtr(d.RedeemedAt),
                ExpireTime: formatTimePtr(d.ExpireTime),
        })
}

// UnmarshalJSON implements custom JSON unmarshaling for Voucher.
func (d *Voucher) UnmarshalJSON(b []byte) error {
        type Alias Voucher
        var tmp = struct {
                *Alias
                RedeemedAt string `json:"redeemed_at"`
                ExpireTime string `json:"expire_time"`
        }{
                Alias: (*Alias)(d),
        }
        if err := json.Unmarshal(b, &tmp); err != nil {
                return err
        }
        if tmp.RedeemedAt != "" {
                t, _ := dateparse.ParseAny(tmp.RedeemedAt)
                d.RedeemedAt = &t
        }
        if tmp.ExpireTime != "" {
                t, _ := dateparse.ParseAny(tmp.ExpireTime)
                d.ExpireTime = &t
        }
        return nil
}

// MarshalJSON implements custom JSON marshaling for HotspotUser.
func (d HotspotUser) MarshalJSON() ([]byte, error) {
        type Alias HotspotUser
        return json.Marshal(&struct {
                Alias
                ExpireTime string `json:"expire_time"`
                LastOnline string `json:"last_online"`
        }{
                Alias:      (Alias)(d),
                ExpireTime: d.ExpireTime.Format(timeutil.YYYYMMDDHHMMSS_LAYOUT),
                LastOnline: d.LastOnline.Format(timeutil.YYYYMMDDHHMM_LAYOUT),
        })
}

// UnmarshalJSON implements custom JSON unmarshaling for HotspotUser.
func (d *HotspotUser) UnmarshalJSON(b []byte) error {
        type Alias HotspotUser
        var tmp = struct {
                *Alias
                ExpireTime string `json:"expire_time"`
                LastOnline string `json:"last_online"`
        }{
                Alias: (*Alias)(d),
        }
        if err := json.Unmarshal(b, &tmp); err != nil {
                return err
        }
        if t, err := time.ParseInLocation(timeutil.YYYYMMDDHHMMSS_LAYOUT, tmp.ExpireTime, time.Local); err == nil {
                d.ExpireTime = t
        } else if tmp.ExpireTime != "" {
                d.ExpireTime, _ = dateparse.ParseAny(tmp.ExpireTime)
        }
        d.LastOnline, _ = dateparse.ParseAny(tmp.LastOnline)
        return nil
}

// MarshalJSON implements custom JSON marshaling for PppoeUser.
func (d PppoeUser) MarshalJSON() ([]byte, error) {
        type Alias PppoeUser
        return json.Marshal(&struct {
                Alias
                ExpireTime string `json:"expire_time"`
                LastOnline string `json:"last_online"`
        }{
                Alias:      (Alias)(d),
                ExpireTime: d.ExpireTime.Format(timeutil.YYYYMMDDHHMMSS_LAYOUT),
                LastOnline: d.LastOnline.Format(timeutil.YYYYMMDDHHMM_LAYOUT),
        })
}

// UnmarshalJSON implements custom JSON unmarshaling for PppoeUser.
func (d *PppoeUser) UnmarshalJSON(b []byte) error {
        type Alias PppoeUser
        var tmp = struct {
                *Alias
                ExpireTime string `json:"expire_time"`
                LastOnline string `json:"last_online"`
        }{
                Alias: (*Alias)(d),
        }
        if err := json.Unmarshal(b, &tmp); err != nil {
                return err
        }
        if t, err := time.ParseInLocation(timeutil.YYYYMMDDHHMMSS_LAYOUT, tmp.ExpireTime, time.Local); err == nil {
                d.ExpireTime = t
        } else if tmp.ExpireTime != "" {
                d.ExpireTime, _ = dateparse.ParseAny(tmp.ExpireTime)
        }
        d.LastOnline, _ = dateparse.ParseAny(tmp.LastOnline)
        return nil
}

// formatTimePtr formats a time pointer to string, returns empty string if nil.
func formatTimePtr(t *time.Time) string {
        if t == nil {
                return ""
        }
        return t.Format(timeutil.YYYYMMDDHHMMSS_LAYOUT)
}
