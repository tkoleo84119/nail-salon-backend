## User Story

作為顧客，我希望能夠取得我的單筆預約資訊，僅可查詢自己的預約。 若帶入非自己的 bookingId 或不存在的 id，皆回傳 404。

---

## Endpoint

**GET** `/api/bookings/{bookingId}`

---

## 說明

- 僅支援已登入顧客（access token 驗證）。
- 回傳對應預約的詳細資訊（包含服務、時段、門市、美甲師、狀態、備註等）。
- 只能查詢自己的預約，若非本人預約或不存在則回傳 404。
- 適用於「預約明細」或點擊預約清單查看詳情場景。

---

## 權限

- 僅顧客本人可查詢（JWT 驗證）。

---

## Request

### Header

```http
Authorization: Bearer <access_token>
```

### Path Parameter

| 參數      | 說明   |
| --------- | ------ |
| bookingId | 預約ID |

---

## Response

### 成功 200 OK

```json
{
  "data": {
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
    "services": [
      { "id": "9000000001", "name": "手部單色" },
      { "id": "9000000002", "name": "基礎保養" }
    ],
    "note": "這次想做奶茶色",
    "status": "SCHEDULED"
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

#### 404 Not Found - 無權限或不存在

```json
{
  "message": "查無此預約或無權限"
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
- `booking_details`
- `services`
- `time_slots`
- `stores`
- `stylists`

---

## Service 邏輯

1. 查詢該 bookingId 是否存在且屬於本人。
2. 回傳對應預約明細。

---

## 注意事項

- 僅允許本人查詢。
- 查無資料或查詢非本人皆回傳 404（防止資料外洩）。

