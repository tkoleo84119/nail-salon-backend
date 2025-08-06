## User Story

1. 作為一位美甲師，我希望可以針對自己的班表（schedule）新增一個時段（time_slot）。
2. 作為一位管理員，我希望可以針對單一美甲師的班表（schedule）新增一個時段（time_slot）。

---

## Endpoint

**POST** `/api/admin/schedules/{scheduleId}/time-slots`

---

## 說明

- 可針對單一班表（schedule）新增一個時段（time_slot）。
- 美甲師只能針對自己的班表新增時段（僅限自己有權限的 store）。
- 管理員可針對任一美甲師的班表新增時段（僅限自己有權限的 store）。
- 新增時段時，需檢查時段是否與既有時段重疊。

---

## 權限

- 需要登入才可使用。
- `SUPER_ADMIN`、`ADMIN`、`MANAGER` 可為任何美甲師的班表新增時段 (僅限自己有權限的 store)。
- `STYLIST` 僅可為自己的班表新增時段 (僅限自己有權限的 store)。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>


### Path Parameter

| 參數       | 說明   |
| ---------- | ------ |
| scheduleId | 班表ID |

### Body

```json
{
  "startTime": "14:00",
  "endTime": "16:00"
}
```

### 驗證規則
| 欄位      | 必填 | 其他規則       | 說明     |
| --------- | ---- | -------------- | -------- |
| startTime | 是   | <li>HH:mm 格式 | 起始時間 |
| endTime   | 是   | <li>HH:mm 格式 | 結束時間 |

---

## Response

### 成功 201 Created

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

| 狀態碼 | 錯誤碼   | 常數名稱                        | 說明                                           |
| ------ | -------- | ------------------------------- | ---------------------------------------------- |
| 401    | E1002    | AuthInvalidCredentials          | 無效的 accessToken，請重新登入                 |
| 401    | E1003    | AuthTokenMissing                | accessToken 缺失，請重新登入                   |
| 401    | E1004    | AuthTokenFormatError            | accessToken 格式錯誤，請重新登入               |
| 401    | E1005    | AuthStaffFailed                 | 未找到有效的員工資訊，請重新登入               |
| 401    | E1006    | AuthContextMissing              | 未找到使用者認證資訊，請重新登入               |
| 403    | E1010    | AuthPermissionDenied            | 權限不足，無法執行此操作                       |
| 400    | E2002    | ValPathParamMissing             | 路徑參數缺失，請檢查                           |
| 400    | E2004    | ValTypeConversionFailed         | 參數類型轉換失敗                               |
| 400    | E2034    | ValFieldTimeFormat              | {field} 格式錯誤，請使用正確的時間格式 (HH:mm) |
| 400    | E3SCH010 | ScheduleCannotCreateBeforeToday | 不能創建過去的班表                             |
| 404    | E3SCH005 | ScheduleNotFound                | 排班不存在或已被刪除                           |
| 404    | E3STY001 | StylistNotFound                 | 美甲師資料不存在                               |
| 404    | E3STO002 | StoreNotFound                   | 門市不存在或已被刪除                           |
| 400    | E3STO001 | StoreNotActive                  | 門市未啟用                                     |
| 409    | E3TMS011 | TimeSlotConflict                | 時段時間區段重疊                               |
| 400    | E3TMS012 | TimeSlotEndBeforeStart          | 結束時間必須在開始時間之後                     |
| 500    | E9001    | SysInternalError                | 系統發生錯誤，請稍後再試                       |
| 500    | E9002    | SysDatabaseError                | 資料庫操作失敗                                 |

---

## 資料表

- `schedules`
- `time_slots`
- `stylists`
- `stores`

---

## Service 邏輯

1. 檢查 startTime/endTime 格式。
2. 確認 startTime 必須在 endTime 之前。
3. 檢查 `scheduleId` 是否存在。
4. 檢查 schedule 日期不早於今天 (不能創建過去的班表)。
5. 檢查 schedule 所屬的 stylist/store 是否存在。
6. 判斷身分是否可操作指定 schedule（員工只能新增自己的班表，管理員可新增任一美甲師班表）。
7. 檢查是否有權限操作該 store。
8. 檢查時間區間是否與該 schedule 下既有 time_slots 重疊。
9. 建立新的 time_slot。
10. 回傳新增結果。

---

## 注意事項

- 員工僅能針對自己的班表新增時段；管理員可針對任一美甲師的班表新增時段。
- 不可新增重疊時段。
- 僅可操作自己有權限的 store。
- 不可新增過去的時段。
