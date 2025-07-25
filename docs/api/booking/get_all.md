## User Story

作為顧客，我希望能夠取得我的預約資訊，並可依狀態分頁查詢，方便管理所有預約。

---

## Endpoint

**GET** `/api/bookings`

---

## 說明

- 僅支援已登入顧客（access token 驗證）。
- 回傳當前顧客所有預約資訊，支援分頁與依狀態查詢。
- 適用於「我的預約」、「預約管理」等場景。

---

## 權限

- 僅顧客本人可查詢（JWT 驗證）。

---

## Request

### Header

```http
Authorization: Bearer <access_token>
```

### Query Parameter

| 參數   | 型別   | 預設值    | 說明                          |
| ------ | ------ | --------- | ----------------------------- |
| limit  | int    | 20        | 單頁筆數                      |
| offset | int    | 0         | 起始筆數                      |
| status | string | SCHEDULED | 預約狀態 (可多選，用逗號分隔) |

- status 支援：`SCHEDULED`, `CANCELLED`, `COMPLETED` 等。
- 可只查詢特定狀態。

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "total": 10,
    "items": [
      {
        "id": "5000000001",
        "storeId": "8000000001",
        "storeName": "大安旗艦店",
        "stylistId": "2000000001",
        "stylistName": "Ava",
        "date": "2025-08-02",
        "timeSlot": {
          "id": "3000000001",
          "startTime": "10:00",
          "endTime": "12:00"
        },
        "status": "SCHEDULED"
      },
    ]
  }
}
```

### 失敗

#### 401 Unauthorized - 未登入/Token 失效

```json
{
  "message": "無效的 accessToken"
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

- `bookings`
- `schedules`
- `time_slots`
- `stores`
- `stylists`

---

## Service 邏輯

1. 依傳入條件查詢 bookings（分頁、狀態等）。
2. 關聯查詢對應時段、美甲師、門市資訊。
3. 回傳分頁結果。

---

## 注意事項

- 僅允許本人查詢。
- 狀態可多選（如 `status=SCHEDULED,COMPLETED`）。
- 順序依預約日期排序(DESC)。
