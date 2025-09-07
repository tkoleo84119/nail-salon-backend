## User Story

作為一位管理員，我希望能更新某家店的某筆支出明細，方便維護支出資訊。

---

## Endpoint

**PATCH** `/api/admin/stores/{storeId}/expenses/{expenseId}/items/{expenseItemId}`

---

## 說明

- 提供後台管理員更新某家店的某筆支出明細功能。

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
  "productId": "1",
  "quantity": 1,
  "price": 100, // 單價
  "expirationDate": "2025-01-01",
  "isArrived": true,
  "arrivalDate": "2025-01-01",
  "storageLocation": "庫存位置",
  "note": "支出備註",
}
```

### 驗證規則

| 欄位            | 必填 | 其他規則                          | 說明         |
| --------------- | ---- | --------------------------------- | ------------ |
| productId       | 否   |                                   | 商品ID       |
| quantity        | 否   | <li>最小值為 0<li>最大值為1000000 | 數量         |
| price           | 否   | <li>最小值為 0<li>最大值為1000000 | 單價         |
| expirationDate  | 否   | <li>格式為YYYY-MM-DD              | 有限期限     |
| isArrived       | 否   |                                   | 是否已到貨   |
| arrivalDate     | 否   | <li>格式為YYYY-MM-DD              | 到貨日期     |
| storageLocation | 否   | <li>最大長度100字元               | 庫存位置     |
| note            | 否   | <li>最大長度255字元               | 支出明細備註 |

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

| 狀態碼 | 錯誤碼   | 常數名稱                | 說明                                                |
| ------ | -------- | ----------------------- | --------------------------------------------------- |
| 401    | E1002    | AuthTokenInvalid        | 無效的 accessToken，請重新登入                      |
| 401    | E1003    | AuthTokenMissing        | accessToken 缺失，請重新登入                        |
| 401    | E1004    | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入                    |
| 401    | E1005    | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入                    |
| 401    | E1006    | AuthContextMissing      | 未找到使用者認證資訊，請重新登入                    |
| 403    | E1010    | AuthPermissionDenied    | 權限不足，無法執行此操作                            |
| 400    | E2001    | ValJsonFormat           | JSON 格式錯誤，請檢查                               |
| 400    | E2002    | ValPathParamMissing     | 路徑參數缺失，請檢查                                |
| 400    | E2003    | ValAllFieldsEmpty       | 至少需要提供一個欄位進行更新                        |
| 400    | E2004    | ValTypeConversionFailed | 參數類型轉換失敗                                    |
| 400    | E2024    | ValFieldStringMaxLength | {field} 長度最多只能有 {param} 個字元               |
| 400    | E2033    | ValFieldDateFormat      | {field} 格式錯誤，請使用正確的日期格式 (YYYY-MM-DD) |
| 400    | E2036    | ValFieldNoBlank         | {field} 不能為空字串                                |
| 400    | E3PRO003 | ProductNotBelongToStore | 產品不屬於指定的門市                                |
| 404    | E3EXP001 | ExpenseNotFound         | 支出不存在或已被刪除                                |
| 404    | E3PRO002 | ProductNotFound         | 產品不存在或已被刪除                                |
| 500    | E9001    | SysInternalError        | 系統發生錯誤，請稍後再試                            |
| 500    | E9002    | SysDatabaseError        | 資料庫操作失敗                                      |

---

## 資料表

- `products`
- `staff_users`
- `expenses`
- `expense_items`

---

## Service 邏輯

1. 確認門市存取權限。
3. 確認 `expense_item` 是否存在。
2. 確認 `expense` 是否存在。
3. 確認 `expense_item` 是否屬於指定的支出。
4. 如果有傳入`productId`，則確認 `productId` 是否存在，並且屬於指定的門市。
5. 更新 `expenses_items` 資料。
6. 如果有更新 `price` 或 `quantity`，則更新 `expenses` 的總金額。
7. 如果有更新 `quantity`，則更新 `products` 的庫存數量 (扣掉原本的 `quantity`，加上新的 `quantity`)。
8. 回傳更新結果。
