## User Story

作為一位超級管理員（`SUPER_ADMIN`）或系統管理員（`ADMIN`），我希望能夠管理後台員工帳號， 包含：

1. 修改其角色（如：由 `STYLIST` 轉為 `MANAGER`），以調整其權限範圍。
2. 停用或重新啟用帳號，以控管系統存取權限。

---

## Endpoint

**PATCH** `/api/admin/staff/{staffId}`

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
  "role": "MANAGER",
  "isActive": false
}
```

- 至少需要提供一個欄位進行更新

### 驗證規則

| 欄位     | 規則                                         | 說明               |
| -------- | -------------------------------------------- | ------------------ |
| role     | <li>可選<li>值只能為 ADMIN、MANAGER、STYLIST | 欲變更的角色       |
| isActive | <li>可選<li>布林值                           | 是否啟用該員工帳號 |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "13984392823",
    "username": "staff_amy",
    "email": "amy@example.com",
    "role": "MANAGER",
    "isActive": false
  }
}
```

### 失敗

#### 400 Bad Request - 驗證錯誤

```json
{
  "message": "輸入驗證失敗",
  "errors": {
    "role": "role只可以傳入特定值"
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

#### 404 Not Found - 員工不存在

```json
{
  "message": "指定的員工不存在"
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

1. 驗證請求至少需要提供一個欄位進行更新
2. 當 `role` 有傳入時，驗證 `role` 是否為合法值
3. 確認目標員工是否存在
4. 根據角色檢查是否可以更新
5. 更新 `staff_users` 的 `role` 與 `is_active` 欄位
6. 回傳更新後資訊

---

## 注意事項

- 不可修改自身帳號的 `role` 或 `is_active`
- 不可修改 `SUPER_ADMIN` 的帳號狀態與角色（僅能由系統預設）

