## User Story

作為一位管理員，我希望能刪除時段範本（template），以維護有效的排班模板。

---

## Endpoint

**DELETE** `/api/admin/time-slot-templates/{templateId}`

---

## 說明

- 刪除範本時會一併刪除該範本下所有時段（time_slot_template_items）。

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

---

## Response

### 成功 204 No Content

```json
{
  "data": {
    "deleted": "6000000011"
  }
}
```

### 錯誤處理

#### 錯誤總覽

| 狀態碼 | 錯誤碼   | 說明                             |
| ------ | -------- | -------------------------------- |
| 401    | E1002    | 無效的 accessToken，請重新登入   |
| 401    | E1003    | accessToken 缺失，請重新登入     |
| 401    | E1004    | accessToken 格式錯誤，請重新登入 |
| 401    | E1005    | 未找到有效的員工資訊，請重新登入 |
| 401    | E1006    | 未找到使用者認證資訊，請重新登入 |
| 403    | E1010    | 權限不足，無法執行此操作         |
| 400    | E2002    | 路徑參數缺失，請檢查             |
| 400    | E2020    | {field} 為必填項目               |
| 404    | E3TMS009 | 範本不存在或已被刪除             |
| 500    | E9001    | 系統發生錯誤，請稍後再試         |
| 500    | E9002    | 資料庫操作失敗                   |

#### 400 Bad Request - 驗證錯誤

```json
{
  "error": {
    "code": "E2002",
    "message": "路徑參數缺失，請檢查"
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

1. 驗證 `templateId` 是否存在。
2. 刪除 `time_slot_templates`。
3. 回傳已刪除 id。

---

## 注意事項

- 刪除範本時會一併移除所有對應的範本時段。
