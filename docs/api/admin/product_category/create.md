## User Story

作為一位管理員，我希望能新增產品分類，方便維護分類資訊。

---

## Endpointˆ

**POST** `/api/admin/product-categories`

---

## 說明

- 提供後台管理員新增產品分類功能。
- 分類名稱須唯一。

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
  "name": "新客優惠-八折"
}
```

### 驗證規則

| 欄位 | 必填 | 其他規則                            | 說明     |
| ---- | ---- | ----------------------------------- | -------- |
| name | 是   | <li>不能為空字串<li>最大長度100字元 | 分類名稱 |

---

## Response

### 成功 201 Created

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

| 狀態碼 | 錯誤碼  | 常數名稱                  | 說明                                  |
| ------ | ------- | ------------------------- | ------------------------------------- |
| 401    | E1002   | AuthTokenInvalid          | 無效的 accessToken，請重新登入        |
| 401    | E1003   | AuthTokenMissing          | accessToken 缺失，請重新登入          |
| 401    | E1004   | AuthTokenFormatError      | accessToken 格式錯誤，請重新登入      |
| 401    | E1005   | AuthStaffFailed           | 未找到有效的員工資訊，請重新登入      |
| 401    | E1006   | AuthContextMissing        | 未找到使用者認證資訊，請重新登入      |
| 403    | E1010   | AuthPermissionDenied      | 權限不足，無法執行此操作              |
| 400    | E2001   | ValJsonFormat             | JSON 格式錯誤，請檢查                 |
| 400    | E2020   | ValFieldRequired          | {field} 為必填項目                    |
| 400    | E2024   | ValFieldStringMaxLength   | {field} 長度最多只能有 {param} 個字元 |
| 400    | E2036   | ValFieldNoBlank           | {field} 不能為空字串                  |
| 409    | E3PC001 | CategoryNameAlreadyExists | 分類名稱已存在，請使用其他名稱        |
| 500    | E9001   | SysInternalError          | 系統發生錯誤，請稍後再試              |
| 500    | E9002   | SysDatabaseError          | 資料庫操作失敗                        |

---

## 資料表

- `product_categories`

---

## Service 邏輯

1. 確認 `name` 是否唯一。
2. 建立 `product_categories` 資料。
4. 回傳新增結果。

---

## 注意事項

- 分類名稱不可重複。
