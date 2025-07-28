## User Story

作為員工，我希望可以幫客戶新增預約（Booking），以便協助客戶安排服務時段與內容。

---

## Endpoint

**POST** `/api/admin/stores/{storeId}/bookings`

---

## 說明

- 所有登入員工皆可使用。
- 由員工操作後台幫客戶新增預約。
- 須指定美甲師、時段、服務項目等必要資訊。

---

## 權限

- 任一已登入員工皆可使用（JWT 驗證）。

---

## Request

### Header

Authorization: Bearer <access_token>

### Path Parameter

| 參數    | 說明    |
| ------- | ------- |
| storeId | 門市 ID |

### Body

```json
{
  "customerId": "2000000001",
  "stylistId": "7000000001",
  "timeSlotId": "9000000001",
  "mainServiceId": "9000000010",
  "subServiceIds": ["9000000010", "9000000012"],
  "isChatEnabled": true,
  "note": "顧客想做法式+跳色"
}
```

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "id": "5000000001",
    "storeId": "8000000001",
    "storeName": "門市名稱",
    "stylistId": "2000000001",
    "stylistName": "美甲師名稱",
    "date": "2025-08-02",
    "timeSlotId": "3000000001",
    "startTime": "10:00",
    "endTime": "11:00",
    "mainServiceName": "主服務項目名稱",
    "subServiceNames": ["副服務項目名稱1", "副服務項目名稱2"],
    "isChatEnabled": true,
    "note": "這次想做奶茶色",
    "status": "SCHEDULED",
    "createdAt": "2025-07-24T10:00:00Z",
    "updatedAt": "2025-07-24T10:00:00Z"
  }
}
```

### 失敗

#### 400 Bad Request - 資料錯誤或時間衝突

```json
{
  "message": "該時段已被預約，請選擇其他時間"
}
```

#### 404 Not Found - 門市/客戶/美甲師/時段不存在

```json
{
  "message": "門市不存在或已被刪除"
}
```

#### 409 Conflict - 預約重複

```json
{
  "message": "該時段已被預約，請重新選擇"
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
- `customers`
- `stylists`
- `time_slots`
- `services`

---

## Service 邏輯

1. 驗證門市是否存在
2. 驗證員工是否有權限操作該門市
3. 驗證美甲師、時段、服務是否存在，且該時段可預約，同時服務時數加起來不能超過時段時數。
4. 建立 `bookings` 主檔與對應的 `booking_details`。
5. 標記 `time_slots.is_available = false`。
6. 回傳新建預約資料。

---

## 注意事項

- 僅支援單一時段對應一筆預約。
