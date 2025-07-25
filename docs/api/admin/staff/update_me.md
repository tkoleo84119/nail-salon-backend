## User Story

作為一位員工，我希望可以更新自己的個人資料，以確保聯絡資訊正確。

---

## Endpoint

**PATCH** `/api/admin/staff/me`

---

## 說明

- 僅允許員工更新自己的資料
- 目前僅支援更新 Email

---

## 權限

- 僅限已登入員工可呼叫

---

## Request

### Header

```http
Content-Type: application/json
Authorization: Bearer <access_token>
```

### Body

```json
{
  "email": "new-email@example.com"
}
```

- 至少需要提供一個欄位進行更新

### 驗證規則

| 欄位  | 規則                  | 說明       |
| ----- | --------------------- | ---------- |
| email | <li>選填<li>email格式 | 員工 Email |

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
    "isActive": true
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "email": "email格式不正確"
  }
}
```

#### 401 Unauthorized - 未登入/Token失效

```json
{
  "message": "無效的 accessToken"
}
```

#### 409 Conflict - Email已存在

```json
{
  "message": "此電子郵件已被註冊"
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

- `staff_users`

---

## Service 邏輯

1. 驗證請求至少需要提供一個欄位進行更新
2. 驗證 `staff_users` 是否存在
3. 如果有傳入 `email`，驗證 Email 是否唯一（其他人不可用，不包含自己）
4. 更新 `staff_users` 的欄位
5. 回傳更新後資訊

---

## 注意事項

- 未來可能會可以更新其他欄位。
- email 欄位需唯一，不可與其他員工重複。
