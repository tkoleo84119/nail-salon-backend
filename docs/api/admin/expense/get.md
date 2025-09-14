## User Story

作為管理員，我希望可以查試單一支出詳細資料，方便管理。

---

## Endpoint

**GET** `/api/admin/stores/{storeId}/expenses/{expenseId}`

---

## 說明

- 用於查詢特定支出的詳細資訊。

---

## 權限

- 需要登入才可使用。
- 僅 `SUPER_ADMIN`、`ADMIN` `MANAGER` 可操作。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameters

| 參數      | 說明    |
| --------- | ------- |
| storeId   | 門市 ID |
| expenseId | 支出 ID |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "id": "9000000001",
    "supplier": {
      "id": "9000000001",
      "name": "供應商A"
    }, // 如果沒有供應商沒有此欄位
    "payer": {
      "id": "9000000001",
      "name": "員工A"
    }, // 如果沒有代墊人沒有此欄位
    "category": "薪資",
    "amount": 200, // 支出金額 (排除其他費用後的金額)
    "otherFee": 10, // 其他費用
    "expenseDate": "2025-01-01",
    "note": "薪資備註",
    "isReimbursed": true, // 如果沒有代墊人沒有此欄位
    "reimbursedAt": "2025-01-01", // 如果沒有代墊人沒有此欄位
    "updater": "員工A", // 如果是 "" 就是系統更新
    "createdAt": "2025-01-01T00:00:00+08:00",
    "updatedAt": "2025-01-01T00:00:00+08:00",
    "items": [
      {
        "id": "9000000001",
        "product": {
          "id": "9000000001",
          "name": "商品A"
        },
        "quantity": 2,
        "price": 100, // 單價
        "expirationDate": "2025-01-01", // 如果沒有限期限沒有此欄位
        "isArrived": true,
        "arrivalDate": "2025-01-01", // 如果還沒有到貨沒有此欄位
        "storageLocation": "庫存位置",
        "note": "支出備註"
      }
    ] // 只有category為"進貨"才有項目
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

| 狀態碼 | 錯誤碼   | 常數名稱                | 說明                             |
| ------ | -------- | ----------------------- | -------------------------------- |
| 401    | E1002    | AuthTokenInvalid        | 無效的 accessToken，請重新登入   |
| 401    | E1003    | AuthTokenMissing        | accessToken 缺失，請重新登入     |
| 401    | E1004    | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入 |
| 401    | E1005    | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入 |
| 401    | E1006    | AuthContextMissing      | 未找到使用者認證資訊，請重新登入 |
| 400    | E2002    | ValPathParamMissing     | 路徑參數缺失，請檢查             |
| 400    | E2004    | ValTypeConversionFailed | 參數類型轉換失敗                 |
| 404    | E3EXP001 | ExpenseNotFound         | 支出不存在或已被刪除             |
| 500    | E9001    | SysInternalError        | 系統發生錯誤，請稍後再試         |
| 500    | E9002    | SysDatabaseError        | 資料庫操作失敗                   |

---

## 實作與流程

### 資料表

- `expenses`
- `suppliers`
- `expense_items`
- `products`
- `staff_users`

---

### Service 邏輯

1. 驗證員工是否有權限查詢該門市。
2. 查詢 `expenses` 表中該筆支出是否存在。
3. 確認該 `expenses` 是否隸屬於該門市。
4. 查詢 `expense_items` 表中該筆支出的詳細資訊。
5. 整理回傳資料。

---

## 注意事項

- `createdAt` 與 `updatedAt` 會是標準 Iso 8601 格式。
