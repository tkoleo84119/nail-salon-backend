## User Story

作為顧客，我希望能夠取得特定門市的服務項目，支援查詢條件，方便預約時選擇。

---

## Endpoint

**GET** `/api/stores/{storeId}/services`

---

## 說明

- 提供顧客查詢特定門市的服務項目。
- 支援基本查詢條件。
- 支援分頁（limit、offset）。
- 支援排序（sort）。
- 僅回傳啟用（`is_active=true`）且可見（`is_visible=true`）的服務。

---

## 權限

- 需要登入才可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameter

| 參數    | 說明   |
| ------- | ------ |
| storeId | 門市ID |

### Query Parameter

| 參數    | 型別   | 必填 | 預設值     | 說明                                             |
| ------- | ------ | ---- | ---------- | ------------------------------------------------ |
| isAddon | bool   | 否   |            | 是否為附加服務                                   |
| limit   | int    | 否   | 20         | 單頁筆數                                         |
| offset  | int    | 否   | 0          | 起始筆數                                         |
| sort    | string | 否   | created_at | 排序欄位 (可以逗號串接，有 `-` 表示 `DESC` 排序) |

### 驗證規則

| 欄位    | 必填 | 其他規則                                     |
| ------- | ---- | -------------------------------------------- |
| isAddon | 否   | <li>布林值                                   |
| limit   | 否   | <li>最小值1<li>最大值100                     |
| offset  | 否   | <li>最小值0<li>最大值1000000                 |
| sort    | 否   | <li>可以為 createdAt, updatedAt (其餘會忽略) |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "total": 10,
    "items": [
      {
        "id": "9000000001",
        "name": "手部單色",
        "price": 1000,
        "durationMinutes": 60,
        "note": "含基礎修型保養"
      },
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

| 狀態碼 | 錯誤碼 | 常數名稱               | 說明                             |
| ------ | ------ | ---------------------- | -------------------------------- |
| 401    | E1002  | AuthInvalidCredentials | 無效的 accessToken，請重新登入   |
| 401    | E1003  | AuthTokenMissing       | accessToken 缺失，請重新登入     |
| 401    | E1004  | AuthTokenFormatError   | accessToken 格式錯誤，請重新登入 |
| 401    | E1006  | AuthContextMissing     | 未找到使用者認證資訊，請重新登入 |
| 401    | E1011  | AuthCustomerFailed     | 未找到有效的顧客資訊，請重新登入 |
| 400    | E2023  | ValFieldMinNumber      | {field} 最小值為 {param}         |
| 400    | E2026  | ValFieldMaxNumber      | {field} 最大值為 {param}         |
| 400    | E2029  | ValFieldBoolean        | {field} 必須是布林值             |
| 500    | E9001  | SysInternalError       | 系統發生錯誤，請稍後再試         |
| 500    | E9002  | SysDatabaseError       | 資料庫操作失敗                   |

---

## 資料表

- `stores`
- `services`

---

## Service 邏輯

1. 驗證 `storeId` 是否存在。
2. 查詢該門市下 `is_visible=true` 且 `is_active=true` 的服務 (同時加上 `isAddon` 條件)。
3. 加入 `limit` 與 `offset` 處理分頁。
4. 加入 `sort` 處理排序。
5. 回傳結果與總筆數。

---

## 注意事項

- 僅回傳前台可見且啟用服務。
