## User Story

作為顧客，我希望能夠取得我的單筆預約資訊。

---

## Endpoint

**GET** `/api/bookings/{bookingId}`

---

## 說明

- 回傳對應預約的詳細資訊（包含服務、時段、門市、美甲師、狀態、備註等）。

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
    "mainService": {
      "id": "1000000001",
      "name": "主服務項目名稱"
    },
    "subServices": [
      {
        "id": "1000000002",
        "name": "副服務項目名稱1"
      },
      {
        "id": "1000000003",
        "name": "副服務項目名稱2"
      }
    ],
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

| 狀態碼 | 錯誤碼  | 常數名稱               | 說明                             |
| ------ | ------- | ---------------------- | -------------------------------- |
| 401    | E1002  | AuthTokenInvalid       | 無效的 accessToken，請重新登入   |
| 401    | E1003   | AuthTokenMissing       | accessToken 缺失，請重新登入     |
| 401    | E1004   | AuthTokenFormatError   | accessToken 格式錯誤，請重新登入 |
| 401    | E1006   | AuthContextMissing     | 未找到使用者認證資訊，請重新登入 |
| 401    | E1011   | AuthCustomerFailed     | 未找到有效的顧客資訊，請重新登入 |
| 400    | E2002   | ValPathParamMissing    | 路徑參數缺失，請檢查             |
| 404    | E3BK001 | BookingNotFound        | 預約不存在或已被取消             |
| 500    | E9001   | SysInternalError       | 系統發生錯誤，請稍後再試         |
| 500    | E9002   | SysDatabaseError       | 資料庫操作失敗                   |

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

1. 查詢該 `bookingId` 是否存在 (一定會帶上 `customer_id`)。
2. 回傳對應預約明細。

---

## 注意事項

- 查無資料回傳 404。
