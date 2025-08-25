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

- Content-Type: application/json

### Body 範例

```json
{
  "refreshToken": "your_refresh_token_here"
}
```

### 驗證規則

| 欄位         | 必填 | 其他規則                            | 說明     |
| ------------ | ---- | ----------------------------------- | -------- |
| refreshToken | 是   | <li>不能為空字串<li>最大長度500字元 | 刷新令牌 |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "accessToken": "new_access_token",
    "expiresIn": 3600
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

| 狀態碼 | 錯誤碼 | 常數名稱                | 說明                                   |
| ------ | ------ | ----------------------- | -------------------------------------- |
| 401    | E1009  | AuthRefreshTokenInvalid | Refresh token 無效或已過期，請重新登入 |
| 400    | E2001  | ValJsonFormat           | JSON 格式錯誤，請檢查                  |
| 400    | E2020  | ValFieldRequired        | {field} 為必填項目                     |
| 400    | E2024  | ValFieldMaxLength       | {field} 長度最多只能有 {param} 個字元  |
| 400    | E2036  | ValFieldNoBlank         | {field} 不能為空字串                   |
| 500    | E9001  | SysInternalError        | 系統發生錯誤，請稍後再試               |
| 500    | E9002  | SysDatabaseError        | 資料庫操作失敗                         |

---

## 實作與流程

### 資料表

- `staff_user_tokens`
- `staff_users`

---

### Service 邏輯

1. 驗證 `refreshToken` 是否存在於資料庫中。
   - `expired_at > now()`
   - `is_revoked = false`
2. 取得員工資料，並驗證是否被停用 `is_active = false`
3. 產生新的 `accessToken`
4. 回傳新的 `accessToken`

---

## 注意事項

- 暫時不考慮 `refreshToken` 的每次使用都撤銷
