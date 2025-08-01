## User Story

作為員工，我希望可以查詢單一門市的詳細資料，以便進行設定與後台管理。

---

## Endpoint

**GET** `/api/admin/stores/{storeId}`

---

## 說明

- 查詢特定門市的完整資訊。

---

## 權限

- 需要登入才可使用。
- 所有角色皆可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameter

| 參數    | 說明    |
| ------- | ------- |
| storeId | 門市 ID |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "8000000001",
    "name": "大安旗艦店",
    "address": "台北市大安區復興南路一段100號",
    "phone": "02-1234-5678",
    "isActive": true,
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
| 404    | E3STO002 | 門市不存在或已被刪除             |
| 500    | E9001    | 系統發生錯誤，請稍後再試         |
| 500    | E9002    | 資料庫操作失敗                   |

#### 400 Bad Request - 輸入驗證失敗

```json
{
  "errors": [
    {
      "code": "E2004",
      "message": "storeId 類型轉換失敗",
      "field": "storeId"
    }
  ]
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

#### 404 Not Found - 查無此門市

```json
{
  "errors": [
    {
      "code": "E3STO002",
      "message": "門市不存在或已被刪除"
    }
  ]
}
```

#### 500 Internal Server Error

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

- `stores`

---

## Service 邏輯

1. 根據 `storeId` 查詢 `stores` 表。
2. 回傳對應門市完整資訊。

---

## 注意事項

- 不論 `is_active` 狀態，皆允許查詢。

