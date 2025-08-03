## User Story

作為員工，我希望可以查詢特定服務資料，以便管理或顯示詳細內容。

---

## Endpoint

**GET** `/api/admin/services/{serviceId}`

---

## 說明

- 用於查詢特定服務的詳細資訊。

---

## 權限

- 需要登入才可使用。
- 所有角色皆可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameters

| 參數      | 說明    |
| --------- | ------- |
| serviceId | 服務 ID |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "9000000001",
    "name": "手部單色",
    "durationMinutes": 60,
    "price": 1200,
    "isAddon": false,
    "isActive": true,
    "isVisible": true,
    "note": "含修型保養",
    "createdAt": "2025-01-01T00:00:00+08:00",
    "updatedAt": "2025-01-01T00:00:00+08:00"
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
| 400    | E2004    | 參數類型轉換失敗                 |
| 400    | E2020    | {field} 為必填項目               |
| 404    | E3SER004 | 服務不存在或已被刪除             |
| 500    | E9001    | 系統發生錯誤，請稍後再試         |
| 500    | E9002    | 資料庫操作失敗                   |

#### 400 Bad Request - 參數類型轉換失敗

```json
{
  "error": {
    "code": "E2004",
    "message": "serviceId 類型轉換失敗",
    "field": "serviceId"
  }
}
```


#### 401 Unauthorized - 未登入/Token失效

```json
{
  "errors": [
    {
      "code": "E1002",
      "message": "無效的 accessToken，請重新登入"
    }
  ]
}
```

#### 403 Forbidden - 無權限

```json
{
  "errors": [
    {
      "code": "E1010",
      "message": "權限不足，無法執行此操作"
    }
  ]
}
```

#### 404 Not Found - 找不到門市或服務

```json
{
  "errors": [
    {
      "code": "E3SER004",
      "message": "服務不存在或已被刪除"
    }
  ]
}
```

#### 500 Internal Server Error - 系統錯誤

```json
{
  "errors": [
    {
      "code": "E9001",
      "message": "系統發生錯誤，請稍後再試"
    }
  ]
}
```

---

## 資料表

- `services`

---

## Service 邏輯

1. 查詢 `services` 表中該筆服務是否存在。
2. 不存在則回傳 `404 Not Found`。
3. 存在則回傳該筆服務詳細內容。
