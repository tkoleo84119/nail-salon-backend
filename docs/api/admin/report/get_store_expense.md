## User Story

作為管理員，我希望可以查詢指定店家的支出，以便了解店家各項支出的情況。

---

## Endpoint

**GET** `/api/admin/reports/expense/store/{storeId}`

---

## 說明

- 用於查詢指定店家的支出。

---

## 權限

- 需要登入才可使用。
- `SUPER_ADMIN` 與 `ADMIN` 可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Query Parameters

| 參數      | 型別   | 必填 | 預設值 | 說明     |
| --------- | ------ | ---- | ------ | -------- |
| startDate | string | 是   |        | 開始日期 |
| endDate   | string | 是   |        | 結束日期 |

### 驗證規則

| 欄位      | 必填 | 其他規則            |
| --------- | ---- | ------------------- |
| startDate | 是   | <li>YYYY-MM-DD 格式 |
| endDate   | 是   | <li>YYYY-MM-DD 格式 |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "startDate": "2025-01-01",
    "endDate": "2025-01-31",
    "summary": {
      "totalCount": 100, // 總支出筆數
      "totalAmount": 200000, // 總支出金額(含其他費用)
      "advanceAmount": 10000, // 代墊總金額(含其他費用)
      "reimbursedAmount": 5000, // 已結清代墊金額
      "pendingAmount": 5000 // 待結清代墊金額
    },
    "category": [
      {
        "category": "租金",
        "count": 10, // 該類別支出筆數
        "amount": 10000, // 該類別總金額(含其他費用)
      }
    ],
    "supplier": [
      {
        "supplierId": "10000000001",
        "supplierName": "蝦皮",
        "count": 10, // 該供應商支出筆數
        "amount": 10000, // 該供應商總金額(含其他費用)
      }
    ],
    "payer": [
      {
        "payerId": "10000000001",
        "payerName": "員工A",
        "advanceCount": 10, // 代墊筆數
        "advanceAmount": 10000, // 代墊總金額(含其他費用)
        "reimbursedAmount": 5000, // 已結清金額
        "pendingAmount": 5000 // 待結清金額
      }
    ]
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

| 狀態碼 | 錯誤碼   | 常數名稱                    | 說明                             |
| ------ | -------- | --------------------------- | -------------------------------- |
| 401    | E1002    | AuthTokenInvalid            | 無效的 accessToken，請重新登入   |
| 401    | E1003    | AuthTokenMissing            | accessToken 缺失，請重新登入     |
| 401    | E1004    | AuthTokenFormatError        | accessToken 格式錯誤，請重新登入 |
| 401    | E1005    | AuthStaffFailed             | 未找到有效的員工資訊，請重新登入 |
| 401    | E1006    | AuthContextMissing          | 未找到使用者認證資訊，請重新登入 |
| 403    | E1010    | AuthPermissionDenied        | 權限不足，無法執行此操作         |
| 400    | E2004    | ValTypeConversionFailed     | 參數類型轉換失敗                 |
| 400    | E3REP002 | ReportDateRangeExceed3Years | 日期範圍不能超過三年             |
| 500    | E9001    | SysInternalError            | 系統發生錯誤，請稍後再試         |
| 500    | E9002    | SysDatabaseError            | 資料庫操作失敗                   |

---

## 資料表

- `expenses`
- `suppliers`
- `staff_users`

---

## Service 邏輯

1. 查詢區段內的支出，使用 `category`、`supplier_id`、`payer_id` 分組。
2. 整理資料，並計算總和。
3. 回傳統計報表資料。

---

## 注意事項

- startDate 與 endDate 期限最長為 3 年。