## User Story

作為顧客，我希望能夠透過 refresh token 來取得新的 access token，方便維持登入狀態。

---

## Endpoint

**POST** `/api/auth/token/refresh`

---

## 說明

- 使用者提供 refresh token，若有效則簽發新的 access token（與 refresh token 一起送出驗證）。
- 若 token 無效、已過期或已被撤銷，則回傳錯誤。

---

## 權限

- 需提供合法的 refresh token。

---

## Request

### Header

Content-Type: application/json

### Body

```json
{
  "refreshToken": "your_refresh_token_here"
}
```

---

## Response

### 成功 200 OK

```json
{
  "accessToken": "new_access_token",
  "expiresIn": 3600
}
```

### 失敗

#### 401 Unauthorized - Token 無效或已過期

```json
{
  "message": "Refresh token 無效或已過期"
}
```

#### 500 Internal Server Error

```json
{
  "message": "系統發生錯誤，請稍後再試"
}
```

---

## 資料表

- `customer_tokens`

---

## Service 邏輯

1. 驗證 refreshToken 是否存在於資料庫中且未撤銷（`is_revoked=false`）。
2. 檢查該 token 是否過期（`expired_at > now()`）。
3. 若皆通過，發出新的 access token。

---

## 注意事項

- 暫時不考慮 refresh token 的每次使用都撤銷
