## User Story

作為員工，我希望可以查詢所有時段模板（Time Slot Templates），並可依名稱與分頁條件篩選，以便套用至排班或複製使用。

---

## Endpoint

**GET** `/api/admin/time-slot-templates`

---

## 說明

- 支援基本查詢條件。
- 支援分頁（limit、offset）。
- 支援排序（sort）。

---

## 權限

- 需要登入才可使用。
- 所有角色皆可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Query Parameters

| 參數   | 型別   | 必填 | 預許值    | 說明                                             |
| ------ | ------ | ---- | --------- | ------------------------------------------------ |
| name   | string | 否   |           | 模糊查詢模板名稱                                 |
| limit  | int    | 否   | 20        | 單頁筆數                                         |
| offset | int    | 否   | 0         | 起始筆數                                         |
| sort   | string | 否   | updatedAt | 排序欄位 (可以逗號串接，有 `-` 表示 `DESC` 排序) |

### 驗證規則

| 欄位   | 必填 | 其他規則                                       |
| ------ | ---- | ---------------------------------------------- |
| name   | 否   | 最大長度100字元                                |
| limit  | 否   | 最小值1<li>最大值100                           |
| offset | 否   | 最小值0<li>最大值1000000                       |
| sort   | 否   | 可以為 createdAt, updatedAt, name (其餘會忽略) |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "total": 2,
    "items": [
      {
        "id": "1000000001",
        "name": "早班模板",
        "note": "適用09:00開工",
        "updater": "1000000001",
        "createdAt": "2025-01-01T00:00:00+08:00",
        "updatedAt": "2025-01-01T00:00:00+08:00"
      },
      {
        "id": "1000000002",
        "name": "午班模板",
        "note": "含午休",
        "updater": "1000000002",
        "createdAt": "2025-01-01T00:00:00+08:00",
        "updatedAt": "2025-01-01T00:00:00+08:00"
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

| 狀態碼 | 錯誤碼 | 常數名稱                | 說明                             |
| ------ | ------ | ----------------------- | -------------------------------- |
| 401    | E1002  | AuthTokenInvalid       | 無效的 accessToken，請重新登入   |
| 401    | E1003  | AuthTokenMissing        | accessToken 缺失，請重新登入     |
| 401    | E1004  | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入 |
| 401    | E1005  | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入 |
| 401    | E1006  | AuthContextMissing      | 未找到使用者認證資訊，請重新登入 |
| 403    | E1010  | AuthPermissionDenied    | 權限不足，無法執行此操作         |
| 400    | E2002  | ValPathParamMissing     | 路徑參數缺失，請檢查             |
| 400    | E2004  | ValTypeConversionFailed | 參數類型轉換失敗                 |
| 400    | E2023  | ValFieldMin             | {field} 最小值為 {param}         |
| 400    | E2024  | ValFieldMaxLength       | {field} 長度最多只能有 {param}   |
| 400    | E2026  | ValFieldMax             | {field} 最大值為 {param}         |
| 500    | E9001  | SysInternalError        | 系統發生錯誤，請稍後再試         |
| 500    | E9002  | SysDatabaseError        | 資料庫操作失敗                   |

---

## 資料表

- `time_slot_templates`

---

## Service 邏輯

1. 根據 `name` 模糊查詢 `time_slot_templates` 表的 `name` 欄位。
2. 加入 `limit` / `offset` 處理分頁。
3. 加入 `sort` 處理排序。
4. 回傳總筆數與模板清單。

---

## 注意事項

- 僅回傳模板主資料（不包含子項 item 時段），若需明細請查詢單筆 API (GET `/api/admin/time-slot-templates/{templateId}`)。
