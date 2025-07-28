## User Story

作為員工，我希望可以查詢所有時段模板（Time Slot Templates），並可依名稱與分頁條件篩選，以便套用至排班或複製使用。

---

## Endpoint

**GET** `/api/admin/time-slot-templates`

---

## 說明

- 所有登入員工皆可查詢。
- 用於查詢系統中建立的時段模板列表。
- 可依名稱 `name` 模糊搜尋。
- 支援分頁（limit、offset）。

---

## 權限

- 任一已登入員工皆可使用（JWT 驗證）。

---

## Request

### Header

Authorization: Bearer <access_token>

### Query Parameters

| 參數   | 型別   | 必填 | 預設值 | 說明             |
| ------ | ------ | ---- | ------ | ---------------- |
| name   | string | 否   |        | 模糊查詢模板名稱 |
| limit  | int    | 否   | 20     | 單頁筆數         |
| offset | int    | 否   | 0      | 起始筆數         |

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
        "createdAt": "2025-06-01T08:00:00Z",
        "updatedAt": "2025-06-20T08:00:00Z"
      },
      {
        "id": "1000000002",
        "name": "午班模板",
        "note": "含午休",
        "updater": "1000000002",
        "createdAt": "2025-06-02T08:00:00Z",
        "updatedAt": "2025-06-21T08:00:00Z"
      }
    ]
  }
}
```

### 失敗

#### 401 Unauthorized - 未登入

```json
{
  "message": "無效的 accessToken"
}
```

#### 500 Internal Server Error

```json
{
  "message": "系統發生錯誤，請稍後再試"
}
```

---

## 資料表

- `time_slot_templates`

---

## Service 邏輯

1. 根據 `name` 模糊查詢 `time_slot_templates` 表的 `name` 欄位。
2. 加入 `limit` / `offset` 處理分頁。
3. 回傳總筆數與模板清單。

---

## 注意事項

- 僅回傳模板主資料（不包含子項 item 時段），若需明細請查詢單筆 API。
