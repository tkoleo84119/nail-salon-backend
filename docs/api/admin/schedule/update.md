## User Story

1. 作為一位美甲師，我希望可以更新我自己的出勤班表（schedules）。
2. 作為一位管理員，我希望可以更新其他美甲師的班表（schedules）。

---

## Endpoint

**PATCH** `/api/admin/store/:storeId/schedules/:scheduleId`

---

## 說明

- 美甲師只能為自己更新班表 (只能更新自己有權限的 `store`)。
- 管理員可為任一美甲師更新班表 (只能更新自己有權限的 `store`)。

---

## 權限

- 需要登入才可使用。
- `SUPER_ADMIN`、`ADMIN`、`MANAGER` 可為任何美甲師更新班表 (只能更新自己有權限的 `store`)。
- `STYLIST` 僅可為自己更新班表 (只能更新自己有權限的 `store`)。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameters

| 參數       | 說明    |
| ---------- | ------- |
| storeId    | 門市 ID |
| scheduleId | 班表 ID |

### Body 範例

```json
{
  "stylistId": "18000000001",
  "workDate": "2024-07-21",
  "note": "早班"
}
```

### 驗證規則

| 欄位      | 必填 | 其他規則            | 說明     |
| --------- | ---- | ------------------- | -------- |
| stylistId | 是   |                     | 美甲師id |
| workDate  | 否   | <li>YYYY-MM-DD 格式 | 班表日期 |
| note      | 否   | <li>最長100字元     | 備註     |

- workDate 與 note 至少要有一個。

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "id": "5000000001",
    "workDate": "2025-08-01",
    "note": "上班全天",
    "timeSlots": [
      {
        "id": "9000000001",
        "startTime": "10:00",
        "endTime": "11:00",
        "isAvailable": true
      },
      {
        "id": "9000000002",
        "startTime": "11:00",
        "endTime": "12:00",
        "isAvailable": false
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

| 狀態碼 | 錯誤碼   | 常數名稱                         | 說明                                                |
| ------ | -------- | -------------------------------- | --------------------------------------------------- |
| 401    | E1002  | AuthTokenInvalid       | 無效的 accessToken，請重新登入                      |
| 401    | E1003    | AuthTokenMissing                 | accessToken 缺失，請重新登入                        |
| 401    | E1004    | AuthTokenFormatError             | accessToken 格式錯誤，請重新登入                    |
| 401    | E1005    | AuthStaffFailed                  | 未找到有效的員工資訊，請重新登入                    |
| 401    | E1006    | AuthContextMissing               | 未找到使用者認證資訊，請重新登入                    |
| 403    | E1010    | AuthPermissionDenied             | 權限不足，無法執行此操作                            |
| 400    | E2002    | ValPathParamMissing              | 路徑參數缺失，請檢查                                |
| 400    | E2003    | ValAllFieldsEmpty                | 至少需要提供一個欄位進行更新                        |
| 400    | E2004    | ValTypeConversionFailed          | 參數類型轉換失敗                                    |
| 400    | E2001    | ValJsonFormat                    | JSON 格式錯誤，請檢查                               |
| 400    | E2020    | ValFieldRequired                 | {field} 為必填欄位                                  |
| 400    | E2024    | ValFieldStringMaxLength          | {field} 長度最多只能有 {param} 個字元               |
| 400    | E2033    | ValFieldDateFormat               | {field} 格式錯誤，請使用正確的日期格式 (YYYY-MM-DD) |
| 400    | E3SCH001 | ScheduleAlreadyBookedDoNotUpdate | 班表已被預約，無法更新                              |
| 400    | E3SCH003 | ScheduleNotBelongToStore         | 部分班表不屬於指定的門市                            |
| 400    | E3SCH004 | ScheduleNotBelongToStylist       | 部分班表不屬於指定的美甲師                          |
| 400    | E3SCH006 | ScheduleAlreadyExists            | 美甲師班表已存在                                    |
| 400    | E3SCH009 | ScheduleDuplicateWorkDateInput   | 輸入的工作日期重複                                  |
| 400    | E3SCH010 | ScheduleCannotCreateBeforeToday  | 不能創建過去的班表                                  |
| 400    | E3SCH011 | ScheduleAlreadyBookedDoNotUpdate | 部分時段已被預約，無法更新                          |
| 404    | E3SCH005 | ScheduleNotFound                 | 排班不存在或已被刪除                                |
| 404    | E3STY001 | StylistNotFound                  | 美甲師資料不存在                                    |
| 500    | E9001    | SysInternalError                 | 系統發生錯誤，請稍後再試                            |
| 500    | E9002    | SysDatabaseError                 | 資料庫操作失敗                                      |

---

## 資料表

- `schedules`
- `stylists`
- `stores`

---

## Service 邏輯

1. 檢查是否有至少一個欄位有值。
2. 檢查 `stylistId` 是否存在。
3. 判斷身分是否可操作指定 stylistId (員工只能更新自己的班表，管理員可更新任一美甲師班表)。
4. 判斷使用者是否有權限操作指定 `storeId`。
5. 驗證 `workDate` 格式是否正確。
6. 檢查新的 `workDate` 是否與已有的班表重複。
7. 檢查 `schedule` 是否存在。
8. 檢查隸屬於 `schedule` 的 `time slot` 是否皆為 `available` (有被預約的時段不能更新)。
9. 更新 `schedules` 資料。
10. 回傳更新結果。

---

## 注意事項

- 員工僅能更新自己的班表；管理員可更新任一美甲師班表。
- 同一天、同店、同美甲師僅能有一筆 schedule。
