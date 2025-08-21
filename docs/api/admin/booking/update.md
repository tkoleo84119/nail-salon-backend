## User Story

作為員工，我希望可以更新某筆顧客預約（Booking），例如更換時段、修改服務、調整備註等，以因應顧客變動需求。

---

## Endpoint

**PATCH** `/api/admin/stores/{storeId}/bookings/{bookingId}`

---

## 說明

- 提供後台管理員更新預約功能。
- 可部分修改指定預約的內容：時段、服務項目、備註、是否啟用聊天。
- 修改時段會檢查是否仍可預約（`is_available=true`）。

---

## 權限

- 需要登入才可使用。
- 所有角色皆可使用。

---

## Request

### Header

- Authorization: Bearer <access_token>
- Content-Type: application/json

### Path Parameters

| 參數      | 說明    |
| --------- | ------- |
| storeId   | 門市 ID |
| bookingId | 預約 ID |

### Body 範例

```json
{
  "stylistId": "9000000001",
  "timeSlotId": "9000000002",
  "mainServiceId": "9000000010",
  "subServiceIds": ["9000000011", "9000000012"],
  "isChatEnabled": false,
  "note": "顧客改約下午並不開啟聊天"
}
```

### 驗證規則

| 欄位          | 必填 | 其他規則      | 說明          |
| ------------- | ---- | ------------- | ------------- |
| stylistId     | 否   |               | 美甲師ID      |
| timeSlotId    | 否   |               | 時段ID        |
| mainServiceId | 否   |               | 主服務項目ID  |
| subServiceIds | 否   | <li>最多10項  | 副服務項目IDs |
| isChatEnabled | 否   | <li>布林值    | 是否要聊天    |
| note          | 否   | <li>最長255字 | 備註說明      |


- stylistId、timeSlotId、mainServiceId、subServiceIds 必須一起傳入
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
| 400    | E2029    | ValFieldBoolean                 | {field} 必須是布林值                         |
| 400    | E3BK002  | BookingStatusNotAllowedToUpdate | 預約狀態不允許更新                           |
| 400    | E3BK007  | BookingUpdateIncomplete         | 預約更新資訊不完整，所有必要資訊必須一起傳入 |
| 400    | E3SER001 | ServiceNotActive                | 服務未啟用                                   |
| 400    | E3SER002 | ServiceNotMainService           | 服務不是主服務                               |
| 400    | E3SER003 | ServiceNotAddon                 | 服務不是附屬服務                             |
| 400    | E3TMS006 | TimeSlotNotEnoughTime           | 時段時間不足                                 |
| 404    | E3BK001  | BookingNotFound                 | 預約不存在或已被取消                         |
| 404    | E3TMS005 | TimeSlotNotFound                | 時段不存在或已被刪除                         |
| 404    | E3SER004 | ServiceNotFound                 | 服務不存在或已被刪除                         |
| 404    | E3STY001 | StylistNotFound                 | 美甲師資料不存在                             |
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

---

## Service 邏輯
1. 驗證是否有傳入任一欄位，若沒有則回傳錯誤
2. 驗證預約是否存在並隸屬於該門市，且預約狀態為 `SCHEDULED`
3. 若有傳入美甲師、時段、服務，則驗證異動後的美甲師、時段、服務。
   1. 驗證門市、美甲師、時段、服務是否存在
   2. 驗證時段是否可用
   3. 驗證服務是否可用
   4. 驗證時段時間是否足夠
   5. 驗證附加服務是否可用
4. 更新預約內容（`bookings`、`booking_details`）。
5. 回傳最新預約資訊。
