## User Story

作為一位管理員（`SUPER_ADMIN` / `ADMIN`），我希望能新增門市（store），方便多據點管理。

- `ADMIN` 新增時，會自動給予該門市權限。
- `SUPER_ADMIN` 原本就擁有所有門市權限。

---

## Endpoint

**POST** `/api/admin/stores`

---

## 說明

- 僅 `SUPER_ADMIN`、`ADMIN` 可建立新門市。
- 門市名稱須唯一。
- `ADMIN` 創建後會自動與該門市建立權限關聯（staff_user_store_access）。

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

### Body

```json
{
  "name": "大安旗艦店",
  "address": "台北市大安區復興南路一段100號",
  "phone": "02-12345678"
}
```

### 驗證規則

| 欄位    | 規則                                             | 說明     |
| ------- | ------------------------------------------------ | -------- |
| name    | <li>必填<li>長度大於1<li>長度小於100<li>唯一     | 門市名稱 |
| address | <li>選填<li>長度小於255                          | 門市地址 |
| phone   | <li>選填<li>長度小於20<li>格式必須為台灣市話號碼 | 電話     |

---

## Response

### 成功 201 Created

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

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "name": "name為必填項目"
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
  "message": "權限不足，僅限管理員新增門市"
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

1. 驗證 `name` 是否唯一。
2. 建立 `stores` 資料。
3. 若為 `ADMIN` 則於 `staff_user_store_access` 關聯該門市。
4. 回傳新增結果。

---

## 注意事項

- 門市名稱不可重複。
- `ADMIN` 新增門市時會自動取得該門市權限。
- `SUPER_ADMIN` 則原本就擁有所有門市權限。

