## User Story

作為員工，我希望可以幫客戶新增預約（Booking），以便協助客戶安排服務時段與內容。

---

## Endpoint

**POST** `/api/admin/stores/{storeId}/bookings`

---

## 說明

- 提供後台管理員新增預約功能。
- 須指定美甲師、時段、服務項目等必要資訊。

---

## 權限

- 需要登入才可使用。
- 所有角色皆可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameter

| 參數    | 說明    |
| ------- | ------- |
| storeId | 門市 ID |

### Body 範例

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

## 驗證規則

| 欄位          | 必填 | 其他規則        | 說明       |
| ------------- | ---- | --------------- | ---------- |
| customerId    | 是   |                 | 顧客 ID    |
| stylistId     | 是   |                 | 美甲師 ID  |
| timeSlotId    | 是   |                 | 時段 ID    |
| mainServiceId | 是   |                 | 主服務 ID  |
| subServiceIds | 是   | <li>最大10筆    | 子服務 IDs |
| isChatEnabled | 是   |                 | 是否要聊天 |
| note          | 選填 | <li>最大長度255 | 備註       |

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "id": "5000000001"
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
| 401    | E1002    | AuthTokenInvalid           | 無效的 accessToken，請重新登入        |
| 401    | E1003    | AuthTokenMissing           | accessToken 缺失，請重新登入          |
| 401    | E1004    | AuthTokenFormatError       | accessToken 格式錯誤，請重新登入      |
| 401    | E1006    | AuthContextMissing         | 未找到使用者認證資訊，請重新登入      |
| 401    | E1011    | AuthCustomerFailed         | 未找到有效的顧客資訊，請重新登入      |
| 400    | E2001    | ValJSONFormatError         | JSON 格式錯誤，請檢查                 |
| 400    | E2020    | ValFieldRequired           | {field} 為必填項目                    |
| 400    | E2024    | ValFieldStringMaxLength    | {field} 長度最多只能有 {param} 個字元 |
| 400    | E2025    | ValFieldArrayMaxLength     | {field} 最多只能有 {param} 個項目     |
| 400    | E3STO001 | StoreNotActive             | 門市未啟用                            |
| 400    | E3SER001 | ServiceNotActive           | 服務未啟用                            |
| 400    | E3SER002 | ServiceNotMainService      | 服務不是主服務                        |
| 400    | E3SER003 | ServiceNotAddon            | 服務不是附屬服務                      |
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
- `customers`
- `stylists`
- `time_slots`
- `services`
- `stores`

---

## Service 邏輯

1. 驗證門市是否存在
2. 驗證員工是否有權限操作該門市
3. 驗證美甲師、時段、服務是否存在，且該時段可預約。
4. 建立 `bookings` 主檔與對應的 `booking_details`。
5. 標記 `time_slots.is_available = false`。
6. 回傳資料。

---

## 注意事項

- 僅支援單一時段對應一筆預約。
