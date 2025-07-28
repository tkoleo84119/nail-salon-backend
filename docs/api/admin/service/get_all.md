## User Story

作為員工，我希望可以查詢某門市底下所有服務資料，並支援條件查詢與分頁，以利管理與設定服務項目。

---

## Endpoint

**GET** `/api/admin/stores/{storeId}/services`

---

## 說明

- 所有登入員工皆可查詢。
- 支援條件搜尋：服務名稱、是否附加服務、啟用狀態、前台是否可見。
- 支援分頁（limit、offset）。

---

## 權限

- 任一已登入員工皆可存取（JWT 驗證）。

---

## Request

### Header

Authorization: Bearer <access_token>

### Path Parameter

| 參數    | 說明    |
| ------- | ------- |
| storeId | 門市 ID |

### Query Parameters

| 參數      | 型別   | 必填 | 預設值 | 說明             |
| --------- | ------ | ---- | ------ | ---------------- |
| name      | string | 否   |        | 模糊查詢服務名稱 |
| isAddon   | bool   | 否   |        | 是否為附加服務   |
| isActive  | bool   | 否   |        | 是否啟用         |
| isVisible | bool   | 否   |        | 前台是否可見     |
| limit     | int    | 否   | 20     | 單頁筆數         |
| offset    | int    | 否   | 0      | 起始筆數         |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "total": 3,
    "items": [
      {
        "id": "9000000001",
        "name": "手部單色",
        "price": 1200,
        "durationMinutes": 60,
        "isAddon": false,
        "isActive": true,
        "isVisible": true,
        "note": "含修型保養"
      }
    ]
  }
}
```

### 失敗

#### 401 Unauthorized - 未登入/Token失效

```json
{
  "message": "無效的 accessToken"
}
```

#### 403 Forbidden - 無權限

```json
{
  "message": "權限不足，無法執行此操作"
}
```

#### 404 Not Found - 門市不存在

```json
{
  "message": "門市不存在或已被刪除"
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

- `stores`
- `services`

---

## Service 邏輯

1. 驗證 `storeId` 是否存在。
2. 驗證 `storeId` 是否為該員工可操作的門市。
3. 根據條件組合查詢 `services`：
   - store_id = storeId
   - name（模糊查詢）
   - is_addon, is_active, is_visible
4. 加入 `limit` / `offset` 處理分頁。
5. 回傳總筆數與項目清單。

---

## 注意事項

- 查詢為該門市綁定的服務，不跨門市資料。

