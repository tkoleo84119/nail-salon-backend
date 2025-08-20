## User Story

作為顧客，我希望能夠取得某筆排班（Schedule）底下全部仍可以預約的時段（Time Slot），以便選擇預約時間。

---

## Endpoint

**GET** `/api/schedules/{scheduleId}/time-slots`

---

## 說明

- 提供顧客查詢指定排班（schedule）下可預約的時段。
- 僅回傳 `is_available=true` 的時段。
- 每筆資料包含：起始時間、結束時間與該時段的總分鐘數（供前端計算服務是否可容納）。
- 若顧客為黑名單（`is_blacklisted=true`），回傳空陣列。
- 若排班不存在，回傳空陣列。
- 按照 `startTime` 升冪排序。

---

## 權限

- 需要登入才可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameter

| 參數       | 說明    |
| ---------- | ------- |
| scheduleId | 排班 ID |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "timeSlots": [
      {
        "id": "9000000001",
        "startTime": "10:00",
        "endTime": "11:00",
        "durationMinutes": 60
      },
      {
        "id": "9000000002",
        "startTime": "11:30",
        "endTime": "12:30",
        "durationMinutes": 60
      }
    ]
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

| 狀態碼 | 錯誤碼   | 常數名稱                | 說明                             |
| ------ | -------- | ----------------------- | -------------------------------- |
| 401    | E1002  | AuthTokenInvalid       | 無效的 accessToken，請重新登入   |
| 401    | E1003    | AuthTokenMissing        | accessToken 缺失，請重新登入     |
| 401    | E1004    | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入 |
| 401    | E1006    | AuthContextMissing      | 未找到使用者認證資訊，請重新登入 |
| 401    | E1011    | AuthCustomerFailed      | 未找到有效的顧客資訊，請重新登入 |
| 400    | E2002    | ValPathParamMissing     | 路徑參數缺失，請檢查             |
| 400    | E2004    | ValTypeConversionFailed | 參數類型轉換失敗                 |
| 404    | E3SCH005 | ScheduleNotFound        | 排班不存在或已被刪除             |
| 500    | E9001    | SysInternalError        | 系統發生錯誤，請稍後再試         |
| 500    | E9002    | SysDatabaseError        | 資料庫操作失敗                   |

---

## 資料表

- `customers`
- `schedules`
- `time_slots`

---

## Service 邏輯

1. 檢查顧客是否為黑名單（`is_blacklisted=true`），若是則回傳空陣列。
2. 驗證指定排班是否存在。
3. 查詢該排班底下 `is_available=true` 的時段，並計算 `durationMinutes = endTime - startTime`（單位：分鐘）。
4. 回傳每筆時段。

---

## 注意事項

- `durationMinutes` 為方便前端驗證服務長度是否合適。
