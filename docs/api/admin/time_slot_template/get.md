## User Story

作為員工，我希望可以查詢特定的時段模板資料，並一併取得該模板所包含的所有時段項目，以便確認內容或套用排班。

---

## Endpoint

**GET** `/api/admin/time-slot-templates/{templateId}`

---

## 說明

- 所有登入員工皆可查詢。
- 回傳模板主資料與所有時間項目（Time Slot Template Items）。

---

## 權限

- 任一已登入員工皆可使用（JWT 驗證）。

---

## Request

### Header

Authorization: Bearer <access_token>

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
    "createdAt": "2025-06-01T08:00:00Z",
    "updatedAt": "2025-06-20T08:00:00Z",
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

### 失敗

#### 401 Unauthorized - 未登入

```json
{
  "message": "無效的 accessToken"
}
```

#### 404 Not Found - 模板不存在

```json
{
  "message": "範本不存在或已被刪除"
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
- `time_slot_template_items`

---

## Service 邏輯

1. 查詢指定的 `templateId` 是否存在於 `time_slot_templates`。
2. 查詢所有對應的 `time_slot_template_items`，依 `start_time` 排序。
3. 回傳主檔與時段項目合併結果。

---

## 注意事項

- 每筆 item 僅含時間區段

