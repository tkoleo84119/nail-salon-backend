## User Story

作為一位員工，我希望能夠更新自己的密碼，以便存取後台功能。

---

## Endpoint

**POST** `/api/admin/auth/update-password`

---

## 說明

提供後台員工更新密碼功能。

---

## 權限

- 需要登入才可使用。
- 所有角色皆可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Body 範例

```json
{
  "staffId": "123",
  "oldPassword": "hunter2",
  "newPassword": "hunter3"
}
```

### 驗證規則

| 欄位        | 必填 | 其他規則                            | 說明    |
| ----------- | ---- | ----------------------------------- | ------- |
| staffId     | 是   |                                     | 員工 ID |
| oldPassword | 否   | <li>不能為空字串<li>最大長度100字元 | 舊密碼  |
| newPassword | 是   | <li>不能為空字串<li>最大長度100字元 | 新密碼  |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "123"
  }
}
```

---

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

| 狀態碼 | 錯誤碼   | 常數名稱                | 說明                                  |
| ------ | -------- | ----------------------- | ------------------------------------- |
| 401    | E1001    | AuthInvalidCredentials  | 帳號或密碼錯誤                        |
| 401    | E1005    | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入      |
| 401    | E1006    | AuthContextMissing      | 未找到使用者認證資訊，請重新登入      |
| 403    | E1010    | AuthPermissionDenied    | 權限不足，無法執行此操作              |
| 400    | E2001    | ValJsonFormat           | JSON 格式錯誤，請檢查                 |
| 400    | E2020    | ValFieldRequired        | {field} 為必填項目                    |
| 400    | E2024    | ValFieldMaxLength       | {field} 長度最多只能有 {param} 個字元 |
| 400    | E2036    | ValFieldNoBlank         | {field} 不能為空字串                  |
| 400    | E2004    | ValTypeConversionFailed | 參數類型轉換失敗                      |
| 404    | E3STA004 | StaffNotFound           | 員工帳號不存在                        |
| 500    | E9001    | SysInternalError        | 系統發生錯誤，請稍後再試              |
| 500    | E9002    | SysDatabaseError        | 資料庫操作失敗                        |

---

## 實作與流程

### 資料表

- `staff_users`

### Service 邏輯

1. 非 `SUPER_ADMIN`，檢查是否為自己的 `staff_id`
2. 非 `SUPER_ADMIN`，檢查是否傳入 `oldPassword` (一定要傳入)
3. 根據 `id` 查詢 `staff_users`
  - 確認是否存在
  - 檢查 `password_hash`（bcrypt）是否與 `oldPassword` 相符
4. 更新 `password_hash` 為 `newPassword`
5. 回傳更新結果

---

## 注意事項

- 密碼以 bcrypt 儲存
