## User Story

1. 作為一位美甲師，我希望可以針對一筆 time_slot 進行更新（時間與是否可被預約）。
2. 作為一位管理員，我希望可以針對單一美甲師的一筆 time_slot 進行更新（時間與是否可被預約）。

---

## Endpoint

**PATCH** `/api/admin/schedules/{scheduleId}/time-slots/{timeSlotId}`

---

## 說明

- 可針對單一時段（time_slot）進行更新，包括起訖時間、是否可被預約。
- 美甲師僅能編輯自己班表的 time_slot（僅限自己有權限的 store）。
- 管理員可編輯任一美甲師的 time_slot（僅限自己有權限的 store）。
- 被預約時段不可更新。
- 要更新時段必須同時傳入 startTime/endTime 兩個欄位，不可單獨傳入。
- 更新時需檢查時間區段是否與同一 schedule 下其他 time_slots 重疊。

---

## 權限

- 需要登入才可使用。
- `SUPER_ADMIN`、`ADMIN`、`MANAGER` 可更新任何美甲師的 time_slot（僅限自己有權限的 store）。
- `STYLIST` 僅可更新自己班表的 time_slot（僅限自己有權限的 store）。

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

### Body 範例

```json
{
  "startTime": "14:00",
  "endTime": "16:00",
  "isAvailable": true
}
```

### 驗證規則

| 欄位        | 必填 | 其他規則       | 說明         |
| ----------- | ---- | -------------- | ------------ |
| startTime   | 否   | <li>HH:mm 格式 | 起始時間     |
| endTime     | 否   | <li>HH:mm 格式 | 結束時間     |
| isAvailable | 否   | <li>布林值     | 是否可被預約 |

- 至少需有一項欄位出現。
- 不可單獨傳入 startTime/endTime，需同時傳入。

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "5000000011",
    "scheduleId": "4000000001",
    "startTime": "14:00",
    "endTime": "16:00",
    "isAvailable": true
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

| 狀態碼 | 錯誤碼   | 常數名稱                         | 說明                                           |
| ------ | -------- | -------------------------------- | ---------------------------------------------- |
| 401    | E1002    | AuthInvalidCredentials           | 無效的 accessToken，請重新登入                 |
| 401    | E1003    | AuthTokenMissing                 | accessToken 缺失，請重新登入                   |
| 401    | E1004    | AuthTokenFormatError             | accessToken 格式錯誤，請重新登入               |
| 401    | E1005    | AuthStaffFailed                  | 未找到有效的員工資訊，請重新登入               |
| 401    | E1006    | AuthContextMissing               | 未找到使用者認證資訊，請重新登入               |
| 403    | E1010    | AuthPermissionDenied             | 權限不足，無法執行此操作                       |
| 400    | E2002    | ValPathParamMissing              | 路徑參數缺失，請檢查                           |
| 400    | E2004    | ValTypeConversionFailed          | 參數類型轉換失敗                               |
| 400    | E2001    | ValJsonFormat                    | JSON 格式錯誤，請檢查                          |
| 400    | E2003    | ValAllFieldsEmpty                | 至少需要提供一個欄位進行更新                   |
| 400    | E2034    | ValFieldTimeFormat               | {field} 格式錯誤，請使用正確的時間格式 (HH:mm) |
| 400    | E3TMS001 | TimeSlotCannotUpdateSeparately   | 時段起始時間和結束時間必須同時傳入             |
| 400    | E3TMS012 | TimeSlotEndBeforeStart           | 結束時間必須在開始時間之後                     |
| 404    | E3TMS008 | TimeSlotNotFound                 | 時段不存在或已被刪除                           |
| 400    | E3TMS002 | TimeSlotNotBelongToSchedule      | 時段不屬於指定的班表                           |
| 400    | E3TMS004 | TimeSlotAlreadyBookedDoNotUpdate | 時段已被預約，無法更新                         |
| 404    | E3SCH005 | ScheduleNotFound                 | 排班不存在或已被刪除                           |
| 404    | E3STY001 | StylistNotFound                  | 美甲師資料不存在                               |
| 404    | E3STO002 | StoreNotFound                    | 門市不存在或已被刪除                           |
| 400    | E3STO001 | StoreNotActive                   | 門市未啟用                                     |
| 409    | E3TMS011 | TimeSlotConflict                 | 時段時間區段重疊                               |
| 500    | E9001    | SysInternalError                 | 系統發生錯誤，請稍後再試                       |
| 500    | E9002    | SysDatabaseError                 | 資料庫操作失敗                                 |

---

## 資料表

- `time_slots`
- `schedules`
- `stylists`
- `stores`

---

## Service 邏輯

1. 驗證至少一個欄位有更新。
2. 若 startTime/endTime 有傳入，另外一個欄位必須一起傳入。
3. 檢查 `timeSlotId` 是否存在。
4. 判斷 time_slot 是否屬於指定 schedule。
5. 檢查 `timeSlotId` 是否已被預約。(被預約時，不可變更任何欄位)
6. 檢查 `scheduleId` 是否存在。
7. 取得 stylist 資訊。
8. 判斷身分是否可操作指定 time_slot（員工僅可編輯自己的 time_slot，管理員可編輯任一美甲師 time_slot）。
9. 檢查是否有權限操作該 store。
10. 若有更新時間，檢查是否時間相關邏輯
    1.  startTime / endTime 格式是否正確。
    2.  startTime 必須在 endTime 之前。
    3.  startTime / endTime 是否與 schedule 下其他 time_slots 重疊。
11. 更新 time_slot。
12. 回傳更新結果。

---

## 注意事項

- 員工僅能針對自己的 time_slot 編輯；管理員可針對任一美甲師的 time_slot 編輯。
- 僅可操作自己有權限的 store。
- 若只更新 isAvailable，則時間不檢查重疊；若有更動時間則必須比對同 schedule 下其他 time_slots 是否重疊。
- 被預約時段不可更新。
