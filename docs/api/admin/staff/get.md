## User Story

作為管理員，我希望可以查詢特定員工資料，並可一併取得對應的美甲師資料（若有），以便後台管理。

---

## Endpoint

**GET** `/api/admin/staff/{staffId}`

---

## 說明

- 僅限 `admin` 以上角色存取。
- 查詢指定 `staffId` 的員工帳號基本資料。
- 若該員工同時為美甲師（`stylists.staff_user_id`），則一併回傳 stylist 資訊。

---

## 權限

- 僅限 `admin` 以上角色 (`SUPER_ADMIN`, `ADMIN`)

---

## Request

### Header

Authorization: Bearer <access_token>

### Path Parameter

| 參數    | 說明        |
| ------- | ----------- |
| staffId | 員工帳號 ID |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "6000000002",
    "username": "stylist88",
    "email": "s88@salon.com",
    "role": "STYLIST",
    "isActive": true,
    "createdAt": "2025-06-01T08:00:00Z",
    "stylist": {
      "id": "7000000001",
      "name": "Bella",
      "goodAtShapes": ["方形"],
      "goodAtColors": ["粉色系"],
      "goodAtStyles": ["簡約風"],
      "isIntrovert": false
    }
  }
}
```

### 失敗

#### 401 Unauthorized - 未登入/Token失效

```json
{
  "message": "無效的 accessToken"
}
```

#### 403 Forbidden - 權限不足

```json
{
  "message": "無權限存取此資源"
}
```

#### 404 Not Found - 查無員工資料

```json
{
  "message": "查無此員工帳號"
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
- `stylists`

---

## Service 邏輯

1. 根據 `staffId` 查詢 `staff_users`。
2. 若該筆帳號存在對應 `stylists.staff_user_id`，一併查詢其美甲師資料。
3. 回傳合併結果。

---

## 注意事項

- 若無對應 stylist 資料，`stylist` 欄位為 null。

