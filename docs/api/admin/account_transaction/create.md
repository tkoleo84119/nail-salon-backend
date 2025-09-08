## User Story

作為一位管理員，我希望能新增帳戶交易紀錄，方便管理帳戶與資金流動。

---

## Endpoint

**POST** `/api/admin/stores/{storeId}/accounts/{accountId}/transactions`

---

## 說明

- 提供後台管理員新增帳戶交易紀錄功能。

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
  "transactionDate": "2025-01-01T00:00:00+08:00",
  "type": "INCOME",
  "amount": 100,
  "note": "備註"
}
```

### 驗證規則

| 欄位            | 必填 | 其他規則                         |
| --------------- | ---- | -------------------------------- |
| transactionDate | 是   | <li>格式為合格的 ISO 8601 格式   |
| type            | 是   | <li>只能為 INCOME 或 EXPENSE     |
| amount          | 是   | <li>最小值為1<li>最大值為1000000 |
| note            | 否   | <li>最大長度255字元              |

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "id": "8000000001"
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

| 狀態碼 | 錯誤碼  | 常數名稱                | 說明                                         |
| ------ | ------- | ----------------------- | -------------------------------------------- |
| 401    | E1002   | AuthTokenInvalid        | 無效的 accessToken，請重新登入               |
| 401    | E1003   | AuthTokenMissing        | accessToken 缺失，請重新登入                 |
| 401    | E1004   | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入             |
| 401    | E1005   | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入             |
| 401    | E1006   | AuthContextMissing      | 未找到使用者認證資訊，請重新登入             |
| 403    | E1010   | AuthPermissionDenied    | 權限不足，無法執行此操作                     |
| 400    | E2001   | ValJsonFormat           | JSON 格式錯誤，請檢查                        |
| 400    | E2002   | ValPathParamMissing     | 路徑參數缺失，請檢查                         |
| 400    | E2004   | ValTypeConversionFailed | 參數類型轉換失敗                             |
| 400    | E2020   | ValFieldRequired        | {field} 欄位為必填項目                       |
| 400    | E2024   | ValFieldStringMaxLength | {field} 長度最多只能有 {param} 個字元        |
| 400    | E2030   | ValFieldOneof           | {field} 必須是 {param} 其中一個值            |
| 400    | E2037   | ValFieldISO8601Format   | {field} 格式錯誤，請使用正確的 ISO 8601 格式 |
| 404    | E3ACC01 | AccountNotFound         | 帳戶不存在或已被刪除                         |
| 400    | E3ACC02 | AccountNotBelongToStore | 帳戶不屬於指定的門市                         |
| 500    | E9001   | SysInternalError        | 系統發生錯誤，請稍後再試                     |
| 500    | E9002   | SysDatabaseError        | 資料庫操作失敗                               |

---

## 資料表

- `accounts`
- `account_transactions`

---

## Service 邏輯

1. 檢查門市權限。
2. 檢查 `accounts` 是否存在。
3. 檢查 `account` 是否屬於該門市。
4. 建立 `account_transactions` 資料。
5. 回傳新增結果。
