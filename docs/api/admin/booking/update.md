## User Story

作為員工，我希望可以更新某筆顧客預約（Booking），例如更換時段、修改服務、調整備註等，以因應顧客變動需求。

---

## Endpoint

**PATCH** `/api/admin/stores/{storeId}/bookings/{bookingId}`

---

## 說明

- 所有登入員工皆可調用。
- 可部分修改指定預約的內容：時段、服務項目、備註、是否啟用聊天、使用產品。
- 修改時段會檢查是否仍可預約（`is_available=true`）。

---

## 權限

- 任一已登入員工皆可使用（JWT 驗證）。

---

## Request

### Header

Authorization: Bearer <access_token>
Content-Type: application/json

### Path Parameters

| 參數      | 說明    |
| --------- | ------- |
| storeId   | 門市 ID |
| bookingId | 預約 ID |

### Body（可選任意欄位）

```json
{
  "timeSlotId": "9000000002",
  "mainServiceId": "9000000010",
  "subServiceIds": ["9000000011", "9000000012"],
  "isChatEnabled": false,
  "note": "顧客改約下午並不開啟聊天"
}
```

### 驗證規則

| 欄位          | 規則                                      | 說明            |
| ------------- | ----------------------------------------- | --------------- |
| timeSlotId    | <li>選填                                  | 時段ID          |
| mainServiceId | <li>選填                                  | 主服務項目ID    |
| subServiceIds | <li>選填<li>陣列<li>最多5項<li>可為空陣列 | 附加服務項目IDs |
| isChatEnabled | <li>選填                                  | 是否要聊天      |
| note          | <li>選填<li>最長500字                     | 備註說明        |


- timeSlotId、mainServiceId、subServiceIds 必須一起傳入
- subServiceIds 可以傳入空陣列 => 代表不加任何附加服務

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "3000000001"
  }
}
```

### 失敗

#### 400 Bad Request - 欄位錯誤或時間衝突

```json
{
  "message": "該時段已被預約，請重新選擇"
}
```

#### 404 Not Found - 門市或預約不存在

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

- `bookings`
- `booking_details`
- `time_slots`
- `services`

---

## Service 邏輯

1. 若有傳入任一欄位，其他兩個也必須傳入(timeSlotId、mainServiceId、subServiceIds)
2. 驗證門市是否存在
3. 驗證員工是否有權限更新該門市預約
4. 驗證預約是否存在並隸屬於該門市，且預約狀態為 `SCHEDULED`
5. 驗證時段是否還可預約，並驗證客戶是否已經預約其他時段
6. 驗證服務項目是否存在
7. 若更新 `timeSlotId`，需：
   - 原時段釋放 `is_available=true`
   - 將新時段標記為 `is_available=false`
8. 若更新 `mainServiceId` 或 `subServiceIds`，則刪除舊有 `booking_details`，並重建 `booking_details`
9. 更新其他欄位：備註、聊天開關。
10. 回傳更新後資料 ID。

---

