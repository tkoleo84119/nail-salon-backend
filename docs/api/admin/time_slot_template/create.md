## User Story

作為一位管理員，我希望能建立時段範本（template），快速複製排班規劃。

---

## Endpoint

**POST** `/api/admin/time-slot-templates`

---

## 說明

- 可建立一組時段範本（template），用於快速複製與套用在班表。
- 一個範本可包含多個時段（time_slots）。

---

## 權限

- 需要登入才可使用。
- 僅 `SUPER_ADMIN`、`ADMIN`、`MANAGER` 可建立。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Body 範例

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

| 欄位                | 必填 | 其他規則                        | 說明     |
| ------------------- | ---- | ------------------------------- | -------- |
| name                | 是   | <li>最大長度50字元              | 範本名稱 |
| note                | 否   | <li>最大長度100字元             | 備註     |
| timeSlots           | 是   | <li>陣列<li>最少1筆<li>最多50筆 | 多個時段 |
| timeSlots.startTime | 是   | <li>HH:mm 格式                  | 起始時間 |
| timeSlots.endTime   | 是   | <li>HH:mm 格式                  | 結束時間 |

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

### 錯誤處理

#### 錯誤總覽

| 狀態碼 | 錯誤碼   | 說明                                           |
| ------ | -------- | ---------------------------------------------- |
| 401    | E1002    | 無效的 accessToken，請重新登入                 |
| 401    | E1003    | accessToken 缺失，請重新登入                   |
| 401    | E1004    | accessToken 格式錯誤，請重新登入               |
| 401    | E1005    | 未找到有效的員工資訊，請重新登入               |
| 401    | E1006    | 未找到使用者認證資訊，請重新登入               |
| 403    | E1010    | 權限不足，無法執行此操作                       |
| 400    | E2001    | JSON 格式錯誤，請檢查                          |
| 400    | E2004    | 參數類型轉換失敗                               |
| 400    | E2020    | {field} 為必填項目                             |
| 400    | E2022    | {field} 至少需要 {param} 個項目                |
| 400    | E2024    | {field} 長度最多只能有 {param} 個字元          |
| 400    | E2025    | {field} 最多只能有 {param} 個項目              |
| 400    | E2034    | {field} 格式錯誤，請使用正確的時間格式 (HH:mm) |
| 400    | E3TMS011 | 時段時間區段重疊                               |
| 400    | E3TMS012 | 結束時間必須在開始時間之後                     |
| 500    | E9001    | 系統發生錯誤，請稍後再試                       |
| 500    | E9002    | 資料庫操作失敗                                 |

#### 400 Bad Request - 驗證錯誤

```json
{
  "error": {
    "code": "E2020",
    "message": "name 為必填項目",
    "field": "name"
  }
}
```

#### 401 Unauthorized - 認證失敗

```json
{
  "error": {
    "code": "E1002",
    "message": "無效的 accessToken"
  }
}
```

#### 403 Forbidden - 權限不足

```json
{
  "error": {
    "code": "E1010",
    "message": "權限不足，無法執行此操作"
  }
  }
```

#### 500 Internal Server Error - 系統錯誤

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

1. 驗證 `timeSlots` 相關邏輯。
   1. 驗證 `startTime`/`endTime` 格式是否正確。
   2. `startTime` 必須在 `endTime` 之前。
   3. 驗證 `timeSlots` 之間不可重疊。
2. 建立 `time_slot_templates` 資料。
3. 建立對應多筆 `time_slot_template_items` 資料。
4. 回傳建立結果。

---

## 注意事項

- 範本下時段不得重疊。
