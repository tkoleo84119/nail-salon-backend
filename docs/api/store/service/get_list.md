## User Story

作為顧客，我希望能夠取得特定門市的服務項目，支援查詢條件，方便預約時選擇。

---

## Endpoint

**GET** `/api/stores/{storeId}/services`

---

## 說明

- 支援已登入顧客查詢指定門市的服務清單。
- 支援分頁（limit、offset）與 `isAddon` 條件（預設為 true）。
- 僅回傳 is_visible=true、is_active=true 的服務。
- 適用於預約服務選單、服務瀏覽等場景。

---

## 權限

- 開放所有顧客查詢，無需 JWT。

---

## Request

### Header

Authorization: Bearer <access_token>

### Path Parameter

| 參數    | 說明   |
| ------- | ------ |
| storeId | 門市ID |

### Query Parameter

| 參數    | 型別 | 預設值 | 說明                  |
| ------- | ---- | ------ | --------------------- |
| limit   | int  | 20     | 單頁筆數              |
| offset  | int  | 0      | 起始筆數              |
| isAddon | bool | true   | 僅查主/附加服務(選填) |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "total": 10,
    "items": [
      {
        "id": "9000000001",
        "name": "手部單色",
        "durationMinutes": 60,
        "isAddon": false,
        "note": "含基礎修型保養"
      },
    ]
  }
}
```

### 失敗

#### 404 Not Found - 門市不存在

```json
{
  "message": "指定門市不存在"
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

- `stores`
- `services`

---

## Service 邏輯

1. 驗證 storeId 是否存在。
2. 查詢該門市下 is_visible=true 且 is_active=true 的服務 (同時加上 isAddon 條件)。
3. 支援分頁（limit/offset）。
4. 回傳分頁資料與總數。

---

## 注意事項

- 僅回傳前台可見且啟用服務。
