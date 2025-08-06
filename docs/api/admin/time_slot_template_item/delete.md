## User Story

作為一位管理員，我希望能刪除時段範本（template）下的一筆 time_slot_template_item，彈性維護排班範本。

---

## Endpoint

**DELETE** `/api/admin/time-slot-templates/{templateId}/items/{itemId}`

---

## 說明

- 可針對指定範本（template）下的特定時段（time_slot_template_item）進行刪除。

---

## 權限

- 需要登入才可使用。
- 僅 `SUPER_ADMIN`、`ADMIN`、`MANAGER` 可建立。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameter

| 參數       | 說明   |
| ---------- | ------ |
| templateId | 範本ID |
| itemId     | 項目ID |

---

## Response

### 成功 204 No Content

```json
{
  "data": {
    "deleted": "6100000003"
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

| 狀態碼 | 錯誤碼   | 常數名稱                     | 說明                             |
| ------ | -------- | ---------------------------- | -------------------------------- |
| 401    | E1002    | AuthInvalidCredentials       | 無效的 accessToken，請重新登入   |
| 401    | E1003    | AuthTokenMissing             | accessToken 缺失，請重新登入     |
| 401    | E1004    | AuthTokenFormatError         | accessToken 格式錯誤，請重新登入 |
| 401    | E1005    | AuthStaffFailed              | 未找到有效的員工資訊，請重新登入 |
| 401    | E1006    | AuthContextMissing           | 未找到使用者認證資訊，請重新登入 |
| 403    | E1010    | AuthPermissionDenied         | 權限不足，無法執行此操作         |
| 400    | E2002    | ValPathParamMissing          | 路徑參數缺失，請檢查             |
| 400    | E2004    | ValTypeConversionFailed      | 參數類型轉換失敗                 |
| 404    | E3TMS010 | TimeSlotTemplateNotFound     | 範本項目不存在或已被刪除         |
| 404    | E3TMS013 | TimeSlotTemplateItemNotFound | 範本項目不存在或已被刪除         |
| 500    | E9001    | SysInternalError             | 系統發生錯誤，請稍後再試         |
| 500    | E9002    | SysDatabaseError             | 資料庫操作失敗                   |

---

## 資料表

- `time_slot_templates`
- `time_slot_template_items`

---

## Service 邏輯

1. 驗證 `templateId` 是否存在。
2. 驗證 `itemId` 是否存在。
3. 刪除 `time_slot_template_items` 資料。
4. 回傳已刪除 id。

---
