## User Story

作為一位管理員，我希望能夠管理後台員工帳號

---

## Endpoint

**PATCH** `/api/admin/staff/{staffId}`

---

## 說明

- 提供管理員更新員工帳號資料的功能。
- 僅允許修改角色、是否啟用。

---

## 權限

- 需要登入才可使用。
- `SUPER_ADMIN` 與 `ADMIN` 可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameter

| 參數    | 說明        |
| ------- | ----------- |
| staffId | 員工帳號 ID |

### Body 範例

```json
{
  "role": "MANAGER",
  "isActive": false
}
```

### 驗證規則

| 欄位     | 必填 | 其他規則                             | 說明               |
| -------- | ---- | ------------------------------------ | ------------------ |
| role     | 否   | <li>值只能為 ADMIN、MANAGER、STYLIST | 欲變更的角色       |
| isActive | 否   | <li>布林值                           | 是否啟用該員工帳號 |

- 至少需要提供一個欄位進行更新

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
    "isActive": false,
    "createdAt": "2025-01-01T00:00:00.000Z",
    "updatedAt": "2025-01-01T00:00:00.000Z"
  }
}
```


### 錯誤處理

#### 錯誤總覽

| 狀態碼 | 錯誤碼   | 說明                              |
| ------ | -------- | --------------------------------- |
| 401    | E1002    | 無效的 accessToken，請重新登入    |
| 401    | E1003    | accessToken 缺失，請重新登入      |
| 401    | E1004    | accessToken 格式錯誤，請重新登入  |
| 401    | E1005    | 未找到有效的員工資訊，請重新登入  |
| 401    | E1006    | 未找到使用者認證資訊，請重新登入  |
| 403    | E1010    | 權限不足，無法執行此操作          |
| 400    | E2001    | JSON 格式錯誤，請檢查             |
| 400    | E2003    | 至少需要提供一個欄位進行更新      |
| 400    | E2004    | 參數類型轉換失敗                  |
| 400    | E2020    | {field} 為必填項目                |
| 400    | E2029    | {field} 必須是布林值              |
| 400    | E2030    | {field} 必須是 {param} 其中一個值 |
| 400    | E3STA001 | 無效的角色                        |
| 403    | E3STA004 | 不可更新自己的帳號                |
| 404    | E3STA005 | 員工帳號不存在                    |
| 500    | E9001    | 系統發生錯誤，請稍後再試          |
| 500    | E9002    | 資料庫操作失敗                    |

#### 400 Bad Request - 驗證錯誤

```json
{
  "error": [
    {
      "code": "E2029",
      "message": "isActive 必須是布林值",
      "field": "isActive"
    }
  ]
}
```

#### 401 Unauthorized - 認證失敗

```json
{
  "error": [
    {
      "code": "E1002",
      "message": "無效的 accessToken，請重新登入"
    }
  ]
}
```

#### 403 Forbidden - 權限不足

```json
{
  "error": [
    {
      "code": "E1010",
      "message": "權限不足，無法執行此操作"
    }
  ]
}
```

#### 404 Not Found - 員工不存在

```json
{
  "error": [
    {
      "code": "E3STA005",
      "message": "員工帳號不存在"
    }
  ]
}
```

#### 500 Internal Server Error

```json
{
  "error": [
    {
      "code": "E9001",
      "message": "系統發生錯誤，請稍後再試"
    }
  ]
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

