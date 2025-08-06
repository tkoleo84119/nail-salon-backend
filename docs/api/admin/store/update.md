## User Story

作為一位管理員，我希望能更新門市（store），以維護門市資訊。

---

## Endpoint

**PATCH** `/api/admin/stores/{storeId}`

---

## 說明

- 提供後台管理員更新門市功能。
- 僅允許修改名稱、地址、電話、是否啟用。
- `ADMIN` 只可修改自己有權限的門市。

---

## 權限

- 需要登入才可使用。
- 僅 `SUPER_ADMIN`、`ADMIN` 可操作。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameter

| 參數    | 說明   |
| ------- | ------ |
| storeId | 門市ID |

### Body 範例

```json
{
  "name": "松江南京分店",
  "address": "台北市中山區松江路123號",
  "phone": "02-88889999",
  "isActive": true
}
```

### 驗證規則

| 欄位     | 規則                                                               | 說明     |
| -------- | ------------------------------------------------------------------ | -------- |
| name     | <li>選填<li>長度大於1<li>長度小於100                               | 門市名稱 |
| address  | <li>選填<li>長度小於255                                            | 門市地址 |
| phone    | <li>選填<li>長度小於20<li>格式必須為台灣市話號碼 (例: 02-12345678) | 電話     |
| isActive | <li>選填<li>必須為布林值                                           | 是否啟用 |

- 欄位皆為選填，但至少需有一項。

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "8000000001",
    "name": "松江南京分店",
    "address": "台北市中山區松江路123號",
    "phone": "02-88889999",
    "isActive": true,
    "createdAt": "2025-01-01T00:00:00.000Z",
    "updatedAt": "2025-01-01T00:00:00.000Z"
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
| 401    | E1002    | AuthInvalidCredentials  | 無效的 accessToken，請重新登入        |
| 401    | E1003    | AuthTokenMissing        | accessToken 缺失，請重新登入          |
| 401    | E1004    | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入      |
| 401    | E1005    | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入      |
| 401    | E1006    | AuthContextMissing      | 未找到使用者認證資訊，請重新登入      |
| 403    | E1010    | AuthPermissionDenied    | 權限不足，無法執行此操作              |
| 400    | E2001    | ValJsonFormat           | JSON 格式錯誤，請檢查                 |
| 400    | E2002    | ValPathParamMissing     | 路徑參數缺失，請檢查                  |
| 400    | E2003    | ValAllFieldsEmpty       | 至少需要提供一個欄位進行更新          |
| 400    | E2004    | ValTypeConversionFailed | 參數類型轉換失敗                      |
| 400    | E2023    | ValFieldMinNumber       | {field} 最小值為 {param}              |
| 400    | E2024    | ValFieldStringMaxLength | {field} 長度最多只能有 {param} 個字元 |
| 400    | E2026    | ValFieldMaxNumber       | {field} 最大值為 {param}              |
| 400    | E2029    | ValFieldBoolean         | {field} 必須是布林值                  |
| 409    | E3STO003 | StoreAlreadyExists      | 門市已存在，請創建其他門市            |
| 500    | E9001    | SysInternalError        | 系統發生錯誤，請稍後再試              |
| 500    | E9002    | SysDatabaseError        | 資料庫操作失敗                        |

---

## 資料表

- `stores`
- `staff_user_store_access`

---

## Service 邏輯

1. 再次驗證 `role` 是否為 `SUPER_ADMIN` 或 `ADMIN`。
2. 驗證至少一個欄位有更新。
3. 驗證 `store` 是否存在。
4. 若 `role` 為 `ADMIN`，則驗證 `store` 是否為自己有權限的門市。
5. 若 `name` 有更新，則驗證 `name` 是否唯一（不包含自己）。
6. 更新 `stores` 資料。
7. 回傳更新結果。

---

## 注意事項

- 門市名稱不可重複（不包含自己）。
- 僅允許 name、address、phone、isActive 欄位修改。
