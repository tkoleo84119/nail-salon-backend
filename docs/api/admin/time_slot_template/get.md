## User Story

作為員工，我希望可以查詢特定的時段模板資料，並一併取得該模板所包含的所有時段項目，以便確認內容或套用排班。

---

## Endpoint

**GET** `/api/admin/time-slot-templates/{templateId}`

---

## 說明

- 回傳模板主資料與所有時間項目（Time Slot Template Items）。

---

## 權限

- 需要登入才可使用。
- 所有角色皆可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameter

| 參數       | 說明    |
| ---------- | ------- |
| templateId | 模板 ID |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "1000000001",
    "name": "早班模板",
    "note": "適用09:00開工",
    "updater": "1000000001",
    "createdAt": "2025-01-01T00:00:00+08:00",
    "updatedAt": "2025-01-01T00:00:00+08:00",
    "items": [
      {
        "id": "1100000001",
        "startTime": "09:00",
        "endTime": "10:00"
      },
      {
        "id": "1100000002",
        "startTime": "10:00",
        "endTime": "11:00"
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

| 狀態碼 | 錯誤碼   | 常數名稱                 | 說明                             |
| ------ | -------- | ------------------------ | -------------------------------- |
| 401    | E1002  | AuthTokenInvalid       | 無效的 accessToken，請重新登入   |
| 401    | E1003    | AuthTokenMissing         | accessToken 缺失，請重新登入     |
| 401    | E1004    | AuthTokenFormatError     | accessToken 格式錯誤，請重新登入 |
| 401    | E1005    | AuthStaffFailed          | 未找到有效的員工資訊，請重新登入 |
| 401    | E1006    | AuthContextMissing       | 未找到使用者認證資訊，請重新登入 |
| 403    | E1010    | AuthPermissionDenied     | 權限不足，無法執行此操作         |
| 400    | E2002    | ValPathParamMissing      | 路徑參數缺失，請檢查             |
| 400    | E2004    | ValTypeConversionFailed  | 參數類型轉換失敗                 |
| 404    | E3TMS010 | TimeSlotTemplateNotFound | 範本項目不存在或已被刪除         |
| 500    | E9001    | SysInternalError         | 系統發生錯誤，請稍後再試         |
| 500    | E9002    | SysDatabaseError         | 資料庫操作失敗                   |

---

## 資料表

- `time_slot_templates`
- `time_slot_template_items`

---

## Service 邏輯

1. 查詢 `time_slot_templates` 與 `time_slot_template_items` 資料。
2. 回傳主檔與時段項目合併結果。

---

## 注意事項

- 每筆 item 僅含時間區段
