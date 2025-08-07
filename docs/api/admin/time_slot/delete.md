## User Story

1. 作為一位美甲師，我希望可以針對一筆 time_slot 進行刪除。
2. 作為一位管理員，我希望可以針對單一美甲師的一筆 time_slot 進行刪除。

---

## Endpoint

**DELETE** `/api/admin/schedules/{scheduleId}/time-slots/{timeSlotId}`

---

## 說明

- 可針對單一時段（time_slot）進行刪除。
- 美甲師僅能刪除自己班表的 time_slot（僅限自己有權限的 store）。
- 管理員可刪除任一美甲師的 time_slot（僅限自己有權限的 store）。
- 已被預約的 time_slot 禁止刪除。

---

## 權限

- 需要登入才可使用。
- `SUPER_ADMIN`、`ADMIN`、`MANAGER` 可刪除任何美甲師的 time_slot（僅限自己有權限的 store）。
- `STYLIST` 僅可刪除自己班表的 time_slot（僅限自己有權限的 store）。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameter

| 參數       | 說明   |
| ---------- | ------ |
| scheduleId | 班表ID |
| timeSlotId | 時段ID |

---

## Response

### 成功 204 No Content

```json
{
  "data": {
    "deleted": "5000000011"
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

| 狀態碼 | 錯誤碼   | 常數名稱                         | 說明                             |
| ------ | -------- | -------------------------------- | -------------------------------- |
| 401    | E1002    | AuthInvalidCredentials           | 無效的 accessToken，請重新登入   |
| 401    | E1003    | AuthTokenMissing                 | accessToken 缺失，請重新登入     |
| 401    | E1004    | AuthTokenFormatError             | accessToken 格式錯誤，請重新登入 |
| 401    | E1005    | AuthStaffFailed                  | 未找到有效的員工資訊，請重新登入 |
| 401    | E1006    | AuthContextMissing               | 未找到使用者認證資訊，請重新登入 |
| 403    | E1010    | AuthPermissionDenied             | 權限不足，無法執行此操作         |
| 400    | E2002    | ValPathParamMissing              | 路徑參數缺失，請檢查             |
| 400    | E2004    | ValTypeConversionFailed          | 參數類型轉換失敗                 |
| 400    | E3TMS004 | TimeSlotAlreadyBookedDoNotDelete | 該時段已被預約，無法刪除         |
| 400    | E3TMS002 | TimeSlotNotBelongToSchedule      | 時段不屬於指定的班表             |
| 404    | E3TMS008 | TimeSlotNotFound                 | 時段不存在或已被刪除             |
| 404    | E3SCH005 | ScheduleNotFound                 | 排班不存在或已被刪除             |
| 404    | E3STY001 | StylistNotFound                  | 美甲師資料不存在                 |
| 500    | E9001    | SysInternalError                 | 系統發生錯誤，請稍後再試         |
| 500    | E9002    | SysDatabaseError                 | 資料庫操作失敗                   |

---

## 資料表

- `time_slots`
- `schedules`
- `stylists`

---

## Service 邏輯

1. 檢查 `timeSlotId` 是否存在。
2. 確認 `timeSlotId` 是否屬於指定 schedule。
3. 確認 `timeSlotId` 是否已被預約。(被預約時，不可刪除)
4. 取得 `schedule` 資訊。
5. 取得 `stylist` 資訊。
6. 判斷身分是否可操作指定 `time_slot`（美甲師僅可刪除自己的 `time_slot`，管理員可刪除任一美甲師 `time_slot`）。
7. 執行刪除。
8. 回傳已刪除 `id`。

---

## 注意事項

- 美甲師僅能針對自己的 `time_slot` 刪除；管理員可針對任一美甲師的 `time_slot` 刪除。
- 僅可操作自己有權限的 `store`。
- 時段一旦被預約（有 `booking` 記錄）禁止刪除。
