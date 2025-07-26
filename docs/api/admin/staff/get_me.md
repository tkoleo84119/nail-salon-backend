## User Story

作為員工，我希望可以取得自己的員工資料，方便查看登入帳號與基本資訊。

---

## Endpoint

**GET** `/api/admin/staff/me`

---

## 說明

- 提供當前登入員工自己的基本資料。
- 用於顯示後台帳號資訊（如個人資料設定頁）。
- 根據 JWT 中的 `staff_user_id` 載入對應資料。

---

## 權限

- 任一登入後的員工皆可使用（需 JWT 驗證）。

---

## Request

### Header

Authorization: Bearer <access_token>

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "6000000001",
    "username": "admin01",
    "email": "admin01@example.com",
    "role": "ADMIN",
    "isActive": true
  }
}
```

### 失敗

#### 401 Unauthorized - 未登入

```json
{
  "message": "請先登入後再操作"
}
```

#### 404 Not Found - 員工資料不存在

```json
{
  "message": "找不到對應的員工帳號"
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

1. 解析 JWT 取得 `staff_user_id`。
2. 查詢 `staff_users` 表是否存在該使用者。
3. 回傳基本資訊（帳號、信箱、角色、狀態、建立時間）。

---

## 注意事項

- 僅可查詢自己的帳號，無其他查詢條件。
