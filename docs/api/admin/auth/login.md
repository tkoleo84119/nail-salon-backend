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

* Content-Type: application/json

### Body 範例

```json
{
  "username": "admin001",
  "password": "hunter2"
}
```

### 驗證規則

| 欄位     | 必填 | 其他規則            | 說明         |
| -------- | ---- | ------------------- | ------------ |
| username | 是   | <li>最大長度100字元 | 帳號（唯一） |
| password | 是   | <li>最大長度100字元 | 密碼明文     |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "accessToken": "<jwt_access_token>",
    "refreshToken": "<secure_refresh_token>",
    "expiresIn": 3600,
    "user": {
      "id": "139842394",
      "username": "admin001",
      "role": "ADMIN",
      "storeList": [
        {
          "id": "1",
          "name": "門市1"
        }
      ]
    }
  }
}
```

---

### 錯誤處理

#### 錯誤總覽

| 狀態碼 | 錯誤碼 | 說明                                  |
| ------ | ------ | ------------------------------------- |
| 400    | E2001  | JSON 格式錯誤，請檢查                 |
| 400    | E2020  | {field} 為必填項目                    |
| 400    | E2024  | {field} 長度最多只能有 {param} 個字元 |
| 401    | E1001  | 帳號或密碼錯誤                        |
| 500    | E9001  | 系統發生錯誤，請稍後再試              |
| 500    | E9002  | 資料庫操作失敗                        |

#### 400 Bad Request - 輸入驗證失敗

```json
{
  "errors": [
    {
      "code": "E2020",
      "message": "username 欄位為必填項目",
      "field": "username"
    },
    {
      "code": "E2024",
      "message": "password 長度最多只能有 100 個字元",
      "field": "password"
    }
  ]
}
```

#### 401 Unauthorized - 認證失敗

```json
{
  "errors": [
    {
      "code": "E1001",
      "message": "帳號或密碼錯誤"
    }
  ]
}
```

#### 500 Internal Server Error - 系統錯誤

```json
{
  "errors": [
    {
      "code": "E9001",
      "message": "系統發生錯誤，請稍後再試"
    }
  ]
}
```

---

## 實作與流程

### 資料表

- `staff_users`
- `staff_user_store_access`
- `stores`
- `staff_user_tokens`

### Service 邏輯

1. 根據 `username` 查詢 `staff_users`
   - 確認是否存在
   - 檢查 `password_hash`（bcrypt）是否與 `password` 相符
   - 確認是否被停用 `is_active = false`
2. 查詢該員工可存取的店家（`staff_user_store_access`）
   - 如果是 `SUPER_ADMIN`，則查詢 `stores` 回傳所有店家
   - 不論店家是否被停用，都會回傳
3. 產生 JWT（`access_token`）與 `refresh_token`（儲存於 `staff_user_tokens`）
4. 回傳登入結果

---

## 注意事項

- 密碼不會回傳，僅用於比對
- `refresh_token` 儲存於 DB 的同時，會設置 `expired_at` / `user_agent` / `ip_address`
- `expires_in` 表示 Access Token 的有效秒數，預設為 1 小時
- `storeList` 為該員工有權限的門市列表
