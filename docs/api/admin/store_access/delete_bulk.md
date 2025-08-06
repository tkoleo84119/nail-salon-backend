## User Story

作為一位管理員，我希望可以刪除某位員工可操作的門市（多筆），以便控管其實際管理範圍。

---

## Endpoint

**DELETE** `/api/admin/staff/{staffId}/store-access/bulk`

---

## 說明

- 提供管理員刪除特定員工的門市存取權限。

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

### Body 範例

```json
{
  "storeIds": ["1", "3"]
}
```

### 驗證規則

| 欄位     | 必填 | 其他規則                | 說明                 |
| -------- | ---- | ----------------------- | -------------------- |
| storeIds | 是   | <li>最小1筆<li>最大20筆 | 欲移除的門市 ID 清單 |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "deleted": ["6000000011", "6000000012"]
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
| 403    | E3STA004 | StaffNotUpdateSelf      | 不可更新自己的帳號               |
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

1. 驗證目標員工是否存在
2. 不能刪除自己的 store access
3. 目標員工不能為 `SUPER_ADMIN`
4. 驗證門市是否存在
5. 檢查門市是否啟用中，且是否是該管理員有權限的門市
6. 執行刪除動作（從 `staff_user_store_access` 中移除多筆）
7. 回傳目前該員工剩餘可操作的門市清單

---

## 注意事項

- `SUPER_ADMIN` 不能被更動其門市權限
