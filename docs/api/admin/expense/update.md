## User Story

作為一位管理員，我希望能更新某家店的支出，方便維護支出資訊。

---

## Endpoint

**PATCH** `/api/admin/stores/{storeId}/expenses/{expenseId}`

---

## 說明

- 提供後台管理員更新某家店的支出功能。

---

## 權限

- 需要登入才可使用。
- 僅 `SUPER_ADMIN`、`ADMIN` `MANAGER` 可操作。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Body 範例

```json
{
  "supplierId": "1",
  "category": "薪資",
  "amount": 10000, // 支出金額 (排除其他費用後的金額)
  "otherFee": null, // 其他費用
  "expenseDate": "2025-01-01",
  "note": "支出備註",
  "payerId": "1", // 如果要取消代墊，請提供 ""
  "isReimbursed": true,
  "reimbursedAt": "2025-01-01"
}
```

### 驗證規則

| 欄位         | 必填 | 其他規則                            | 說明                            |
| ------------ | ---- | ----------------------------------- | ------------------------------- |
| supplierId   | 否   |                                     | 供應商ID                        |
| category     | 否   | <li>不能為空字串<li>最大長度100字元 | 支出類別                        |
| amount       | 否   | <li>最小值為 0<li>最大值為1000000   | 支出金額 (排除其他費用後的金額) |
| otherFee     | 否   | <li>最小值為 0<li>最大值為1000000   | 其他費用                        |
| expenseDate  | 否   | <li>格式為YYYY-MM-DD                | 支出日期                        |
| note         | 否   | <li>最大長度255字元                 | 支出備註                        |
| payerId      | 否   |                                     | 代墊人ID                        |
| isReimbursed | 否   |                                     | 是否已結清代墊款                |
| reimbursedAt | 否   | <li>格式為YYYY-MM-DD                | 結清日期                        |

- 至少需要提供一個欄位進行更新

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "9000000001"
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

| 狀態碼 | 錯誤碼   | 常數名稱                                     | 說明                                                |
| ------ | -------- | -------------------------------------------- | --------------------------------------------------- |
| 401    | E1002    | AuthTokenInvalid                             | 無效的 accessToken，請重新登入                      |
| 401    | E1003    | AuthTokenMissing                             | accessToken 缺失，請重新登入                        |
| 401    | E1004    | AuthTokenFormatError                         | accessToken 格式錯誤，請重新登入                    |
| 401    | E1005    | AuthStaffFailed                              | 未找到有效的員工資訊，請重新登入                    |
| 401    | E1006    | AuthContextMissing                           | 未找到使用者認證資訊，請重新登入                    |
| 403    | E1010    | AuthPermissionDenied                         | 權限不足，無法執行此操作                            |
| 400    | E2001    | ValJsonFormat                                | JSON 格式錯誤，請檢查                               |
| 400    | E2002    | ValPathParamMissing                          | 路徑參數缺失，請檢查                                |
| 400    | E2003    | ValAllFieldsEmpty                            | 至少需要提供一個欄位進行更新                        |
| 400    | E2004    | ValTypeConversionFailed                      | 參數類型轉換失敗                                    |
| 400    | E2020    | ValFieldRequired                             | {field} 為必填項目                                  |
| 400    | E2024    | ValFieldStringMaxLength                      | {field} 長度最多只能有 {param} 個字元               |
| 400    | E2033    | ValFieldDateFormat                           | {field} 格式錯誤，請使用正確的日期格式 (YYYY-MM-DD) |
| 400    | E2036    | ValFieldNoBlank                              | {field} 不能為空字串                                |
| 400    | E3EXP002 | ExpenseNotUpdateReimbursedInfoWithoutPayerID | 沒有代墊人，無法更新結清資訊                        |
| 400    | E3EXP003 | ExpenseNotUpdateAmountWithExpenseItems       | 有產品明細，無法單獨更新金額                        |
| 404    | E3SUP002 | SupplierNotFound                             | 供應商不存在或已被刪除                              |
| 404    | E3STA004 | StaffNotFound                                | 員工帳號不存在                                      |
| 404    | E3EXP001 | ExpenseNotFound                              | 支出不存在或已被刪除                                |
| 500    | E9001    | SysInternalError                             | 系統發生錯誤，請稍後再試                            |
| 500    | E9002    | SysDatabaseError                             | 資料庫操作失敗                                      |

---

## 資料表

- `suppliers`
- `staff_users`
- `expenses`
- `expense_items`

---

## Service 邏輯

1. 確認門市存取權限。
2. 如果有傳入`supplierId`，則確認 `supplierId` 是否存在。
3. 如果有傳入`payerId`，則確認 `payerId` 是否存在，並且擁有該店權限(`staff_user_store_access`)。
4. 如果有傳入`isReimbursed` 或 `reimbursedAt`，則確認 `payerId` 是否存在。
5. 如果有修改 `amount`，則確認是否有 `expense_items` 資料，有的話不允許透過該 API 修改。
6. 如果有傳入`payerId` 且為空字串，則將 `payerId` 和 `isReimbursed` 和 `reimbursedAt` 設為 `null`。
7. 更新 `expenses` 資料。
8. 回傳更新結果。
