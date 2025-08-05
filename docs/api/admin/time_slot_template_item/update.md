## User Story

作為一位管理員，我希望能更新時段範本（template）下的某筆 time_slot_template_item，彈性調整預設時段。

---

## Endpoint

**PATCH** `/api/admin/time-slot-templates/{templateId}/items/{itemId}`

---

## 說明

- 僅支援同時更新 `startTime` 與 `endTime`。

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
| itemId     | 項目ID |

### Body 範本

```json
{
  "startTime": "14:00",
  "endTime": "18:00"
}
```

### 驗證規則

| 欄位      | 必填 | 其他規則       | 說明     |
| --------- | ---- | -------------- | -------- |
| startTime | 是   | <li>HH:mm 格式 | 起始時間 |
| endTime   | 是   | <li>HH:mm 格式 | 結束時間 |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "6100000003",
    "templateId": "6000000011",
    "startTime": "14:00",
    "endTime": "18:00"
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
| 400    | E2002    | 路徑參數缺失，請檢查                           |
| 400    | E2004    | 參數類型轉換失敗                               |
| 400    | E2020    | {field} 為必填項目                             |
| 400    | E2034    | {field} 格式錯誤，請使用正確的時間格式 (HH:mm) |
| 400    | E3TMS011 | 時段時間區段重疊                               |
| 400    | E3TMS012 | 結束時間必須在開始時間之後                     |
| 404    | E3TMS009 | 範本不存在或已被刪除                           |
| 404    | E3TMS010 | 範本項目不存在或已被刪除                       |
| 500    | E9001    | 系統發生錯誤，請稍後再試                       |
| 500    | E9002    | 資料庫操作失敗                                 |

#### 400 Bad Request - 驗證錯誤

```json
{
  "error": {
    "code": "E2001",
    "message": "JSON 格式錯誤，請檢查",
    "field": "body"
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
- `time_slot_template_items`

---

## Service 邏輯

1. 驗證 `startTime`/`endTime` 格式。
2. 驗證 `startTime` 必須在 `endTime` 之前。
3. 驗證 `templateId` 是否存在。
4. 驗證 `itemId` 是否存在。
5. 驗證新時間是否與同範本其他時段重疊。
6. 更新 `time_slot_template_items` 資料。
7. 回傳更新結果。

---

## 注意事項

- 僅允許 `startTime` 與 `endTime` 欄位一起修改。
- 調整後不得與其他 `time_slot_template_item` 重疊。
