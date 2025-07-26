## User Story

作為管理員，我希望可以查詢特定員工的門市存取權限（Store Access），以便了解其可操作的門市範圍。

---

## Endpoint

**GET** `/api/admin/staff/{staffId}/store-access`

---

## 說明

- 僅限 `admin` 以上角色可查詢。
- 回傳該員工可操作的門市列表。

---

## 權限

- 僅限 `admin` 以上角色（`SUPER_ADMIN`, `ADMIN`）

---

## Request

### Header

Authorization: Bearer <access_token>

### Path Parameter

| 參數    | 說明        |
| ------- | ----------- |
| staffId | 員工帳號 ID |

---

## Response

### 成功 200 OK

```json
{
  "data": [
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
```

### 失敗

#### 401 Unauthorized - 未登入/Token失效

```json
{
  "message": "無效的 accessToken"
}
```

#### 403 Forbidden - 權限不足

```json
{
  "message": "無權限存取此資源"
}
```

#### 404 Not Found - 查無員工資料

```json
{
  "message": "查無此員工帳號"
}
```

#### 500 Internal Server Error

```json
{
  "message": "系統發生錯誤，請稍後再試"
}
```

---

## 資料表

- `staff_users`
- `staff_user_store_access`
- `stores`

---

## Service 邏輯

1. 確認 `staffId` 是否存在。
2. 查詢 `staff_user_store_access` 表取得所有該員工可操作的門市 ID。
   - 如果是 `SUPER_ADMIN` 角色，則回傳所有門市。
3. 回傳門市清單。

---

## 注意事項

- 僅回傳 `stores.is_active=true` 的門市。
- 若無授權任何門市，回傳空陣列。

