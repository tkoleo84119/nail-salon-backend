## User Story

作為一位管理員（`SUPER_ADMIN` / `ADMIN`），我希望能更新門市（store），以維護門市資訊。

- `ADMIN` 僅能操作自己有權限的門市。
- `SUPER_ADMIN` 可操作所有門市。

---

## Endpoint

**PATCH** `/api/admin/stores/{storeId}`

---

## 說明

- 僅限 `SUPER_ADMIN`、`ADMIN` 可更新門市。
- 僅允許修改名稱、地址、電話。
- 門市名稱須唯一。
- `ADMIN` 僅能操作有權限的門市。

---

## 權限

- 僅 `SUPER_ADMIN`、`ADMIN` 可操作。

---

## Request

### Header

```http
Content-Type: application/json
Authorization: Bearer <access_token>
```

### Path Parameter

| 參數    | 說明   |
| ------- | ------ |
| storeId | 門市ID |

### Body

```json
{
  "name": "松江南京分店",
  "address": "台北市中山區松江路123號",
  "phone": "02-88889999",
  "isActive": true
}
```

### 驗證規則

| 欄位     | 規則                                             | 說明     |
| -------- | ------------------------------------------------ | -------- |
| name     | <li>選填<li>長度大於1<li>長度小於100<li>唯一     | 門市名稱 |
| address  | <li>選填<li>長度小於255                          | 門市地址 |
| phone    | <li>選填<li>長度小於20<li>格式必須為台灣市話號碼 | 電話     |
| isActive | <li>選填                                         | 是否啟用 |

- 欄位皆為選填，但至少需有一項。

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "8000000001",
    "name": "松江南京分店",
    "address": "台北市中山區松江路123號",
    "phone": "02-88889999",
    "isActive": true
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "phone": "必須為有效的台灣市話號碼格式 (例: 02-12345678)"
  }
}
```

#### 401 Unauthorized - 未登入/Token失效

```json
{
  "message": "無效的 accessToken"
}
```

#### 403 Forbidden - 權限不足

```json
{
  "message": "權限不足，僅限有權限的管理員操作"
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
- `staff_user_store_access`

---

## Service 邏輯

1. 驗證至少一個欄位有更新。
2. 驗證 `storeId` 是否存在。
3. 驗證 `name` 是否唯一（若有更改，不包含自己）。
4. 更新 `stores` 資料。
5. 回傳更新結果。

---

## 注意事項

- `ADMIN` 僅能操作自己有權限的門市。
- 門市名稱不可重複（不包含自己）。
- 僅允許 name、address、phone 欄位修改。
