## User Story

作為一位員工，我希望能夠登入系統，以便存取後台功能。

---

## Endpoint

**POST** `/api/admin/auth/login`

---

## 說明

提供後台員工登入功能，並取得 `Access Token` 與 `Refresh Token`。

---

## 權限

- 不須預先認證

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Body 範例

```json
{
  "username": "admin001",
  "password": "hunter2"
}
```

### 驗證規則

| 欄位     | 必填 | 其他規則                            | 說明         |
| -------- | ---- | ----------------------------------- | ------------ |
| username | 是   | <li>不能為空字串<li>最大長度100字元 | 帳號（唯一） |
| password | 是   | <li>不能為空字串<li>最大長度100字元 | 密碼明文     |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "accessToken": "<jwt_access_token>",
    "refreshToken": "<secure_refresh_token>",
    "expiresIn": 3600
  }
}
```

---

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
| 401    | E1001  | AuthInvalidCredentials  | 帳號或密碼錯誤                        |
| 401    | E1005  | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入      |
| 401    | E1006  | AuthContextMissing      | 未找到使用者認證資訊，請重新登入      |
| 400    | E2001  | ValJsonFormat           | JSON 格式錯誤，請檢查                 |
| 400    | E2020  | ValFieldRequired        | {field} 為必填項目                    |
| 400    | E2024  | ValFieldMaxLength       | {field} 長度最多只能有 {param} 個字元 |
| 400    | E2036  | ValFieldNoBlank         | {field} 不能為空字串                  |
| 400    | E2004  | ValTypeConversionFailed | 參數類型轉換失敗                      |
| 500    | E9001  | SysInternalError        | 系統發生錯誤，請稍後再試              |
| 500    | E9002  | SysDatabaseError        | 資料庫操作失敗                        |

---

## 實作與流程

### 資料表

- `staff_users`
- `staff_user_tokens`

### Service 邏輯

1. 根據 `username` 查詢 `staff_users`
   - 確認是否存在
   - 檢查 `password_hash`（bcrypt）是否與 `password` 相符
   - 確認是否被停用 `is_active = false`
2. 產生 JWT（`access_token`）與 `refresh_token`（儲存於 `staff_user_tokens`）
3. 回傳登入結果

---

## 注意事項

- 密碼不會回傳，僅用於比對
- `refresh_token` 儲存於 DB 的同時，會設置 `expired_at` / `user_agent` / `ip_address`
- `expires_in` 表示 Access Token 的有效秒數，預設為 1 小時
- `storeList` 為該員工有權限的門市列表
