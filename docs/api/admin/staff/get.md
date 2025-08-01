## User Story

作為管理員，我希望可以查詢特定員工資料，以便後台管理。

---

## Endpoint

**GET** `/api/admin/staff/{staffId}`

---

## 說明

- 查詢指定 `staffId` 的員工帳號基本資料。
- 若該員工同時為美甲師（`stylists.staff_user_id`），則一併回傳 stylist 資訊。

---

## 權限

- 需要登入才可使用。
- `SUPER_ADMIN` 與 `ADMIN` 可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

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
    "updatedAt": "2025-06-01T08:00:00Z",
    "stylist": {
      "id": "7000000001",
      "name": "Bella",
      "goodAtShapes": ["方形"],
      "goodAtColors": ["粉色系"],
      "goodAtStyles": ["簡約風"],
      "isIntrovert": false,
      "createdAt": "2025-06-01T08:00:00Z",
      "updatedAt": "2025-06-01T08:00:00Z"
    }
  }
}
```

### 錯誤處理

#### 錯誤總覽

| 狀態碼 | 錯誤碼   | 說明                             |
| ------ | -------- | -------------------------------- |
| 401    | E1002    | 無效的 accessToken，請重新登入   |
| 401    | E1003    | accessToken 缺失，請重新登入     |
| 401    | E1004    | accessToken 格式錯誤，請重新登入 |
| 401    | E1005    | 未找到有效的員工資訊，請重新登入 |
| 401    | E1006    | 未找到使用者認證資訊，請重新登入 |
| 403    | E1010    | 權限不足，無法執行此操作         |
| 400    | E2004    | 參數類型轉換失敗                 |
| 400    | E2020    | {field} 為必填項目               |
| 404    | E3STA005 | 員工帳號不存在                   |
| 404    | E3STY001 | 美甲師資料不存在                 |
| 500    | E9001    | 系統發生錯誤，請稍後再試         |
| 500    | E9002    | 資料庫操作失敗                   |

#### 400 Bad Request - 參數類型轉換失敗

```json
{
  "error": {
    "code": "E2004",
    "message": "staffId 類型轉換失敗",
    "field": "staffId"
  }
}
```

#### 401 Unauthorized - 未登入/Token失效

```json
{
  "error": {
    "code": "E1002",
    "message": "無效的 accessToken，請重新登入"
  }
}
```

#### 403 Forbidden - 權限不足

```json
{
  "error": {
    "code": "E1010",
    "message": "權限不足，無法執行此操作"
  }
}
```

#### 404 Not Found - 查無員工資料

```json
{
  "error": {
    "code": "E3STA005",
    "message": "員工帳號不存在"
  }
}
```

#### 500 Internal Server Error - 系統發生錯誤

```json
{
  "error": {
    "code": "E9001",
    "message": "系統發生錯誤，請稍後再試"
  }
}
```

---

## 資料表

- `staff_users`
- `stylists`

---

## Service 邏輯

1. 根據 `staffId` 查詢 `staff_users`。
2. 若該帳號非 `SUPER_ADMIN`，則查詢 `stylists` 資料。
3. 回傳合併結果。

---

## 注意事項

- 若無對應 stylist 資料，`stylist` 欄位為 null。
