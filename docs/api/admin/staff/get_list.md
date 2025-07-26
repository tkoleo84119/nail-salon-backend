## User Story

作為管理員，我希望可以查詢所有員工資料，並支援查詢條件與分頁，方便管理與篩選。

---

## Endpoint

**GET** `/api/admin/staff`

---

## 說明

- 僅限 admin 以上權限使用（`SUPER_ADMIN`, `ADMIN`）。
- 支援基本查詢條件，如：姓名、帳號、角色、啟用狀態等。
- 支援分頁（`limit`, `offset`）。
- 回傳項目包含：ID、帳號、信箱、角色、啟用狀態、建立時間。

---

## 權限

- 僅限 `admin` 以上角色（JWT + RBAC 權限驗證）。

---

## Request

### Header

Authorization: Bearer <access_token>

### Query Parameter

| 參數     | 型別   | 必填 | 說明                            |
| -------- | ------ | ---- | ------------------------------- |
| keyword  | string | 否   | 可模糊搜尋帳號/信箱/姓名        |
| role     | string | 否   | 欲篩選角色（ADMIN、STYLIST...） |
| isActive | bool   | 否   | 是否啟用帳號                    |
| limit    | int    | 否   | 單頁筆數（預設 20）             |
| offset   | int    | 否   | 起始筆數（預設 0）              |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "total": 3,
    "items": [
      {
        "id": "6000000001",
        "username": "admin01",
        "email": "admin01@example.com",
        "role": "ADMIN",
        "isActive": true,
        "createdAt": "2025-05-01T08:00:00Z"
      },
      {
        "id": "6000000002",
        "username": "stylist88",
        "email": "s88@salon.com",
        "role": "STYLIST",
        "isActive": false,
        "createdAt": "2025-06-01T08:00:00Z"
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

- `staff_users`

---

## Service 邏輯

1. 根據查詢條件動態查詢資料（模糊搜尋、角色、啟用狀態），並帶入 `limit` 與 `offset` 參數。
2. 回傳 `items` 與 `total`。

---

## 注意事項

- keyword 模糊搜尋帳號/信箱/姓名
