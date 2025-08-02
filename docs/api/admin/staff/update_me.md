## User Story

作為一位員工，我希望可以更新自己的個人資料，以確保聯絡資訊正確。

---

## Endpoint

**PATCH** `/api/admin/staff/me`

---

## 說明

- 僅允許員工更新自己的資料

---

## 權限

- 需要登入才可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Body 範例

```json
{
  "email": "new-email@example.com"
}
```

### 驗證規則

| 欄位 | 必填 | 其他規則      | 說明       |
| ---- | ---- | ------------- | ---------- |
| role | 否   | <li>email格式 | 員工 Email |

- 至少需要提供一個欄位進行更新

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "13984392823",
    "username": "staff_amy",
    "email": "new-email@example.com",
    "role": "STYLIST",
    "isActive": true,
    "createdAt": "2025-06-01T08:00:00+08:00",
    "updatedAt": "2025-06-01T08:00:00+08:00"
  }
}
```

### 錯誤處理

#### 錯誤總覽

| 狀態碼 | 錯誤碼   | 說明                                       |
| ------ | -------- | ------------------------------------------ |
| 401    | E1002    | 無效的 accessToken，請重新登入             |
| 401    | E1003    | accessToken 缺失，請重新登入               |
| 401    | E1004    | accessToken 格式錯誤，請重新登入           |
| 401    | E1005    | 未找到有效的員工資訊，請重新登入           |
| 401    | E1006    | 未找到使用者認證資訊，請重新登入           |
| 400    | E2001    | JSON 格式錯誤，請檢查                      |
| 400    | E2003    | 至少需要提供一個欄位進行更新               |
| 400    | E2004    | 參數類型轉換失敗                           |
| 400    | E2027    | {field} 格式錯誤，請使用正確的電子郵件格式 |
| 404    | E3STA005 | 員工帳號不存在                             |
| 500    | E9001    | 系統發生錯誤，請稍後再試                   |
| 500    | E9002    | 資料庫操作失敗                             |

#### 400 Bad Request - 驗證錯誤

```json
{
  "errors": [
    {
      "code": "E2003",
      "message": "至少需要提供一個欄位進行更新"
    }
  ]
}
```

#### 401 Unauthorized - 未登入/Token失效

```json
{
  "errors": [
    {
      "code": "E1002",
      "message": "無效的 accessToken"
    }
  ]
}
```

#### 404 Not Found - 員工帳號不存在

```json
{
  "errors": [
    {
      "code": "E3STA005",
      "message": "員工帳號不存在"
    }
  ]
}
```

#### 500 Internal Server Error - 系統發生錯誤

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

## 資料表

- `staff_users`

---

## Service 邏輯

1. 驗證請求至少需要提供一個欄位進行更新。
2. 驗證 `staff_users` 是否存在。
3. 更新 `staff_users` 的欄位。
4. 回傳更新後資訊。

---

## 注意事項

- 未來可能會可以更新其他欄位。
