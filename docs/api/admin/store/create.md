## User Story

作為一位管理員，我希望能新增門市（store），方便多據點管理。

---

## Endpoint

**POST** `/api/admin/stores`

---

## 說明

- 提供後台管理員新增門市功能。
- `ADMIN` 新增門市時會自動取得該門市權限。

---

## 權限

- 需要登入才可使用。
- 僅 `SUPER_ADMIN`、`ADMIN` 可操作。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Body 範例

```json
{
  "name": "大安旗艦店",
  "address": "台北市大安區復興南路一段100號",
  "phone": "02-12345678"
}
```

### 驗證規則

| 欄位    | 必填 | 其他規則                              |
| ------- | ---- | ------------------------------------- |
| name    | 是   | <li>最大長度100字元                   |
| address | 否   | <li>最大長度255字元                   |
| phone   | 否   | <li>支援台灣市話格式 <li>支援手機格式 |

---

## Response

### 成功 201 Created

```json
{
  "data": {
    "id": "8000000001",
    "name": "大安旗艦店",
    "address": "台北市大安區復興南路一段100號",
    "phone": "02-12345678",
    "isActive": true
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

| 狀態碼 | 錯誤碼   | 常數名稱                | 說明                                                                       |
| ------ | -------- | ----------------------- | -------------------------------------------------------------------------- |
| 401    | E1002    | AuthTokenInvalid        | 無效的 accessToken，請重新登入                                             |
| 401    | E1003    | AuthTokenMissing        | accessToken 缺失，請重新登入                                               |
| 401    | E1004    | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入                                           |
| 401    | E1005    | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入                                           |
| 401    | E1006    | AuthContextMissing      | 未找到使用者認證資訊，請重新登入                                           |
| 403    | E1010    | AuthPermissionDenied    | 權限不足，無法執行此操作                                                   |
| 400    | E2001    | ValJsonFormat           | JSON 格式錯誤，請檢查                                                      |
| 400    | E2020    | ValFieldRequired        | {field} 欄位為必填項目                                                     |
| 400    | E2024    | ValFieldStringMaxLength | {field} 長度最多只能有 {param} 個字元                                      |
| 400    | E2031    | ValFieldTaiwanPhone     | {field} 格式錯誤，請使用正確的台灣電話號碼格式 (0X-XXXXXXXX 或 09XXXXXXXX) |
| 409    | E3STO003 | StoreAlreadyExists      | 門市已存在，請創建其他門市                                                 |
| 500    | E9001    | SysInternalError        | 系統發生錯誤，請稍後再試                                                   |
| 500    | E9002    | SysDatabaseError        | 資料庫操作失敗                                                             |

---

## 資料表

- `stores`
- `staff_user_store_access`

---

## Service 邏輯

1. 檢查 `name` 是否已存在。
2. 建立 `stores` 資料。
3. 若為 `ADMIN` 則於 `staff_user_store_access` 關聯該門市。
4. 回傳新增結果。

---

## 注意事項

- 門市名稱不可重複。
- `ADMIN` 新增門市時會自動取得該門市權限。
