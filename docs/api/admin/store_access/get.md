## User Story

作為管理員，我希望可以查詢特定員工的門市存取權限（Store Access），以便了解其可操作的門市範圍。

---

## Endpoint

**GET** `/api/admin/staff/{staffId}/store-access`

---

## 說明

- 提供管理員查詢特定員工的門市存取權限。

---

## 權限

- 需要登入才可使用。
- `SUPER_ADMIN` 與 `ADMIN` 可使用。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Path Parameter

| 參數    | 說明        |
| ------- | ----------- |
| staffId | 員工帳號 ID |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "storeList": [
      {
        "storeId": "8000000001",
        "name": "大安旗艦店"
      },
      {
        "storeId": "8000000003",
        "name": "信義分店"
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

| 狀態碼 | 錯誤碼   | 常數名稱                | 說明                             |
| ------ | -------- | ----------------------- | -------------------------------- |
| 401    | E1002    | AuthInvalidCredentials  | 無效的 accessToken，請重新登入   |
| 401    | E1003    | AuthTokenMissing        | accessToken 缺失，請重新登入     |
| 401    | E1004    | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入 |
| 401    | E1005    | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入 |
| 401    | E1006    | AuthContextMissing      | 未找到使用者認證資訊，請重新登入 |
| 403    | E1010    | AuthPermissionDenied    | 權限不足，無法執行此操作         |
| 400    | E2002    | ValPathParamMissing     | 路徑參數缺失，請檢查             |
| 400    | E2004    | ValTypeConversionFailed | 參數類型轉換失敗                 |
| 404    | E3STA005 | StaffNotFound           | 員工帳號不存在                   |
| 500    | E9001    | SysInternalError        | 系統發生錯誤，請稍後再試         |
| 500    | E9002    | SysDatabaseError        | 資料庫操作失敗                   |

---

## 資料表

- `staff_users`
- `staff_user_store_access`
- `stores`

---

## Service 邏輯

1. 確認 `staffId` 是否存在。
2. 查詢 `staff_user_store_access` 表取得所有該員工可操作的門市資料。
3. 回傳門市清單。

---

## 注意事項

- 若無授權任何門市，回傳空陣列。
