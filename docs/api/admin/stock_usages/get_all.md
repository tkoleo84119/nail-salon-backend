## User Story

作為員工，我希望可以查詢某門市下的所有產品庫存使用紀錄，並支援條件查詢與分頁，以利管理與設定產品庫存使用紀錄。

---

## Endpoint

**GET** `/api/admin/stores/{storeId}/stock-usages`

---

## 說明

- 提供員工查詢某門市下的所有產品庫存使用紀錄資料。
- 支援基本查詢條件。
- 支援分頁（limit、offset）。
- 支援排序（sort）。

---

## 權限

- 需要登入才可使用。
- 僅 `SUPER_ADMIN`、`ADMIN` `MANAGER` 可操作。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Query Parameters

| 參數      | 型別   | 必填 | 預設值     | 說明                                             |
| --------- | ------ | ---- | ---------- | ------------------------------------------------ |
| productId | string | 否   |            | 查詢指定產品ID                                   |
| name      | string | 否   |            | 模糊查詢產品名稱                                 |
| isInUse   | bool   | 否   |            | 查詢是否正在使用                                 |
| limit     | int    | 否   | 20         | 單頁筆數                                         |
| offset    | int    | 否   | 0          | 起始筆數                                         |
| sort      | string | 否   | -createdAt | 排序欄位 (可以逗號串接，有 `-` 表示 `DESC` 排序) |

### 驗證規則

| 欄位      | 必填 | 其他規則                                              |
| --------- | ---- | ----------------------------------------------------- |
| productId | 否   |                                                       |
| name      | 否   | <li>不能為空字串<li>最大長度100字元                   |
| isInUse   | 否   |                                                       |
| limit     | 否   | <li>最小值1<li>最大值100                              |
| offset    | 否   | <li>最小值0<li>最大值1000000                          |
| sort      | 否   | <li>可以為 createdAt, updatedAt, isInUse (其餘會忽略) |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "total": 3,
    "items": [
      {
        "id": "9000000001",
        "product": {
          "id": "9000000001",
          "name": "利德瑪"
        },
        "quantity": 10,
        "isInUse": true,
        "expiration": "2027-01-01", // ""表示無期限
        "usageStarted": "2025-01-01",
        "usageEndedAt": "2025-01-01", // ""表示正在使用
        "createdAt": "2025-01-01T00:00:00+08:00",
        "updatedAt": "2025-01-01T00:00:00+08:00"
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

| 狀態碼 | 錯誤碼 | 常數名稱                | 說明                                  |
| ------ | ------ | ----------------------- | ------------------------------------- |
| 401    | E1002  | AuthTokenInvalid        | 無效的 accessToken，請重新登入        |
| 401    | E1003  | AuthTokenMissing        | accessToken 缺失，請重新登入          |
| 401    | E1004  | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入      |
| 401    | E1005  | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入      |
| 401    | E1006  | AuthContextMissing      | 未找到使用者認證資訊，請重新登入      |
| 403    | E1010  | AuthPermissionDenied    | 權限不足，無法執行此操作              |
| 400    | E2004  | ValTypeConversionFailed | 參數類型轉換失敗                      |
| 400    | E2023  | ValFieldMinNumber       | {field} 最小值為 {param}              |
| 400    | E2024  | ValFieldStringMaxLength | {field} 長度最多只能有 {param} 個字元 |
| 400    | E2026  | ValFieldMaxNumber       | {field} 最大值為 {param}              |
| 400    | E2036  | ValFieldNoBlank         | {field} 不能為空字串                  |
| 500    | E9001  | SysInternalError        | 系統發生錯誤，請稍後再試              |
| 500    | E9002  | SysDatabaseError        | 資料庫操作失敗                        |

---

## 資料表

- `products`
- `stock_usages`

---

## Service 邏輯

1. 確認門市存取權限。
2. 根據 `productId`（產品ID）、`name`（名稱）、`isInUse`（是否正在使用）條件動態查詢。
3. 加入 `limit` 與 `offset` 處理分頁。
4. 加入 `sort` 處理排序。
5. 回傳總筆數與項目清單。

---

## 注意事項

- createdAt 與 updatedAt 會是標準 Iso 8601 格式。
