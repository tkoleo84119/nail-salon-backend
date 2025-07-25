## User Story

作為顧客，我希望能夠取得所有門市的資料，支援分頁，方便快速查找可預約據點。

---

## Endpoint

**GET** `/api/stores`

---

## 說明

- 支援已登入顧客查詢所有門市。
- 支援分頁（limit、offset）。
- 僅回傳 `is_active=true` 的門市。
- 適用於預約選點、據點查詢等場景。

---

## 權限

- 僅顧客可查詢（JWT 驗證）。

---

## Request

### Header

Authorization: Bearer <access_token>

### Query Parameter

| 參數   | 型別 | 預設值 | 說明     |
| ------ | ---- | ------ | -------- |
| limit  | int  | 20     | 單頁筆數 |
| offset | int  | 0      | 起始筆數 |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "total": 10,
    "items": [
      {
        "id": "8000000001",
        "name": "大安旗艦店",
        "address": "台北市大安區復興南路一段100號",
        "phone": "02-1234-5678"
      },
    ]
  }
}
```

### 失敗

#### 500 Internal Server Error

```json
{
  "message": "系統發生錯誤，請稍後再試"
}
```

---

## 資料表

- `stores`

---

## Service 邏輯

1. 查詢 `is_active=true` 的門市。
2. 支援分頁（limit/offset）。
3. 回傳分頁資料與總數。

---

## 注意事項

- 僅回傳啟用門市。
