## User Story

作為一位管理員（`SUPER_ADMIN` / `ADMIN` / `MANAGER`），我希望能建立時段範本（template），快速複製排班規劃。

---

## Endpoint

**POST** `/api/admin/time-slot-templates`

---

## 說明

- 可建立一組時段範本（template），用於快速複製與套用在班表。
- 一個範本可包含多個時段（time_slots），僅設置時間，無須關聯 schedule。
- 範本僅限管理員建立。

---

## 權限

- 僅 `SUPER_ADMIN`、`ADMIN`、`MANAGER` 可建立。

---

## Request

### Header

```http
Content-Type: application/json
Authorization: Bearer <access_token>
```

### Body

```json
{
  "name": "標準早班",
  "note": "適用平日",
  "timeSlots": [
    { "startTime": "09:00", "endTime": "12:00" },
    { "startTime": "13:00", "endTime": "18:00" }
  ]
}
```

### 驗證規則

| 欄位                | 規則                                | 說明     |
| ------------------- | ----------------------------------- | -------- |
| name                | <li>必填<li>長度大於1<li>長度小於50 | 範本名稱 |
| note                | <li>選填<li>長度小於100             | 備註     |
| timeSlots           | <li>必填<li>陣列<li>至少一筆        | 多個時段 |
| timeSlots.startTime | <li>必填<li>HH:mm 格式              | 起始時間 |
| timeSlots.endTime   | <li>必填<li>HH:mm 格式              | 結束時間 |

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "id": "6000000011",
    "name": "標準早班",
    "note": "適用平日",
    "timeSlots": [
      { "id": "6100000001", "startTime": "09:00", "endTime": "12:00" },
      { "id": "6100000002", "startTime": "13:00", "endTime": "18:00" }
    ]
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "name": "name為必填"
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
  "message": "權限不足，僅限管理員建立時段範本"
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

1. 驗證 `timeSlots` 相關邏輯。
   1. 驗證 startTime/endTime 格式是否正確。
   2. startTime 必須在 endTime 之前。
   3. 驗證 timeSlots 之間不可重疊。
2. 建立 `time_slot_templates` 資料。
3. 建立對應多筆 `time_slot_template_items` 資料。
4. 回傳建立結果。

---

## 注意事項

- 僅管理員可建立時段範本。
- 範本下時段不得重疊。

