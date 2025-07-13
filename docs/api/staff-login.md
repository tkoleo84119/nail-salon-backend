## User Story

作為一位員工（`SUPER_ADMIN` / `ADMIN` / `MANAGER` / `STYLIST`），我希望能夠登入系統，以便存取後台功能。

---

## Endpoint

**POST** `/api/staff/login`

---

## 說明

供後台使用者（`staff`）登入系統，並取得 `Access Token` 與 `Refresh Token`。

---

## 權限

- 不須認證

---

## Request

### Header

```
Content-Type: application/json
```

### Body

```json
{
  "username": "admin001",
  "password": "hunter2"
}
```

### 驗證規則

| 欄位       | 規則       | 說明                  |
| -------- | -------- | ------------------- |
| username | required | 帳號（唯一）              |
| password | required | 密碼明文（將與 DB hash 比對） |

---

## Response

### 成功 200 OK

```json
{
  "access_token": "<jwt_access_token>",
  "refresh_token": "<secure_refresh_token>",
  "expires_in": 3600,
  "user": {
    "id": "139842394",
    "username": "admin001",
    "role": "ADMIN",
    "store_list": [
      {
        "id": 1,
        "name": "門市1"
      }
    ]
  }
}
```

### 失敗

#### 401 Unauthorized

```json
{
  "error": "invalid username or password"
}
```

---

## 實作

### 資料表

- `staff_users`
- `staff_user_tokens`
- `staff_user_store_access`
- `stores`

### Service 邏輯

1. 根據 `username` 和 `is_active = true` 查詢 `staff_users`
2. 檢查 `password_hash`（bcrypt）
3. 產生 JWT（`access_token`）與 `refresh_token`（儲存於 `staff_user_tokens`）
4. 查詢該員工可存取的店家（`staff_user_store_access`）
  - 如果是 `SUPER_ADMIN`，則查詢 `stores` 回傳所有店家
5. 回傳登入結果

---

## 注意事項

- 密碼不可回傳（僅比對）
- `refresh_token` 儲存於 DB，並設置 `expired_at` / `user_agent` / `ip`
- `expires_in` 表示 Access Token 的有效秒數
- `store_ids` 為該員工可查看的門市列表
