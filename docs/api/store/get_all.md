## User Story

作為顧客，我希望能夠取得所有門市的資料，支援分頁，方便快速查找可預約據點。

---

## Endpoint

**GET** `/api/stores`

---

## 說明

- 提供顧客查詢所有門市資料。
- 支援分頁（limit、offset）。
- 支援排序（sort）。
- 僅回傳 `is_active=true` 的門市。

---

## 權限

- 需要登入才可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Query Parameter

| 參數   | 型別   | 必填 | 預設值     | 說明                                             |
| ------ | ------ | ---- | ---------- | ------------------------------------------------ |
| limit  | int    | 否   | 20         | 單頁筆數                                         |
| offset | int    | 否   | 0          | 起始筆數                                         |
| sort   | string | 否   | created_at | 排序欄位 (可以逗號串接，有 `-` 表示 `DESC` 排序) |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "total": 10,
    "items": [
      {
        "id": "8000000001",
        "name": "大安旗艦店",
        "address": "台北市大安區復興南路一段100號",
        "phone": "02-1234-5678"
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
| 500    | E9001  | SysInternalError       | 系統發生錯誤，請稍後再試         |
| 500    | E9002  | SysDatabaseError       | 資料庫操作失敗                   |

---

## 資料表

- `stores`

---

## Service 邏輯

1. 查詢 `is_active=true` 的門市。
2. 加入 `limit` 與 `offset` 處理分頁。
3. 加入 `sort` 處理排序。
4. 回傳結果與總筆數。

---

## 注意事項

- 僅回傳啟用門市。
