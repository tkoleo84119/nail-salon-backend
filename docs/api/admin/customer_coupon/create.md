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
  "period": "unlimited"
}
```

### 驗證規則

| 欄位       | 必填 | 其他規則                                                      | 說明     |
| ---------- | ---- | ------------------------------------------------------------- | -------- |
| customerId | 是   |                                                               | 顧客ID   |
| couponId   | 是   |                                                               | 優惠券ID |
| period     | 是   | <li>值可以為 `unlimited` `1month` `3months` `6months` `1year` | 有效期間 |

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

| 狀態碼 | 錯誤碼    | 常數名稱                    | 說明                              |
| ------ | --------- | --------------------------- | --------------------------------- |
| 401    | E1002     | AuthTokenInvalid            | 無效的 accessToken，請重新登入    |
| 401    | E1003     | AuthTokenMissing            | accessToken 缺失，請重新登入      |
| 401    | E1004     | AuthTokenFormatError        | accessToken 格式錯誤，請重新登入  |
| 401    | E1005     | AuthStaffFailed             | 未找到有效的員工資訊，請重新登入  |
| 401    | E1006     | AuthContextMissing          | 未找到使用者認證資訊，請重新登入  |
| 403    | E1010     | AuthPermissionDenied        | 權限不足，無法執行此操作          |
| 400    | E2001     | ValJsonFormat               | JSON 格式錯誤，請檢查             |
| 400    | E2004     | ValTypeConversionFailed     | 參數類型轉換失敗                  |
| 400    | E2020     | ValFieldRequired            | {field} 為必填項目                |
| 400    | E2030     | ValFieldOneof               | {field} 必須是 {param} 其中一個值 |
| 404    | E3C001    | CustomerNotFound            | 客戶不存在                        |
| 404    | E3COU004  | CouponNotFound              | 優惠券不存在或已被刪除            |
| 409    | E3CCOU004 | CustomerCouponAlreadyExists | 客戶優惠券已存在                  |
| 500    | E9001     | SysInternalError            | 系統發生錯誤，請稍後再試          |
| 500    | E9002     | SysDatabaseError            | 資料庫操作失敗                    |

---

## 資料表

- `customer_coupons`
- `customers`
- `coupons`

---

## Service 邏輯

1. 確認 `customer_id` 與 `coupon_id` 是否存在。
2. 確認 `customer_id` 和 `coupon_id` 組合是否已存在。
3. 根據 `period` 設定 `valid_from` 與 `valid_to`。
4. 建立 `customer_coupons` 資料。
5. 回傳新增結果。

---

## 注意事項

- `period` 未來可能會再擴充。
