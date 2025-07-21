## User Story

作為顧客，我希望可以預約喜歡的美甲師與服務項目，方便安排我的美甲時段。

---

## Endpoint

**POST** `/api/bookings/me`

---

## 說明

- 僅支援已登入顧客（access token 驗證）。
- 顧客可指定門市、美甲師、日期、時段與服務項目進行預約。
- 支援備註欄位，可讓顧客描述本次需求。
- 預約後會自動檢查時段可用性與服務時數。

---

## 權限

- 僅顧客本人可預約（JWT 驗證）。

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
  "storeId": "8000000001",
  "stylistId": "2000000001",
  "timeSlotId": "3000000001",
  "mainServiceId": "9000000001",
  "subServiceIds": ["9000000002", "9000000003"],
  "isChatEnabled": true,
  "note": "這次想做奶茶色"
}
```

### 驗證規則

| 欄位          | 規則                        | 說明          |
| ------------- | --------------------------- | ------------- |
| storeId       | <li>必填                    | 預約門市      |
| stylistId     | <li>必填                    | 美甲師ID      |
| timeSlotId    | <li>必填                    | 時段ID        |
| mainServiceId | <li>必填                    | 主服務項目ID  |
| subServiceIds | <li>選填<li>陣列<li>最多5項 | 副服務項目IDs |
| isChatEnabled | <li>選填                    | 是否要聊天    |
| note          | <li>選填<li>最長500字       | 備註說明      |

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
    "status": "SCHEDULED"
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "storeId": "storeId為必填項目"
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
  "message": "權限不足，僅限本人預約"
}
```

#### 404 Not Found - 指定資源不存在

```json
{
  "message": "指定時段或服務不存在"
}
```

#### 409 Conflict - 時段已被預約

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
- `time_slots`
- `services`
- `stylists`
- `stores`

---

## Service 邏輯

1. 驗證門市、美甲師、時段、服務是否存在。
2. 驗證時段可預約（不可重複預約）。
3. 驗證時段可做該項主服務。
4. 建立預約資料（bookings、booking_details）。
5. 回傳預約資訊。

---

## 注意事項

- 僅支援本人預約。
