## User Story

作為員工，我希望可以查詢特定的排班資料，並一併取得排班底下的時段（Time Slots），以便確認排班內容。

---

## Endpoint

**GET** `/api/admin/stores/{storeId}/schedules/{scheduleId}`

---

## 說明

- 回傳指定排班資料與對應的時段資訊。
- 時段依 `start_time` 排序。

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

| 參數       | 說明    |
| ---------- | ------- |
| storeId    | 門市 ID |
| scheduleId | 排班 ID |

---

## Response

### 成功 200 OK

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

| 狀態碼 | 錯誤碼   | 常數名稱                | 說明                             |
| ------ | -------- | ----------------------- | -------------------------------- |
| 401    | E1002    | AuthInvalidCredentials  | 無效的 accessToken，請重新登入   |
| 401    | E1003    | AuthTokenMissing        | accessToken 缺失，請重新登入     |
| 401    | E1004    | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入 |
| 401    | E1005    | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入 |
| 401    | E1006    | AuthContextMissing      | 未找到使用者認證資訊，請重新登入 |
| 403    | E1010    | AuthPermissionDenied    | 權限不足，無法執行此操作         |
| 400    | E2002    | ValPathParamMissing     | 路徑參數缺失，請檢查             |
| 400    | E2004    | ValTypeConversionFailed | 參數類型轉換失敗                 |
| 404    | E3SCH005 | ScheduleNotFound        | 排班不存在或已被刪除             |
| 500    | E9001    | SysInternalError        | 系統發生錯誤，請稍後再試         |
| 500    | E9002    | SysDatabaseError        | 資料庫操作失敗                   |

---

## 資料表

- `stores`
- `schedules`
- `stylists`
- `time_slots`

---

## Service 邏輯

1. 驗證員工是否有權限存取該門市。
2. 查詢指定 `scheduleId` 同時 JOIN `time_slots` 表取得時段資料。
3. 回傳合併資訊。

---

## 注意事項

- 一律回傳所有時段。
