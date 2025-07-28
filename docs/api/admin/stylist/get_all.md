## User Story

作為員工，我希望可以查詢某門市下的所有美甲師資料，並支援條件查詢與分頁，以便管理與選擇可排班對象。

---

## Endpoint

**GET** `/api/admin/stores/{storeId}/stylists`

---

## 說明

- 所有登入員工皆可查詢。
- 支援條件搜尋：姓名、是否內向
- 支援分頁（limit、offset）。

---

## 權限

- 任一已登入員工皆可存取（JWT 驗證）。

---

## Request

### Header

Authorization: Bearer <access_token>

### Path Parameters

| 參數    | 說明    |
| ------- | ------- |
| storeId | 門市 ID |

### Query Parameters

| 參數        | 型別   | 必填 | 預設值 | 說明                |
| ----------- | ------ | ---- | ------ | ------------------- |
| name        | string | 否   |        | 模糊查詢姓名        |
| isIntrovert | bool   | 否   |        | 是否為內向者（I人） |
| limit       | int    | 否   | 20     | 單頁筆數            |
| offset      | int    | 否   | 0      | 起始筆數            |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "total": 2,
    "items": [
      {
        "id": "7000000001",
        "staffUserId": "6000000010",
        "name": "Ariel",
        "goodAtShapes": ["方形"],
        "goodAtColors": ["裸色系"],
        "goodAtStyles": ["簡約風"],
        "isIntrovert": false
      }
    ]
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

#### 403 Forbidden - 無權限

```json
{
  "message": "權限不足，無法執行此操作"
}
```

#### 404 Not Found - 門市不存在

```json
{
  "message": "門市不存在或已被刪除"
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
- `store_access`
- `stores`

---

## Service 邏輯

1. 驗證 `storeId` 是否存在。
2. 驗證員工是否擁有該門市存取權限。
3. 透過 `store_access` 表查詢與該門市有關聯的 `staff_user_id`。
4. JOIN `stylists` 表並依查詢條件過濾：
   - 關鍵字模糊搜尋姓名
   - is_introvert 過濾
5. 加入 `limit` / `offset` 處理分頁。
6. 回傳總筆數與清單。

---
