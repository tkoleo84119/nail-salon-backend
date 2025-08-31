## User Story

作為顧客，我希望可以取消自己的預約，並傳入取消原因，讓店家知悉我的狀況。

---

## Endpoint

**PATCH** `/api/bookings/{bookingId}/cancel`

---

## 說明

- 提供顧客取消自己的預約。
- 可傳入取消原因（文字），供後台記錄。
- 預約狀態將變更為 CANCELLED。

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
  "hasChatPermission": true,
  "cancelReason": "臨時有事無法前往，抱歉！"
}
```

### 驗證規則

| 欄位              | 必填 | 其他規則      | 說明                   |
| ----------------- | ---- | ------------- | ---------------------- |
| hasChatPermission | 是   |               | 是否擁有line聊天室權限 |
| cancelReason      | 否   | <li>最長255字 | 取消原因               |

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
    "status": "CANCELLED",
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

| 狀態碼 | 錯誤碼   | 常數名稱                   | 說明                                  |
| ------ | -------- | -------------------------- | ------------------------------------- |
| 401    | E1002    | AuthTokenInvalid           | 無效的 accessToken，請重新登入        |
| 401    | E1003    | AuthTokenMissing           | accessToken 缺失，請重新登入          |
| 401    | E1004    | AuthTokenFormatError       | accessToken 格式錯誤，請重新登入      |
| 401    | E1006    | AuthContextMissing         | 未找到使用者認證資訊，請重新登入      |
| 401    | E1011    | AuthCustomerFailed         | 未找到有效的顧客資訊，請重新登入      |
| 403    | E1010    | AuthPermissionDenied       | 權限不足，無法執行此操作              |
| 400    | E2001    | ValJSONFormatError         | JSON 格式錯誤，請檢查                 |
| 400    | E2020    | ValFieldRequired           | {field} 為必填項目                    |
| 400    | E2024    | ValFieldStringMaxLength    | {field} 長度最多只能有 {param} 個字元 |
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
- `time_slots`

---

## Service 邏輯

1. 驗證預約是否存在且屬於本人，且狀態為 `SCHEDULED`。
2. 記錄取消原因，變更狀態為 `CANCELLED`。
3. 將舊時段狀態更新為可預約。
4. 若顧客沒有聊天室權限 (代表前端沒辦法發送訊息給顧客)，則後端協助發送預約取消通知到 LINE。
5. 回傳結果。

---

## 注意事項

- 僅支援本人預約取消。
- 取消時若有傳入原因，則記錄取消原因。
- 狀態不可重複取消。

