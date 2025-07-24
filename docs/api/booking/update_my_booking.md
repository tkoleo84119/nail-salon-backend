## User Story

作為顧客，我希望可以更改自己的預約（包含不同門市、不同美甲師、不同服務、不同附加服務），彈性調整預約內容。

---

## Endpoint

**PATCH** `/api/bookings/{bookingId}/me`

---

## 說明

- 僅支援已登入顧客（access token 驗證）。
- 可變更門市、美甲師、日期、時段、服務（可含多個）、備註等。
- 異動時會自動檢查新時段與服務項目的可用性與衝突。

---

## 權限

- 僅顧客本人可修改自己預約（JWT 驗證）。

---

## Request

### Header

```http
Content-Type: application/json
Authorization: Bearer <access_token>
```

### Path Parameter

| 參數      | 說明   |
| --------- | ------ |
| bookingId | 預約ID |

### Body

```json
{
  "storeId": "8000000002",
  "stylistId": "2000000005",
  "timeSlotId": "3000000020",
  "mainServiceId": "9000000001",
  "subServiceIds": ["9000000002", "9000000003"],
  "isChatEnabled": true,
  "note": "這次想做奶茶色"
}
```

- 欄位皆為選填，至少需要傳入一個欄位

### 驗證規則

| 欄位          | 規則                                      | 說明            |
| ------------- | ----------------------------------------- | --------------- |
| storeId       | <li>選填                                  | 預約門市        |
| stylistId     | <li>選填                                  | 美甲師ID        |
| timeSlotId    | <li>選填                                  | 時段ID          |
| mainServiceId | <li>選填                                  | 主服務項目ID    |
| subServiceIds | <li>選填<li>陣列<li>最多5項<li>可為空陣列 | 附加服務項目IDs |
| isChatEnabled | <li>選填                                  | 是否要聊天      |
| note          | <li>選填<li>最長500字                     | 備註說明        |


- storeId、stylistId、timeSlotId、mainServiceId、subServiceIds 必須一起傳入
- subServiceIds 可以傳入空陣列 => 代表不加任何附加服務

---

## Response

### 成功 200 OK

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

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "note": "note長度最多500個字元"
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
  "message": "權限不足，僅限本人操作"
}
```

#### 404 Not Found - 指定資源不存在

```json
{
  "message": "指定時段、服務或預約不存在"
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

---

## Service 邏輯

1. 驗證預約是否存在且屬於本人，並且預約狀態為 `SCHEDULED`。
2. 若有傳入門市、美甲師、時段、服務，則驗證異動後的門市、美甲師、時段、服務。
   1. 驗證門市、美甲師、時段、服務是否存在
   2. 驗證時段是否可用
   3. 驗證服務是否可用
   4. 驗證時段時間是否足夠
   5. 驗證附加服務是否可用
3. 更新預約內容（bookings、booking_details）。
4. 回傳最新預約資訊。

---

## 注意事項

- 僅支援本人預約內容異動。
- 異動時需重新檢查時段、服務是否可用。
