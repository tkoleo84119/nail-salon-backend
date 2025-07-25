## User Story

作為一位員工（`ADMIN` / `MANAGER` / `STYLIST`，不含 `SUPER_ADMIN`），我希望可以為自己新增美甲師個人資料（stylists），讓顧客預約時能看到我的專長與風格。

---

## Endpoint

**POST** `/api/admin/stylists/me`

---

## 說明

每位員工（`staff_user`）只能有一筆對應的 `stylist` 資料。

---

## 權限

- 僅限已登入之非 `SUPER_ADMIN` 員工（`ADMIN` / `MANAGER` / `STYLIST`）可呼叫

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
  "stylistName": "Jane 美甲師",
  "goodAtShapes": ["方形", "圓形"],
  "goodAtColors": ["裸色系", "粉嫩系"],
  "goodAtStyles": ["手繪", "簡約"],
  "isIntrovert": false
}
```

### 驗證規則

| 欄位         | 規則                                | 說明           |
| ------------ | ----------------------------------- | -------------- |
| stylistName  | <li>必填<li>長度大於1<li>長度小於50 | 美甲師顯示姓名 |
| goodAtShapes | <li>可選                            | 擅長指型       |
| goodAtColors | <li>可選                            | 擅長色系       |
| goodAtStyles | <li>可選                            | 擅長款式       |
| isIntrovert  | <li>可選<li>布林值                  | 是否I人        |

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "id": "18000000001",
    "staffUserId": "13984392823",
    "stylistName": "Jane 美甲師",
    "goodAtShapes": ["方形", "圓形"],
    "goodAtColors": ["裸色系", "粉嫩系"],
    "goodAtStyles": [],
    "isIntrovert": false
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "stylistName": "stylistName為必填"
  }
}
```

#### 401 Unauthorized - 未登入/Token失效

```json
{
  "message": "無效的 accessToken"
}
```

#### 403 Forbidden - 權限不足

```json
{
  "message": "權限不足，無法執行此操作"
}
```

#### 409 Conflict - 已存在

```json
{
  "message": "該員工已建立過美甲師資料，請使用修改功能"
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

- `stylists`
- `staff_users`

---

## Service 邏輯

1. 檢查 `stylists` 是否已有該 `staff_user_id` 對應資料：
   - 若已存在，回傳 409 Conflict。
   - 若尚未建立，執行新增。
2. 新增一筆 stylist 資料，並回傳。

---

## 注意事項

- 一個 staff_user 只可對應一筆 stylist 資料（由 DB unique 約束）。
- 欲修改資料時，請呼叫 PATCH 而非重複新增。

---
