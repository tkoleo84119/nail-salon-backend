## User Story

作為管理員，我希望可以查詢單一門市的詳細資料，以便進行設定與後台管理。

---

## Endpoint

**GET** `/api/admin/stores/{storeId}`

---

## 說明

- 僅限 `admin` 以上角色使用。
- 查詢特定門市的完整資訊。

---

## 權限

- 僅限 `admin` 以上角色 (`SUPER_ADMIN`, `ADMIN`)。

---

## Request

### Header

Authorization: Bearer <access_token>

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
    "isActive": true
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

#### 403 Forbidden - 權限不足

```json
{
  "message": "無權限存取此資源"
}
```

#### 404 Not Found - 查無此門市

```json
{
  "message": "查無此門市"
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

---

## Service 邏輯

1. 根據 `storeId` 查詢 `stores` 表。
2. 回傳對應門市完整資訊。

---

## 注意事項

- 不論 `is_active` 狀態，皆允許查詢。

