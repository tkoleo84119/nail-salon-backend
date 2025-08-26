## User Story

作為一位員工，我希望能新增顧客優惠券，方便維護顧客優惠券資訊。

---

## Endpoint

**POST** `/api/admin/customer_coupons`

---

## 說明

- 提供員工新增顧客優惠券功能。
- 可設定有效期間。

---

## 權限

- 需要登入才可使用。
- 所有角色皆可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Body 範例

```json
{
  "customerId": "1234567890",
  "couponId": "1234567890",
  "validFrom": "2021-01-01T00:00:00+08:00",
  "validTo": "2021-01-01T00:00:00+08:00"
}
```

### 驗證規則

| 欄位       | 必填 | 其他規則 | 說明                             |
| ---------- | ---- | -------- | -------------------------------- |
| customerId | 是   |          | 顧客ID                           |
| couponId   | 是   |          | 優惠券ID                         |
| validFrom  | 是   |          | 有效期間開始                     |
| validTo    | 否   |          | 有效期間結束，若不填則為永久有效 |

- `validFrom` 與 `validTo` 要傳入標準 Iso 8601 格式。

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "id": "9000000001"
  }
}
```

### 錯誤處理

全部 API 皆回傳如下結構，請參考錯誤總覽。

```json
{
  "errors": [
    {
      "code": "EXXXX",
      "message": "錯誤訊息",
      "field": "錯誤欄位名稱"
    }
  ]
}
```

- 欄位說明：
  - errors: 錯誤陣列（支援多筆同時回報）
  - code: 錯誤代碼，唯一對應每種錯誤
  - message: 中文錯誤訊息（可參照錯誤總覽）
  - field: 參數欄位名稱（僅部分驗證錯誤有）

| 狀態碼 | 錯誤碼    | 常數名稱                            | 說明                                         |
| ------ | --------- | ----------------------------------- | -------------------------------------------- |
| 401    | E1002     | AuthTokenInvalid                    | 無效的 accessToken，請重新登入               |
| 401    | E1003     | AuthTokenMissing                    | accessToken 缺失，請重新登入                 |
| 401    | E1004     | AuthTokenFormatError                | accessToken 格式錯誤，請重新登入             |
| 401    | E1005     | AuthStaffFailed                     | 未找到有效的員工資訊，請重新登入             |
| 401    | E1006     | AuthContextMissing                  | 未找到使用者認證資訊，請重新登入             |
| 403    | E1010     | AuthPermissionDenied                | 權限不足，無法執行此操作                     |
| 400    | E2001     | ValJsonFormat                       | JSON 格式錯誤，請檢查                        |
| 400    | E2004     | ValTypeConversionFailed             | 參數類型轉換失敗                             |
| 400    | E2020     | ValFieldRequired                    | {field} 為必填項目                           |
| 400    | E2037     | ValFieldISO8601Format               | {field} 格式錯誤，請使用正確的 ISO 8601 格式 |
| 400    | E3CCOU001 | CustomerCouponValidFromBeforeNow    | 有效期間開始不能小於當前時間                 |
| 400    | E3CCOU002 | CustomerCouponValidFromAfterValidTo | 有效期間開始不能大於有效期間結束             |
| 400    | E3CCOU003 | CustomerCouponValidToBeforeNow      | 有效期間結束不能小於當前時間                 |
| 404    | E3C001    | CustomerNotFound                    | 客戶不存在                                   |
| 404    | E3COU004  | CouponNotFound                      | 優惠券不存在或已被刪除                       |
| 500    | E9001     | SysInternalError                    | 系統發生錯誤，請稍後再試                     |
| 500    | E9002     | SysDatabaseError                    | 資料庫操作失敗                               |

---

## 資料表

- `customer_coupons`
- `customers`
- `coupons`

---

## Service 邏輯

1. 確認 `customer_id` 與 `coupon_id` 是否存在。
2. 確認 `valid_from` 不能小於當前時間。
3. 若 `valid_to` 不為空，則確認以下邏輯：
   - `valid_from` 不能大於 `valid_to`。
   - `valid_to` 不能小於當前時間。
4. 建立 `customer_coupons` 資料。
5. 回傳新增結果。

---

## 注意事項

- `valid_from` 與 `valid_to` 會是標準 Iso 8601 格式。
