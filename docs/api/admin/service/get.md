## User Story

作為員工，我希望可以查詢某門市底下的特定服務資料，以便管理或顯示詳細內容。

---

## Endpoint

**GET** `/api/admin/stores/{storeId}/services/{serviceId}`

---

## 說明

- 所有登入員工皆可查詢。
- 用於查詢特定服務的詳細資訊。
- 僅限該門市底下綁定的服務資料。

---

## 權限

- 任一已登入員工皆可使用（JWT 驗證）。

---

## Request

### Header

Authorization: Bearer <access_token>

### Path Parameters

| 參數      | 說明    |
| --------- | ------- |
| storeId   | 門市 ID |
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
    "note": "含修型保養"
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

#### 404 Not Found - 找不到門市或服務

```json
{
  "message": "門市不存在或已被刪除"
}
```

```json
{
  "message": "查無此門市或服務"
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
2. 驗證員工是否擁有該門市存取權限。
3. 查詢 `services` 表中該筆服務是否存在，且 `store_id=storeId`。
4. 不存在則回傳 `404 Not Found`。
5. 存在則回傳該筆服務詳細內容。

---

## 注意事項

- 僅能查詢特定門市底下的服務（跨門市查詢無效）。

