## User Story

作為一位管理員，我希望能新增優惠券，方便維護優惠券資訊。

---

## Endpointˆ

**POST** `/api/admin/coupons`

---

## 說明

- 提供後台管理員新增優惠券功能。
- 優惠券名稱須唯一。
- 可設定優惠券顯示名稱、折扣率、折扣金額、是否啟用、備註。

---

## 權限

- 需要登入才可使用。
- 僅 `SUPER_ADMIN`、`ADMIN` 可操作。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Body 範例

```json
{
  "name": "新客優惠-八折",
  "displayName": "新客優惠",
  "code": "NEW_CUSTOMER_80",
  "discountRate": 0.8,
  "discountAmount": 100,
  "note": "只適用於新客"
}
```

- discountRate 和 discountAmount 至少傳入一個，但不能同時填寫。

### 驗證規則

| 欄位           | 必填 | 其他規則                            | 說明           |
| -------------- | ---- | ----------------------------------- | -------------- |
| name           | 是   | <li>不能為空字串<li>最大長度100字元 | 優惠券名稱     |
| displayName    | 是   | <li>不能為空字串<li>最大長度100字元 | 優惠券顯示名稱 |
| code           | 是   | <li>不能為空字串<li>最大長度100字元 | 優惠券代碼     |
| discountRate   | 否   | <li>最小值0.1<li>最大值0.99         | 折扣率         |
| discountAmount | 否   | <li>最小值1<li>最大值1000000        | 折扣金額       |
| note           | 否   | <li>最大長度255                     | 備註           |

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

| 狀態碼 | 錯誤碼   | 常數名稱                | 說明                                  |
| ------ | -------- | ----------------------- | ------------------------------------- |
| 401    | E1002    | AuthTokenInvalid        | 無效的 accessToken，請重新登入        |
| 401    | E1003    | AuthTokenMissing        | accessToken 缺失，請重新登入          |
| 401    | E1004    | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入      |
| 401    | E1005    | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入      |
| 401    | E1006    | AuthContextMissing      | 未找到使用者認證資訊，請重新登入      |
| 403    | E1010    | AuthPermissionDenied    | 權限不足，無法執行此操作              |
| 400    | E2001    | ValJsonFormat           | JSON 格式錯誤，請檢查                 |
| 400    | E2004    | ValTypeConversionFailed | 參數類型轉換失敗                      |
| 400    | E2020    | ValFieldRequired        | {field} 為必填項目                    |
| 400    | E2023    | ValFieldMinNumber       | {field} 最小值為 {param}              |
| 400    | E2024    | ValFieldStringMaxLength | {field} 長度最多只能有 {param} 個字元 |
| 400    | E2026    | ValFieldMaxNumber       | {field} 最大值為 {param}              |
| 400    | E2036    | ValFieldNoBlank         | {field} 不能為空字串                  |
| 400    | E3COU002 | CouponDiscountRequired  | 折數或折扣金額至少需要提供一個        |
| 400    | E3COU003 | CouponDiscountExclusive | 折數和折扣金額不能同時填寫            |
| 409    | E3COU005 | CouponNameAlreadyExists | 優惠券名稱已存在，請使用其他名稱      |
| 409    | E3COU006 | CouponCodeAlreadyExists | 優惠券代碼已存在，請使用其他代碼      |
| 500    | E9001    | SysInternalError        | 系統發生錯誤，請稍後再試              |
| 500    | E9002    | SysDatabaseError        | 資料庫操作失敗                        |

---

## 資料表

- `coupons`

---

## Service 邏輯

1. 確認 `name` 是否唯一。
2. 確認 `code` 是否唯一。
3. 建立 `coupons` 資料。
4. 回傳新增結果。

---

## 注意事項

- 優惠券名稱不可重複。
- 優惠券代碼不可重複。
