## User Story

作為一位管理員（`SUPER_ADMIN` / `ADMIN` / `MANAGER`），我希望能刪除時段範本（template）下的一筆 time_slot_template_item，彈性維護排班範本。

---

## Endpoint

**DELETE** `/api/time-slot-templates/{templateId}/items/{itemId}`

---

## 說明

- 僅限管理員可針對指定範本（template）下的特定時段（time_slot_template_item）進行刪除。

---

## 權限

- 僅 `SUPER_ADMIN`、`ADMIN`、`MANAGER` 可操作。

---

## Request

### Header

```http
Content-Type: application/json
Authorization: Bearer <access_token>
```

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
    "deleted": ["6100000003"]
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗"
}
```

#### 401 Unauthorized - 未登入/Token失效

```json
{
  "message": "無效的 accessToken"
}
```

#### 403 Forbidden - 權限不足

```json
{
  "message": "權限不足，僅限管理員操作"
}
```

#### 404 Not Found - 範本或項目不存在

```json
{
  "message": "範本或項目不存在或已被刪除"
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

1. 驗證 `templateId` 是否存在。
2. 驗證 `itemId` 是否存在。
3. 驗證 `itemId` 是否屬於 `templateId`。
4. 刪除 `time_slot_template_items` 資料。
5. 回傳已刪除 id。

---

## 注意事項

- 僅管理員可操作。

