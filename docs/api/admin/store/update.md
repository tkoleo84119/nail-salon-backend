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

#### 錯誤總覽

| 狀態碼 | 錯誤碼   | 說明                                                         |
| ------ | -------- | ------------------------------------------------------------ |
| 401    | E1002    | 無效的 accessToken，請重新登入                               |
| 401    | E1003    | accessToken 缺失，請重新登入                                 |
| 401    | E1004    | accessToken 格式錯誤，請重新登入                             |
| 401    | E1005    | 未找到有效的員工資訊，請重新登入                             |
| 401    | E1006    | 未找到使用者認證資訊，請重新登入                             |
| 403    | E1010    | 權限不足，無法執行此操作                                     |
| 400    | E2001    | JSON 格式錯誤，請檢查                                        |
| 400    | E2003    | 至少需要提供一個欄位進行更新                                 |
| 400    | E2004    | 參數類型轉換失敗                                             |
| 400    | E2020    | {field} 為必填項目                                           |
| 400    | E2021    | {field} 長度至少需要 {param} 個字元                          |
| 400    | E2024    | {field} 長度最多只能有 {param} 個字元                        |
| 400    | E2031    | {field} 格式錯誤，請使用正確的台灣電話號碼格式 (0X-XXXXXXXX) |
| 400    | E2029    | {field} 必須是布林值                                         |
| 409    | E3STO003 | 門市已存在，請創建其他門市                                   |
| 500    | E9001    | 系統發生錯誤，請稍後再試                                     |
| 500    | E9002    | 資料庫操作失敗                                               |

#### 400 Bad Request - 驗證錯誤

```json
{
  "error": [
    {
      "code": "E2031",
      "message": "phone 格式錯誤，請使用正確的台灣電話號碼格式 (0X-XXXXXXXX)",
      "field": "phone"
    }
  ]
}
```

#### 401 Unauthorized - 未登入/Token失效

```json
{
  "error": [
    {
      "code": "E1002",
      "message": "無效的 accessToken，請重新登入"
    }
  ]
}
```

#### 403 Forbidden - 權限不足

```json
{
  "error": [
    {
      "code": "E1010",
      "message": "權限不足，無法執行此操作"
    }
  ]
}
```

#### 404 Not Found - 門市不存在

```json
{
  "error": [
    {
      "code": "E3STO003",
      "message": "門市已存在，請創建其他門市"
    }
  ]
}
```

#### 500 Internal Server Error - 系統錯誤

```json
{
  "error": [
    {
      "code": "E9001",
      "message": "系統發生錯誤，請稍後再試"
    }
  ]
}
```

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
