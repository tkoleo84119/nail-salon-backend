## User Story

作為員工，我希望可以查詢系統活動紀錄，以便了解系統內的各種操作記錄和用戶行為。

---

## Endpoint

**GET** `/api/admin/activity-logs`

---

## 說明

- 提供員工查詢系統活動紀錄。
- 支援分頁（limit）。
- 記錄會依時間倒序排列（最新的在前）。
- 系統會自動保留最新 50 筆記錄，每筆記錄保存 24 小時。

---

## 權限

- 需要登入才可使用。
- 所有員工都可操作。

---

## Request

### Header

- Content-Type: application/json
- Authorization: Bearer <access_token>

### Query Parameters

| 參數  | 型別 | 必填 | 預設值 | 說明     |
| ----- | ---- | ---- | ------ | -------- |
| limit | int  | 否   | 20     | 單頁筆數 |

### 驗證規則

| 欄位  | 必填 | 其他規則                |
| ----- | ---- | ----------------------- |
| limit | 否   | <li>最小值1<li>最大值50 |

---

## Response

### 成功 200 OK

```json
{
  "data": {
    "total": 4,
    "items": [
      {
        "id": "1735689600001",
        "type": "CUSTOMER_REGISTER",
        "message": "顧客 王小美 完成註冊",
        "timestamp": "2025-01-01T10:30:00+08:00"
      },
      {
        "id": "1735689500002",
        "type": "ADMIN_BOOKING_CREATE",
        "message": "員工 張美甲師 為顧客 李小花 建立預約",
        "timestamp": "2025-01-01T10:28:20+08:00"
      },
      {
        "id": "1735689400003",
        "type": "CUSTOMER_BOOKING_UPDATE",
        "message": "顧客 陳小雅 修改預約",
        "timestamp": "2025-01-01T10:26:40+08:00"
      },
      {
        "id": "1735689300004",
        "type": "ADMIN_BOOKING_COMPLETED",
        "message": "員工 林經理 為預約完成結帳",
        "timestamp": "2025-01-01T10:25:00+08:00"
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

| 狀態碼 | 錯誤碼 | 常數名稱                | 說明                             |
| ------ | ------ | ----------------------- | -------------------------------- |
| 401    | E1002  | AuthTokenInvalid        | 無效的 accessToken，請重新登入   |
| 401    | E1003  | AuthTokenMissing        | accessToken 缺失，請重新登入     |
| 401    | E1004  | AuthTokenFormatError    | accessToken 格式錯誤，請重新登入 |
| 401    | E1005  | AuthStaffFailed         | 未找到有效的員工資訊，請重新登入 |
| 401    | E1006  | AuthContextMissing      | 未找到使用者認證資訊，請重新登入 |
| 403    | E1010  | AuthPermissionDenied    | 權限不足，無法執行此操作         |
| 400    | E2004  | ValTypeConversionFailed | 參數類型轉換失敗                 |
| 400    | E2023  | ValFieldMinNumber       | {field} 最小值為 {param}         |
| 400    | E2026  | ValFieldMaxNumber       | {field} 最大值為 {param}         |
| 500    | E9001  | SysInternalError        | 系統發生錯誤，請稍後再試         |
| 500    | E9003  | SysCacheError           | 快取操作失敗                     |

---

## 資料來源

- Redis List：`activity_logs`

---

## Service 邏輯

1. 從 Redis List 中取得指定數量的活動紀錄。
2. 記錄按時間倒序排列（最新的在前）。
3. 回傳項目清單。

---

## 活動類型說明

系統會自動追蹤以下 8 種活動類型：

| 活動類型                  | 說明               | 觸發時機           |
| ------------------------- | ------------------ | ------------------ |
| `CUSTOMER_REGISTER`       | 顧客註冊           | 顧客完成 LINE 註冊 |
| `CUSTOMER_BOOKING`        | 顧客建立預約       | 顧客自行建立預約   |
| `CUSTOMER_BOOKING_UPDATE` | 顧客修改預約       | 顧客自行修改預約   |
| `CUSTOMER_BOOKING_CANCEL` | 顧客取消預約       | 顧客自行取消預約   |
| `ADMIN_BOOKING_CREATE`    | 管理員協助建立預約 | 員工代客戶建立預約 |
| `ADMIN_BOOKING_UPDATE`    | 管理員協助修改預約 | 員工代客戶修改預約 |
| `ADMIN_BOOKING_CANCEL`    | 管理員協助取消預約 | 員工代客戶取消預約 |
| `ADMIN_BOOKING_COMPLETED` | 管理員完成預約結帳 | 員工完成結帳流程   |

---

## 注意事項

- timestamp 會是標準 ISO 8601 格式。
- 系統只保留最新 50 筆記錄，超過的記錄會自動移除。
- 每筆記錄會在 24 小時後自動過期。
- 記錄的 message 欄位會以「誰做了什麼」的格式呈現。
- 由於使用 Redis 儲存，此 API 不支援複雜的查詢條件，僅支援基本的數量限制。