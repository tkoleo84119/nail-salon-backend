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

#### 錯誤總覽

| 狀態碼 | 錯誤碼 | 說明                                  |
| ------ | ------ | ------------------------------------- |
| 401    | E1002  | 無效的 accessToken，請重新登入        |
| 401    | E1003  | accessToken 缺失，請重新登入          |
| 401    | E1004  | accessToken 格式錯誤，請重新登入      |
| 401    | E1005  | 未找到有效的員工資訊，請重新登入      |
| 401    | E1006  | 未找到使用者認證資訊，請重新登入      |
| 400    | E2002  | 路徑參數缺失，請檢查                  |
| 400    | E2023  | {field} 最小值為 {param}              |
| 400    | E2024  | {field} 長度最多只能有 {param} 個字元 |
| 400    | E2026  | {field} 最大值為 {param}              |
| 500    | E9001  | 系統發生錯誤，請稍後再試              |
| 500    | E9002  | 資料庫操作失敗                        |

#### 400 Bad Request - 參數類型轉換失敗

```json
{
  "error": {
    "code": "E2002",
    "message": "路徑參數缺失，請檢查"
  }
}
```

#### 401 Unauthorized - 未登入/Token失效

```json
{
  "error": {
    "code": "E1002",
    "message": "無效的 accessToken，請重新登入"
  }
}
```

#### 500 Internal Server Error - 系統發生錯誤

```json
{
  "error": {
    "code": "E9001",
    "message": "系統發生錯誤，請稍後再試"
  }
}
```

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
