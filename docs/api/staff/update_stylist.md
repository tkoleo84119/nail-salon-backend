## User Story

作為一位員工（`ADMIN` / `MANAGER` / `STYLIST`，不含 `SUPER_ADMIN`），我希望可以更新自己的美甲師個人資料（stylists），讓顧客可以看到我最新的專長與風格。

---

## Endpoint

**PATCH** `/api/stylists/me`

---

## 說明

每位員工（`staff_user`）僅能更新自己所屬的 `stylist` 資料。

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

### Body（可傳任一欄位，僅更新指定內容）

```json
{
  "stylistName": "Jane 美甲師",
  "goodAtShapes": ["方形", "橢圓形"],
  "goodAtColors": ["粉嫩系"],
  "goodAtStyles": ["簡約", "法式"],
  "isIntrovert": true
}
```

### 驗證規則

| 欄位         | 規則                                | 說明           |
| ------------ | ----------------------------------- | -------------- |
| stylistName  | <li>可選<li>長度大於1<li>長度小於50 | 美甲師顯示姓名 |
| goodAtShapes | <li>可選                            | 擅長指型       |
| goodAtColors | <li>可選                            | 擅長色系       |
| goodAtStyles | <li>可選                            | 擅長款式       |
| isIntrovert  | <li>可選<li>布林值                  | 是否I人        |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "18000000001",
    "staffUserId": "13984392823",
    "stylistName": "Jane 美甲師",
    "goodAtShapes": ["方形", "橢圓形"],
    "goodAtColors": ["粉嫩系"],
    "goodAtStyles": ["簡約", "法式"],
    "isIntrovert": true
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "stylistName": "美甲師姓名長度不可超過50字"
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

#### 404 Not Found - 尚未建立美甲師資料

```json
{
  "message": "尚未建立美甲師資料，請先新增"
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

1. 檢查 `stylists` 表是否有該 `staff_user_id` 對應資料：
   - 若尚未建立，回傳 404 Not Found。
2. 驗證更新欄位（如: 姓名長度）。
3. 僅更新傳入的指定欄位。
4. 回傳更新後的 stylist 資料。

---

## 注意事項

- 一個 `staff_user` 只可對應一筆 `stylist` 資料。
