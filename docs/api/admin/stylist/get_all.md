## User Story

作為員工，我希望可以查詢某門市下的所有美甲師資料，並支援條件查詢與分頁，以便管理與選擇可排班對象。

---

## Endpoint

**GET** `/api/admin/stores/{storeId}/stylists`

---

## 說明

- 支援基本查詢條件。
- 支援分頁（limit、offset）。
- 支援排序（sort）。

---

## 權限

- 需要登入才可使用。
- 所有角色皆可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameters

| 參數    | 說明    |
| ------- | ------- |
| storeId | 門市 ID |

### Query Parameters

| 參數        | 型別   | 必填 | 預設值    | 說明                                             |
| ----------- | ------ | ---- | --------- | ------------------------------------------------ |
| name        | string | 否   |           | 模糊查詢姓名                                     |
| isIntrovert | bool   | 否   |           | 是否為內向者（I人）                              |
| isActive    | bool   | 否   |           | 是否為啟用員工                                   |
| limit       | int    | 否   | 20        | 單頁筆數                                         |
| offset      | int    | 否   | 0         | 起始筆數                                         |
| sort        | string | 否   | createdAt | 排序欄位 (可以逗號串接，有 `-` 表示 `DESC` 排序) |

### 驗證規則

| 欄位        | 必填 | 其他規則                                                        |
| ----------- | ---- | --------------------------------------------------------------- |
| name        | 否   | <li>最大長度100字元                                             |
| isIntrovert | 否   | <li>是否是布林值                                                |
| isActive    | 否   | <li>是否是布林值                                                |
| limit       | 否   | <li>最小值1<li>最大值100                                        |
| offset      | 否   | <li>最小值0<li>最大值1000000                                    |
| sort        | 否   | <li>可以為 createdAt, updatedAt, isIntrovert, name (其餘會忽略) |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "total": 2,
    "items": [
      {
        "id": "7000000001",
        "staffUserId": "6000000010",
        "name": "Ariel",
        "goodAtShapes": ["方形"],
        "goodAtColors": ["裸色系"],
        "goodAtStyles": ["簡約風"],
        "isIntrovert": false,
        "isActive": true
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

| 狀態碼 | 錯誤碼 | 常數名稱                | 說明                                  |
| ------ | ------ | ----------------------- | ------------------------------------- |
| 401    | E1002  | AuthTokenInvalid        | 無效的 accessToken，請重新登入        |
| 401    | E1003  | AuthTokenMissing        | accessToken 缺失，請重新登入          |
| 401    | E1004  | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入      |
| 401    | E1005  | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入      |
| 401    | E1006  | AuthContextMissing      | 未找到使用者認證資訊，請重新登入      |
| 403    | E3001  | AuthPermissionDenied    | 權限不足，無法執行此操作              |
| 400    | E2002  | AuthPathParamMissing    | 路徑參數缺失，請檢查                  |
| 400    | E2004  | AuthParamTypeConversion | 參數類型轉換失敗                      |
| 400    | E2023  | ValFieldMinNumber       | {field} 最小值為 {param}              |
| 400    | E2024  | ValFieldMaxLength       | {field} 長度最多只能有 {param} 個字元 |
| 400    | E2026  | ValFieldMaxNumber       | {field} 最大值為 {param}              |
| 400    | E2029  | ValFieldBoolean         | {field} 必須是布林值                  |
| 500    | E9001  | SysInternalError        | 系統發生錯誤，請稍後再試              |
| 500    | E9002  | SysDatabaseError        | 資料庫操作失敗                        |

---

## 資料表

- `stylists`
- `staff_users`
- `staff_user_store_access`
- `stores`

---

## Service 邏輯

1. 驗證員工是否擁有該門市存取權限。
2. 透過 `store_access` 表查詢與該門市有關聯的 `staff_user_id`。
3. JOIN `stylists` 表並依查詢條件過濾：
   - `name` 過濾
   - `is_introvert` 過濾
   - `staff_users.is_active = true` 過濾
5. 加入 `limit` / `offset` 處理分頁。
6. 加入 `sort` 處理排序。
7. 回傳總筆數與清單。

---

## 注意事項

- 預設會回傳所有員工。