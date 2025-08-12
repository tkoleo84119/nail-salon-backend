## User Story

作為顧客，我希望可以預約喜歡的美甲師與服務項目，方便安排我的美甲時段。

---

## Endpoint

**POST** `/api/bookings`

---

## 說明

- 提供顧客預約美甲師與服務項目。
- 顧客可指定門市、美甲師、時段與服務項目進行預約。

---

## 權限

- 需要登入才可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Body 範例

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

| 欄位          | 必填 | 其他規則      | 說明          |
| ------------- | ---- | ------------- | ------------- |
| storeId       | 是   | <li>必填      | 預約門市      |
| stylistId     | 是   | <li>必填      | 美甲師ID      |
| timeSlotId    | 是   | <li>必填      | 時段ID        |
| mainServiceId | 是   | <li>必填      | 主服務項目ID  |
| subServiceIds | 否   | <li>最多5項   | 副服務項目IDs |
| isChatEnabled | 否   | <li>布林值    | 是否要聊天    |
| note          | 否   | <li>最長255字 | 備註說明      |

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

| 狀態碼 | 錯誤碼   | 常數名稱                   | 說明                                  |
| ------ | -------- | -------------------------- | ------------------------------------- |
| 401    | E1002    | AuthInvalidCredentials     | 無效的 accessToken，請重新登入        |
| 401    | E1003    | AuthTokenMissing           | accessToken 缺失，請重新登入          |
| 401    | E1004    | AuthTokenFormatError       | accessToken 格式錯誤，請重新登入      |
| 401    | E1006    | AuthContextMissing         | 未找到使用者認證資訊，請重新登入      |
| 401    | E1011    | AuthCustomerFailed         | 未找到有效的顧客資訊，請重新登入      |
| 400    | E2001    | ValJSONFormatError         | JSON 格式錯誤，請檢查                 |
| 400    | E2020    | ValFieldRequired           | {field} 為必填項目                    |
| 400    | E2024    | ValFieldStringMaxLength    | {field} 長度最多只能有 {param} 個字元 |
| 400    | E2025    | ValFieldArrayMaxLength     | {field} 最多只能有 {param} 個項目     |
| 400    | E2029    | ValFieldBoolean            | {field} 必須是布林值                  |
| 400    | E2030    | ValFieldOneof              | {field} 必須是 {param} 其中一個值     |
| 400    | E3STO001 | StoreNotActive             | 門市未啟用                            |
| 400    | E3SER001 | ServiceNotActive           | 服務未啟用                            |
| 400    | E3SER002 | ServiceNotMainService      | 服務不是主服務                        |
| 400    | E3SER003 | ServiceNotAddon            | 服務不是附屬服務                      |
| 400    | E3TMS006 | TimeSlotNotEnoughTime      | 時段時間不足                          |
| 404    | E3STO002 | StoreNotFound              | 門市不存在或已被刪除                  |
| 404    | E3TMS005 | TimeSlotNotFound           | 時段不存在或已被刪除                  |
| 404    | E3SER004 | ServiceNotFound            | 服務不存在或已被刪除                  |
| 404    | E3STY001 | StylistNotFound            | 美甲師資料不存在                      |
| 409    | E3BK006  | BookingTimeSlotUnavailable | 該時段已被預約，請重新選擇            |
| 500    | E9001    | SysInternalError           | 系統發生錯誤，請稍後再試              |
| 500    | E9002    | SysDatabaseError           | 資料庫操作失敗                        |

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
3. 驗證時段時間是否足夠支援服務（主服務+副服務）。
4. 建立預約資料（bookings、booking_details）。
5. 更新時段狀態為不可預約。
6. 回傳預約資訊。

---

## 注意事項

- 僅支援本人預約。
