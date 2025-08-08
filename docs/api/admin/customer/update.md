## User Story

作為一位員工，我希望能更新顧客（customer），以維護顧客資訊。

---

## Endpoint

**PATCH** `/api/admin/customers/{customerId}`

---

## 說明

- 提供後台員工更新顧客功能。
- 僅允許修改門市備註、顧客等級、是否列入黑名單。

---

## 權限

- 需要登入才可使用。
- 所有角色皆可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameter

| 參數       | 說明   |
| ---------- | ------ |
| customerId | 顧客ID |

### Body 範例

```json
{
  "storeNote": "門市備註",
  "level": "VIP",
  "isBlacklisted": true
}
```

### 驗證規則

| 欄位          | 必填 | 其他規則                         | 說明           |
| ------------- | ---- | -------------------------------- | -------------- |
| storeNote     | 否   | <li>長度小於255                  | 門市備註       |
| level         | 否   | <li>格式必須為 NORMAL, VIP, VVIP | 顧客等級       |
| isBlacklisted | 否   | <li>必須為布林值                 | 是否列入黑名單 |

- 欄位皆為選填，但至少需有一項。

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "8000000001",
    "name": "王小明",
    "phone": "0912345678",
    "birthday": "2000-01-01",
    "city": "台北市",
    "level": "NORMAL",
    "isBlacklisted": false,
    "lastVisitAt": "2025-01-01T00:00:00+08:00",
    "updatedAt": "2025-01-01T00:00:00+08:00"
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

| 狀態碼 | 錯誤碼 | 常數名稱                | 說明                                  |
| ------ | ------ | ----------------------- | ------------------------------------- |
| 401    | E1002  | AuthInvalidCredentials  | 無效的 accessToken，請重新登入        |
| 401    | E1003  | AuthTokenMissing        | accessToken 缺失，請重新登入          |
| 401    | E1004  | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入      |
| 401    | E1005  | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入      |
| 401    | E1006  | AuthContextMissing      | 未找到使用者認證資訊，請重新登入      |
| 403    | E1010  | AuthPermissionDenied    | 權限不足，無法執行此操作              |
| 400    | E2001  | ValJsonFormat           | JSON 格式錯誤，請檢查                 |
| 400    | E2002  | ValPathParamMissing     | 路徑參數缺失，請檢查                  |
| 400    | E2003  | ValAllFieldsEmpty       | 至少需要提供一個欄位進行更新          |
| 400    | E2004  | ValTypeConversionFailed | 參數類型轉換失敗                      |
| 400    | E2024  | ValFieldStringMaxLength | {field} 長度最多只能有 {param} 個字元 |
| 400    | E2029  | ValFieldBoolean         | {field} 必須是布林值                  |
| 400    | E2030  | ValFieldOneof           | {field} 必須是 {param} 其中一個值     |
| 404    | E3C001 | CustomerNotFound        | 客戶不存在                            |
| 500    | E9001  | SysInternalError        | 系統發生錯誤，請稍後再試              |
| 500    | E9002  | SysDatabaseError        | 資料庫操作失敗                        |

---

## 資料表

- `customers`

---

## Service 邏輯

1. 驗證至少一個欄位有更新。
2. 驗證客戶是否存在。
3. 更新 `customers` 資料。
3. 回傳更新結果。

---

## 注意事項

- 僅允許 storeNote、level、isBlacklisted 欄位修改。
