## User Story

作為一位超級管理員（`SUPER_ADMIN`）或系統管理員（`ADMIN`），我希望能夠新增後台員工帳號（如管理員、門市主管或美甲師），以便指派角色與存取特定門市後台功能。

---

## Endpoint

**POST** `/api/staff`

---

## 說明

建立新的員工帳號，並指定其角色與可存取之門市。

---

## 權限

- 僅限 `SUPER_ADMIN` 或 `ADMIN` 存取

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
  "username": "stylist_jane",
  "email": "jane@example.com",
  "password": "hunter2",
  "role": "STYLIST",
  "storeIds": ["1", "2"]
}
```

### 驗證規則

| 欄位     | 規則                                        | 說明                   |
| -------- | ------------------------------------------- | ---------------------- |
| username | <li>必填<li>唯一<li>長度大於1<li>長度小於30 | 員工帳號               |
| email    | <li>必填<li>email格式                       | 員工 Email             |
| password | <li>必填<li>長度大於1<li>長度小於50         | 登入密碼（將加密儲存） |
| role     | <li>必填<li>值只能為ADMIN、MANAGER、STYLIST | 角色                   |
| storeIds | <li>必填<li>至少一筆                        | 有權限的門市 ID 清單   |

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "id": "13984392823",
    "username": "stylist_jane",
    "email": "jane@example.com",
    "role": "STYLIST",
    "storeList": [
      {
        "id": "1",
        "name": "台北忠孝店"
      },
      {
        "id": "2",
        "name": "新竹巨城店"
      }
    ]
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "username": "帳號已存在",
    "password": "密碼至少需 6 個字元"
  }
}
```

#### 401 Unauthorized - 認證失敗

```json
{
  "message": "無效的 accessToken"
}
```

#### 403 Forbidden - 權限不足

```json
{
  "message": "權限不足，無法執行此操作"
}
```

#### 500 Internal Server Error

```json
{
  "message": "系統發生錯誤，請稍後再試"
}
```

---

## 實作

### 資料表

- `staff_users`
- `staff_user_store_access`
- `stores`

### Service 邏輯

1. 檢查傳入的 `storeIds` 是否是該管理員有權限的門市
2. 檢查帳號與 Email 是否唯一（`username`, `email`）
3. 檢查 `storeIds` 是否存在且為啟用中（`is_active = true`）
4. 將密碼加密（bcrypt）後儲存至 `staff_users`
5. 新增 `staff_user_store_access` 多筆紀錄，綁定可控門市
6. 回傳創建成功之帳號資訊（不含密碼）

---

## 注意事項

- 密碼以 bcrypt 儲存，禁止明文存取
- 不可新增 `SUPER_ADMIN` 帳號
- response 中不包含 `password`
