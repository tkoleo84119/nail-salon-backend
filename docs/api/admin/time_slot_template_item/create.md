## User Story

作為一位管理員，我希望能在時段範本（template）下新增一筆時段（time_slot_template_item），靈活調整範本內容。

---

## Endpoint

**POST** `/api/admin/time-slot-templates/{templateId}/items`

---

## 說明

- 可針對指定範本（template）新增一筆時段（time_slot_template_item）。

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

### Body 範例

```json
{
  "startTime": "10:00",
  "endTime": "13:00"
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
    "id": "6100000003",
    "templateId": "6000000011",
    "startTime": "10:00",
    "endTime": "13:00"
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

| 狀態碼 | 錯誤碼   | 常數名稱                 | 說明                                           |
| ------ | -------- | ------------------------ | ---------------------------------------------- |
| 401    | E1002  | AuthTokenInvalid       | 無效的 accessToken，請重新登入                 |
| 401    | E1003    | AuthTokenMissing         | accessToken 缺失，請重新登入                   |
| 401    | E1004    | AuthTokenFormatError     | accessToken 格式錯誤，請重新登入               |
| 401    | E1005    | AuthStaffFailed          | 未找到有效的員工資訊，請重新登入               |
| 401    | E1006    | AuthContextMissing       | 未找到使用者認證資訊，請重新登入               |
| 403    | E1010    | AuthPermissionDenied     | 權限不足，無法執行此操作                       |
| 400    | E2002    | ValPathParamMissing      | 路徑參數缺失，請檢查                           |
| 400    | E2004    | ValTypeConversionFailed  | 參數類型轉換失敗                               |
| 400    | E2001    | ValJsonFormat            | JSON 格式錯誤，請檢查                          |
| 400    | E2020    | ValFieldRequired         | {field} 為必填項目                             |
| 400    | E2034    | ValFieldTimeFormat       | {field} 格式錯誤，請使用正確的時間格式 (HH:mm) |
| 400    | E3TMS011 | TimeSlotConflict         | 時段時間區段重疊                               |
| 400    | E3TMS012 | TimeSlotEndBeforeStart   | 結束時間必須在開始時間之後                     |
| 404    | E3TMS010 | TimeSlotTemplateNotFound | 範本項目不存在或已被刪除                       |
| 500    | E9001    | SysInternalError         | 系統發生錯誤，請稍後再試                       |
| 500    | E9002    | SysDatabaseError         | 資料庫操作失敗                                 |

---

## 資料表

- `time_slot_templates`
- `time_slot_template_items`

---

## Service 邏輯

1. 驗證 `startTime`/`endTime` 格式。
2. 驗證 `templateId` 是否存在。
3. 驗證 `startTime` 必須在 `endTime` 之前。
4. 驗證新時段是否與範本內其他時段重疊。
5. 建立 `time_slot_template_item` 資料。
6. 回傳建立結果。

---

## 注意事項

- 時段不得重疊。
