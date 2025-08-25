## User Story

作為一位管理員，我希望能夠新增後台員工帳號（如管理員、門市主管或美甲師），以便指派角色與存取特定門市後台功能。

---

## Endpoint

**POST** `/api/admin/staff`

---

## 說明

- 提供管理員可建立新的員工帳號功能，並指定其角色與可存取之門市。

---

## 權限

- 需要登入才可使用。
- `SUPER_ADMIN` 與 `ADMIN` 可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Body 範例

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

| 欄位     | 必填 | 其他規則                            | 說明                 |
| -------- | ---- | ----------------------------------- | -------------------- |
| username | 是   | <li>不能為空字串<li>最大長度50字元  | 帳號（唯一）         |
| password | 是   | <li>不能為空字串<li>最大長度50字元  | 密碼明文             |
| email    | 是   | <li>email格式                       | 信箱                 |
| role     | 是   | <li>值只能為ADMIN、MANAGER、STYLIST | 角色                 |
| storeIds | 是   | <li>最少1筆<li>最多10筆             | 有權限的門市 ID 清單 |

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
    "isActive": true,
    "createdAt": "2025-01-01T00:00:00+08:00",
    "updatedAt": "2025-01-01T00:00:00+08:00"
  }
}
```


### 錯誤處理

全部 API 皆回傳如下結構，請參考錯誤總覽。

```json
{
  "errors": [
    {
      "code": "EXXXX",
      "message": "錯誤訊息",
      "field": "錯誤欄位名稱"
    }
  ]
}
```

- 欄位說明：
  - errors: 錯誤陣列（支援多筆同時回報）
  - code: 錯誤代碼，唯一對應每種錯誤
  - message: 中文錯誤訊息（可參照錯誤總覽）
  - field: 參數欄位名稱（僅部分驗證錯誤有）

| 狀態碼 | 錯誤碼   | 常數名稱                | 說明                                       |
| ------ | -------- | ----------------------- | ------------------------------------------ |
| 401    | E1002    | AuthTokenInvalid        | 無效的 accessToken，請重新登入             |
| 401    | E1003    | AuthTokenMissing        | accessToken 缺失，請重新登入               |
| 401    | E1004    | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入           |
| 401    | E1005    | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入           |
| 401    | E1006    | AuthContextMissing      | 未找到使用者認證資訊，請重新登入           |
| 403    | E1010    | AuthPermissionDenied    | 權限不足，無法執行此操作                   |
| 400    | E2001    | ValJsonFormat           | JSON 格式錯誤，請檢查                      |
| 400    | E2002    | ValPathParamMissing     | 路徑參數缺失，請檢查                       |
| 400    | E2004    | ValTypeConversionFailed | 參數類型轉換失敗                           |
| 400    | E2020    | ValFieldRequired        | {field} 為必填項目                         |
| 400    | E2022    | ValFieldArrayMinLength  | {field} 至少需要 {param} 個項目            |
| 400    | E2024    | ValFieldStringMaxLength | {field} 長度最多只能有 {param} 個字元      |
| 400    | E2025    | ValFieldArrayMaxLength  | {field} 最多只能有 {param} 個項目          |
| 400    | E2027    | ValFieldInvalidEmail    | {field} 格式錯誤，請使用正確的電子郵件格式 |
| 400    | E2030    | ValFieldOneof           | {field} 必須是 {param} 其中一個值          |
| 400    | E2036    | ValFieldNoBlank         | {field} 不能為空字串                       |
| 400    | E3STA001 | StaffInvalidRole        | 無效的角色                                 |
| 409    | E3STA007 | StaffAlreadyExists      | 帳號或Email已存在                          |
| 400    | E3STO001 | StoreNotActive          | 門市未啟用                                 |
| 404    | E3STO002 | StoreNotFound           | 門市不存在或已被刪除                       |
| 500    | E9001    | SysInternalError        | 系統發生錯誤，請稍後再試                   |
| 500    | E9002    | SysDatabaseError        | 資料庫操作失敗                             |

---

## 實作流程

### 資料表

- `staff_users`
- `staff_user_store_access`

### Service 邏輯

1. 檢查 `role` 是否為合法值
2. 檢查 `role` 不可為 `SUPER_ADMIN`
3. 根據 creator 的 `role` 檢查是否可以新增 `targetRole` 的帳號
4. 檢查傳入的 `storeIds` 是否是該管理員有權限的門市
5. 檢查 `username` 是否唯一
6. 檢查 `storeIds` 是否存在且為啟用中（`is_active = true`）
7. 將密碼加密（bcrypt）後儲存至 `staff_users`
8. 新增 `staff_user_store_access` 多筆紀錄，綁定有權限的門市
9. 新增 `stylists` 資料，綁定 `staff_user_id`
10. 回傳創建成功之帳號資訊（不含密碼）

---

## 注意事項

- 密碼以 bcrypt 儲存
- 不可新增 `SUPER_ADMIN` 帳號
- response 中不包含 `password`
