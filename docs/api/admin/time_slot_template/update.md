## User Story

作為一位管理員，我希望能更新時段範本（template）的名稱與備註，方便後續管理。

---

## Endpoint

**PATCH** `/api/admin/time-slot-templates/{templateId}`

---

## 說明

- 僅支援名稱與備註修改，範本下的時段（`time_slot_template_items`）不在本 API 內調整。

---

## 權限

- 需要登入才可使用。
- 僅 `SUPER_ADMIN`、`ADMIN`、`MANAGER` 可操作。

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
  "name": "新標準早班",
  "note": "夏季適用"
}
```

### 驗證規則

| 欄位 | 必填 | 其他規則                            | 說明     |
| ---- | ---- | ----------------------------------- | -------- |
| name | 否   | <li>最短長度1字元<li>最大長度50字元 | 範本名稱 |
| note | 否   | <li>最大長度100字元                 | 備註     |

- 兩者皆為選填，但至少要有一項。

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "6000000011",
    "name": "新標準早班",
    "note": "夏季適用",
    "updater": "1000000001",
    "createdAt": "2025-01-01T00:00:00+08:00",
    "updatedAt": "2025-01-01T00:00:00+08:00"
  }
}
```

### 錯誤處理

#### 錯誤總覽

| 狀態碼 | 錯誤碼   | 說明                                  |
| ------ | -------- | ------------------------------------- |
| 401    | E1002    | 無效的 accessToken，請重新登入        |
| 401    | E1003    | accessToken 缺失，請重新登入          |
| 401    | E1004    | accessToken 格式錯誤，請重新登入      |
| 401    | E1005    | 未找到有效的員工資訊，請重新登入      |
| 401    | E1006    | 未找到使用者認證資訊，請重新登入      |
| 403    | E1010    | 權限不足，無法執行此操作              |
| 400    | E2001    | JSON 格式錯誤，請檢查                 |
| 400    | E2002    | 路徑參數缺失，請檢查                  |
| 400    | E2003    | 至少需要提供一個欄位進行更新          |
| 400    | E2004    | 參數類型轉換失敗                      |
| 400    | E2020    | {field} 為必填項目                    |
| 400    | E2024    | {field} 長度最多只能有 {param} 個字元 |
| 404    | E3TMS009 | 範本不存在或已被刪除                  |
| 500    | E9001    | 系統發生錯誤，請稍後再試              |
| 500    | E9002    | 資料庫操作失敗                        |

#### 400 Bad Request - 驗證錯誤

```json
{
  "error": {
    "code": "E2024",
    "message": "name 長度最多只能有 50 個字元",
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

#### 404 Not Found - 範本不存在或已被刪除

```json
{
  "error": {
    "code": "E3TMS009",
    "message": "範本不存在或已被刪除"
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

---

## Service 邏輯

1. 驗證至少一個欄位有更新。
2. 驗證 `templateId` 是否存在。
3. 更新 `time_slot_templates` 資料。
4. 回傳更新結果。

---

## 注意事項

- 僅允許 `name` 與 `note` 欄位修改。
