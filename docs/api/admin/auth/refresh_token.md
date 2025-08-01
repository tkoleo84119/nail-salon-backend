## User Story

作為一位員工，我希望能夠透過 `refreshToken` 來取得新的 `accessToken`，方便維持登入狀態。

---

## Endpoint

**POST** `/api/admin/auth/token/refresh`

---

## 說明

- 員工提供 `refreshToken`，若有效則簽發新的 `accessToken` 與 `user` 資料。
- 此為後台管理端使用，與前台顧客端獨立。

---

## 權限

- 需提供合法的 `refreshToken`。

---

## Request

### Header

* Content-Type: application/json

### Body 範例

```json
{
  "refreshToken": "your_refresh_token_here"
}
```

### 驗證規則

| 欄位         | 必填 | 其他規則            | 說明     |
| ------------ | ---- | ------------------- | -------- |
| refreshToken | 是   | <li>最大長度500字元 | 刷新令牌 |

---

## Response

### 成功 200 OK

```json
{
  "accessToken": "new_access_token",
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
```

### 錯誤總覽

| 狀態碼 | 錯誤碼 | 說明                                  |
| ------ | ------ | ------------------------------------- |
| 400    | E2001  | JSON 格式錯誤                         |
| 400    | E2020  | {field} 為必填項目                    |
| 400    | E2024  | {field} 長度最多只能有 {param} 個字元 |
| 401    | E1009  | Refresh token 無效或已過期            |
| 500    | E9001  | 系統發生錯誤，請稍後再試              |
| 500    | E9002  | 資料庫操作失敗，請稍後再試            |

#### 400 Bad Request - 輸入驗證失敗

```json
{
  "errors": [
    {
      "code": "E2020",
      "message": "refreshToken 為必填項目",
      "field": "refreshToken"
    },
    {
      "code": "E2024",
      "message": "refreshToken 長度最多只能有 500 個字元",
      "field": "refreshToken"
    }
  ]
}
```

#### 401 Unauthorized - Token 無效或已過期

```json
{
  "errors": [
    {
      "code": "E1009",
      "message": "Refresh token 無效或已過期"
    }
  ]
}
```

#### 500 Internal Server Error

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

- `staff_user_tokens`
- `staff_users`
- `staff_user_store_access`
- `stores`

---

### Service 邏輯

1. 驗證 refreshToken 是否存在於資料庫中。
   - `expired_at > now()`
   - `is_revoked = false`
2. 取得員工資料
3. 取得員工可存取的店家
4. 產生新的 access token
5. 回傳新的 access token

---

## 注意事項

- 暫時不考慮 refresh token 的每次使用都撤銷
