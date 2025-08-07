## User Story

1. 作為一位美甲師，我希望可以安排我自己的出勤班表（schedules），每筆班表可包含多個 time_slot。
2. 作為一位管理員，我希望可以安排其他美甲師的班表（schedules），每筆班表可包含多個 time_slot。

---

## Endpoint

**POST** `/api/admin/store/:storeId/schedules/bulk`

---

## 說明

- 一次只能針對同一位美甲師、同一家門市，新增多日班表（schedules），每筆班表對應一個日期，可包含多個時段（time_slots）。
- 美甲師只能為自己建立班表 (只能建立自己有權限的 `store`)。
- 管理員可為任一美甲師建立班表 (只能建立自己有權限的 `store`)。

---

## 權限

- 需要登入才可使用。
- `SUPER_ADMIN`、`ADMIN`、`MANAGER` 可為任何美甲師建立班表 (只能建立自己有權限的 `store`)。
- `STYLIST` 僅可為自己建立班表 (只能建立自己有權限的 `store`)。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameters

| 參數    | 說明    |
| ------- | ------- |
| storeId | 門市 ID |

### Body 範例

```json
{
  "stylistId": "18000000001",
  "schedules": [
    {
      "workDate": "2024-07-21",
      "note": "早班",
      "timeSlots": [
        { "startTime": "09:00", "endTime": "12:00" },
        { "startTime": "13:00", "endTime": "18:00" }
      ]
    },
    {
      "workDate": "2024-07-22",
      "timeSlots": [
        { "startTime": "09:00", "endTime": "12:00" },
        { "startTime": "13:00", "endTime": "18:00" }
      ]
    }
  ]
}
```

### 驗證規則

| 欄位                          | 必填 | 其他規則                | 說明         |
| ----------------------------- | ---- | ----------------------- | ------------ |
| stylistId                     | 是   |                         | 美甲師id     |
| schedules                     | 是   | <li>最小1筆<li>最大31筆 | 多日班表     |
| schedules.workDate            | 是   | <li>YYYY-MM-DD 格式     | 班表日期     |
| schedules.note                | 否   | <li>最長100字元         | 備註         |
| schedules.timeSlots           | 是   | <li>最小1筆<li>最大20筆 | 當日多個時段 |
| schedules.timeSlots.startTime | 是   | <li>HH:mm 格式          | 起始時間     |
| schedules.timeSlots.endTime   | 是   | <li>HH:mm 格式          | 結束時間     |

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "schedules": [
      {
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

| 狀態碼 | 錯誤碼   | 常數名稱                        | 說明                                                |
| ------ | -------- | ------------------------------- | --------------------------------------------------- |
| 401    | E1002    | AuthInvalidCredentials          | 無效的 accessToken，請重新登入                      |
| 401    | E1003    | AuthTokenMissing                | accessToken 缺失，請重新登入                        |
| 401    | E1004    | AuthTokenFormatError            | accessToken 格式錯誤，請重新登入                    |
| 401    | E1005    | AuthStaffFailed                 | 未找到有效的員工資訊，請重新登入                    |
| 401    | E1006    | AuthContextMissing              | 未找到使用者認證資訊，請重新登入                    |
| 403    | E1010    | AuthPermissionDenied            | 權限不足，無法執行此操作                            |
| 400    | E2001    | ValJsonFormat                   | JSON 格式錯誤，請檢查                               |
| 400    | E2002    | ValPathParamMissing             | 路徑參數缺失，請檢查                                |
| 400    | E2004    | ValTypeConversionFailed         | 參數類型轉換失敗                                    |
| 400    | E2034    | ValFieldTimeFormat              | {field} 格式錯誤，請使用正確的時間格式 (HH:mm)      |
| 400    | E2033    | ValFieldDateFormat              | {field} 格式錯誤，請使用正確的日期格式 (YYYY-MM-DD) |
| 400    | E3SCH006 | ScheduleAlreadyExists           | 美甲師班表已存在                                    |
| 400    | E3SCH010 | ScheduleCannotCreateBeforeToday | 不能創建過去的班表                                  |
| 400    | E3SCH009 | ScheduleDuplicateWorkDateInput  | 輸入的工作日期重複                                  |
| 409    | E3TMS011 | TimeSlotConflict                | 時段時間區段重疊                                    |
| 400    | E3TMS012 | TimeSlotEndBeforeStart          | 結束時間必須在開始時間之後                          |
| 404    | E3STY001 | StylistNotFound                 | 美甲師資料不存在                                    |
| 500    | E9001    | SysInternalError                | 系統發生錯誤，請稍後再試                            |
| 500    | E9002    | SysDatabaseError                | 資料庫操作失敗                                      |

---

## 資料表

- `schedules`
- `time_slots`
- `stylists`
- `stores`

---

## Service 邏輯

1. 檢查 `stylistId` 是否存在。
2. 判斷身分是否可操作指定 stylistId (員工只能建立自己的班表，管理員可建立任一美甲師班表)。
3. 判斷是否有權限操作指定 `storeId`。
4. 驗證每筆 `schedule` 的 `workDate`、`timeSlots`
   - 驗證 `workDate` 格式是否正確。
   - 不可創建過去的班表。
   - 不可傳入相同的 `workDate`。
   - 驗證 `timeSlots` 的 `startTime`、`endTime` 格式是否正確。
   - 驗證 `timeSlots` 的 `startTime` 必須在 `endTime` 之前。
   - 驗證 `timeSlots` 的 `startTime`、`endTime` 不得重疊。
5. 檢查同一天同店同美甲師是否已有班表（不可重複排班）。
6. 新增 `schedules` 資料。
7. 批次建立對應的多筆 `time_slots`。
8. 回傳新增結果。

---

## 注意事項

- 員工僅能建立自己的班表；管理員可建立任一美甲師班表。
- 同一天、同店、同美甲師僅能有一筆 schedule。
- 每個 schedule 需至少一筆 time_slot，且時間區段不得重疊。
