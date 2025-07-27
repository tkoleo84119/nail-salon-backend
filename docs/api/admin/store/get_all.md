## User Story

作為管理員，我希望可以查詢全部門市資料，並支援查詢條件與分頁，方便管理與篩選。

---

## Endpoint

**GET** `/api/admin/stores`

---

## 說明

- 僅限 `admin` 以上角色可查詢。
- 支援基本查詢條件，如門市名稱、啟用狀態等。
- 支援分頁（limit、offset）。

---

## 權限

- 僅限 `admin` 以上角色（`SUPER_ADMIN`, `ADMIN`）。

---

## Request

### Header

Authorization: Bearer <access_token>

### Query Parameters

| 參數     | 型別   | 必填 | 說明                    |
| -------- | ------ | ---- | ----------------------- |
| keyword  | string | 否   | 模糊查詢門市名稱 / 地址 |
| isActive | bool   | 否   | 是否啟用門市            |
| limit    | int    | 否   | 單頁筆數（預設 20）     |
| offset   | int    | 否   | 起始筆數（預設 0）      |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "total": 2,
    "items": [
      {
        "id": "8000000001",
        "name": "大安旗艦店",
        "address": "台北市大安區復興南路一段100號",
        "phone": "02-1234-5678",
        "isActive": true
      },
      {
        "id": "8000000003",
        "name": "信義分店",
        "address": "台北市信義區松壽路9號",
        "phone": "02-3333-8888",
        "isActive": false
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

#### 403 Forbidden - 權限不足

```json
{
  "message": "無權限存取此資源"
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

1. 根據 `keyword`（名稱/地址）與 `is_active` 條件動態查詢。
2. 加入 `limit` 與 `offset` 處理分頁。
3. 回傳結果與總筆數。

---

## 注意事項

- 未來可能添加排序等。
