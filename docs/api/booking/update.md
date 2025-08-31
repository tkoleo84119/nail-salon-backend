## User Story

作為顧客，我希望可以更改自己的預約（包含不同門市、不同美甲師、不同服務、不同附加服務），彈性調整預約內容。

---

## Endpoint

**PATCH** `/api/bookings/{bookingId}`

---

## 說明

- 提供顧客修改自己的預約。
- 可變更門市、美甲師、時段、服務（可含多個）、備註等。
- 異動時會自動檢查新時段與服務項目的可用性與衝突。
- 異動時段會更新舊有時段狀態為可預約。

---

## 權限

- 需要登入才可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameter

| 參數      | 說明   |
| --------- | ------ |
| bookingId | 預約ID |

### Body 範例

```json
{
  "timeSlotId": "3000000020",
  "mainServiceId": "9000000001",
  "subServiceIds": ["9000000002", "9000000003"],
  "isChatEnabled": true,
  "hasChatPermission": true,
  "note": "這次想做奶茶色"
}
```

### 驗證規則

| 欄位              | 必填 | 其他規則      | 說明                   |
| ----------------- | ---- | ------------- | ---------------------- |
| hasChatPermission | 是   |               | 是否擁有line聊天室權限 |
| timeSlotId        | 否   |               | 時段ID                 |
| mainServiceId     | 否   |               | 主服務項目ID           |
| subServiceIds     | 否   | <li>最多5項   | 副服務項目IDs          |
| isChatEnabled     | 否   |               | 是否要聊天             |
| note              | 否   | <li>最長255字 | 備註說明               |

- 欄位皆為選填，至少需要傳入一個欄位
- timeSlotId、mainServiceId、subServiceIds 必須一起傳入

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
    "customerName": "顧客名稱",
    "customerPhone": "顧客電話",
    "date": "2025-08-02",
    "timeSlotId": "3000000001",
    "startTime": "10:00",
    "endTime": "11:00",
    "mainServiceName": "主服務項目名稱",
    "subServiceNames": ["副服務項目名稱1", "副服務項目名稱2"],
    "isChatEnabled": true,
    "note": "這次想做奶茶色",
    "status": "SCHEDULED",
    "createdAt": "2025-01-01T00:00:00+08:00",
    "updatedAt": "2025-01-01T00:00:00+08:00"
  }
}
```

### 錯誤處理

全部 API 皆回傳如下結構，請參考錯誤總覽。

```json
{
  "errors": [
    {
      "code": "EXXXX",
      "message": "錯誤訊息",
      "field": "錯誤欄位名稱"
    }
  ]
}
```

- 欄位說明：
  - errors: 錯誤陣列（支援多筆同時回報）
  - code: 錯誤代碼，唯一對應每種錯誤
  - message: 中文錯誤訊息（可參照錯誤總覽）
  - field: 參數欄位名稱（僅部分驗證錯誤有）

| 狀態碼 | 錯誤碼   | 常數名稱                        | 說明                                         |
| ------ | -------- | ------------------------------- | -------------------------------------------- |
| 401    | E1002    | AuthTokenInvalid                | 無效的 accessToken，請重新登入               |
| 401    | E1003    | AuthTokenMissing                | accessToken 缺失，請重新登入                 |
| 401    | E1004    | AuthTokenFormatError            | accessToken 格式錯誤，請重新登入             |
| 401    | E1006    | AuthContextMissing              | 未找到使用者認證資訊，請重新登入             |
| 401    | E1011    | AuthCustomerFailed              | 未找到有效的顧客資訊，請重新登入             |
| 403    | E1010    | AuthPermissionDenied            | 權限不足，無法執行此操作                     |
| 400    | E2001    | ValJSONFormatError              | JSON 格式錯誤，請檢查                        |
| 400    | E2002    | ValPathParamMissing             | 路徑參數缺失，請檢查                         |
| 400    | E2003    | ValAllFieldsEmpty               | 至少需要提供一個欄位進行更新                 |
| 400    | E2004    | ValTypeConversionFailed         | 參數類型轉換失敗                             |
| 400    | E2024    | ValFieldStringMaxLength         | {field} 長度最多只能有 {param} 個字元        |
| 400    | E2025    | ValFieldArrayMaxLength          | {field} 最多只能有 {param} 個項目            |
| 400    | E3BK002  | BookingStatusNotAllowedToUpdate | 預約狀態不允許更新                           |
| 400    | E3BK007  | BookingUpdateIncomplete         | 預約更新資訊不完整，所有必要資訊必須一起傳入 |
| 400    | E3SER001 | ServiceNotActive                | 服務未啟用                                   |
| 400    | E3SER002 | ServiceNotMainService           | 服務不是主服務                               |
| 400    | E3SER003 | ServiceNotAddon                 | 服務不是附屬服務                             |
| 400    | E3TMS006 | TimeSlotNotEnoughTime           | 時段時間不足                                 |
| 404    | E3BK001  | BookingNotFound                 | 預約不存在或已被取消                         |
| 404    | E3TMS005 | TimeSlotNotFound                | 時段不存在或已被刪除                         |
| 404    | E3SER004 | ServiceNotFound                 | 服務不存在或已被刪除                         |
| 409    | E3BK006  | BookingTimeSlotUnavailable      | 該時段已被預約，請重新選擇                   |
| 500    | E9001    | SysInternalError                | 系統發生錯誤，請稍後再試                     |
| 500    | E9002    | SysDatabaseError                | 資料庫操作失敗                               |

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

1. 驗證預約是否存在且屬於本人，並且預約狀態為 `SCHEDULED`。
2. 若有傳入時段、服務，則驗證異動後的時段、服務。
   1. 驗證時段、服務是否存在
   2. 驗證時段是否可用
   3. 驗證服務是否可用
   4. 驗證時段時間是否足夠
   5. 驗證附加服務是否可用
3. 更新預約內容（`bookings`、`booking_details`）。
4. 若異動了時段，則更新舊有時段狀態為可預約。
5. 若異動了時段，則更新新時段狀態為不可預約。
6. 若顧客沒有聊天室權限 (代表前端沒辦法發送訊息給顧客)，且異動了時段，則後端協助發送預約通知到 LINE。
7. 回傳最新預約資訊。

---

## 注意事項

- 僅支援本人預約內容異動。
- 異動時需重新檢查時段、服務是否可用。
