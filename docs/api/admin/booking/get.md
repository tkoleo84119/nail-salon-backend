## User Story

作為員工，我希望可以查詢單一預約的詳細資料，方便管理。

---

## Endpoint

**GET** `/api/admin/stores/:storeId/bookings/:bookingId`

---

## 說明

- 用於查詢特定預約的詳細資訊。

---

## 權限

- 需要登入才可使用。
- 所有角色皆可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameters

| 參數      | 說明    |
| --------- | ------- |
| storeId   | 門市 ID |
| bookingId | 預約 ID |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "3000000001",
    "customer": {
      "id": "2000000001",
      "name": "小美"
    },
    "stylist": {
      "id": "7000000001",
      "name": "Ariel"
    },
    "timeSlot": {
      "id": "9000000001",
      "workDate": "2025-08-01",
      "startTime": "10:00",
      "endTime": "11:00"
    },
    "isChatEnabled": true,
    "status": "SCHEDULED",
    "note": "這是備註",
    "createdAt": "2025-01-01T00:00:00+08:00",
    "updatedAt": "2025-01-01T00:00:00+08:00",
    "bookingDetails": [
      {
        "id": "9000000001",
        "service": {
          "id": "9000000010",
          "name": "法式美甲",
          "is_addon": false,
        },
        "rawPrice": 1000,
        "price": 800
      },
      {
        "id": "9000000001",
        "service": {
          "id": "9000000010",
          "name": "卸甲",
          "is_addon": true,
        },
        "rawPrice": 300,
        "price": 300
      }
    ],
    "checkout": {
      "id": "9000000001",
      "paymentMethod": "cash",
      "totalAmount": 1300,
      "finalAmount": 1100,
      "paidAmount": 1100,
      "checkoutUser": "admin",
      "coupon": {
        "id": "1000000001",
        "name": "優惠券",
        "code": "TEXT"
      }
    }
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

| 狀態碼 | 錯誤碼  | 常數名稱                | 說明                             |
| ------ | ------- | ----------------------- | -------------------------------- |
| 401    | E1002   | AuthTokenInvalid        | 無效的 accessToken，請重新登入   |
| 401    | E1003   | AuthTokenMissing        | accessToken 缺失，請重新登入     |
| 401    | E1004   | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入 |
| 401    | E1005   | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入 |
| 401    | E1006   | AuthContextMissing      | 未找到使用者認證資訊，請重新登入 |
| 400    | E2002   | ValPathParamMissing     | 路徑參數缺失，請檢查             |
| 400    | E2004   | ValTypeConversionFailed | 參數類型轉換失敗                 |
| 404    | E3BK001 | BookingNotFound         | 預約不存在或已被取消             |
| 500    | E9001   | SysInternalError        | 系統發生錯誤，請稍後再試         |
| 500    | E9002   | SysDatabaseError        | 資料庫操作失敗                   |

---

## 實作與流程

### 資料表

- `bookings`
- `customers`
- `stylists`
- `time_slots`
- `services`
- `booking_details`
- `checkouts`
- `coupons`

---

### Service 邏輯

1. 驗證員工是否有權限查詢該門市。
2. 查詢 `bookings` 表中該筆預約是否存在。
3. 確認該 `booking` 是否隸屬於該門市。
4. 查詢 `booking_details` 表中該筆預約的詳細資訊。
5. 如果是已結帳的預約，則查詢 `checkouts` 表中該筆預約的結帳資訊。
6. 整理回傳資料。

---

## 注意事項

- `createdAt` 與 `updatedAt` 會是標準 Iso 8601 格式。
