## User Story

作為一位管理員（`SUPER_ADMIN` / `ADMIN` / `MANAGER`），我希望能在時段範本（template）下新增一筆 time_slot_template_item，靈活調整範本內容。

---

## Endpoint

**POST** `/api/time-slot-templates/{templateId}/items`

---

## 說明

- 僅限管理員可針對指定範本（template）新增一筆時段（time_slot_template_item）。
- 用於動態增添時段，調整排班模板。

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

### Body

```json
{
  "startTime": "10:00",
  "endTime": "13:00"
}
```

### 驗證規則

| 欄位      | 規則                   | 說明     |
| --------- | ---------------------- | -------- |
| startTime | <li>必填<li>HH:mm 格式 | 起始時間 |
| endTime   | <li>必填<li>HH:mm 格式 | 結束時間 |

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

### 失敗

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "startTime": "startTime為必填項目"
  }
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

#### 404 Not Found - 範本不存在

```json
{
  "message": "範本不存在或已被刪除"
}
```

#### 409 Conflict - 時段重疊

```json
{
  "message": "時段重疊，不可新增重複時間"
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
2. 驗證 `startTime`/`endTime` 格式。
3. 驗證 `startTime` 必須在 `endTime` 之前。
4. 驗證新時段是否與範本內其他時段重疊。
5. 建立 `time_slot_template_item`。
6. 回傳建立結果。

---

## 注意事項

- 僅管理員可操作。
- 時段不得重疊。
