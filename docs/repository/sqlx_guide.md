# SQLX Repository 開發指南

本文件詳細說明專案中 SQLX Repository 層的開發規範、最佳實務和撰寫模式。

## 概覽

SQLX Repository 用於處理需要動態 SQL 查詢的場景，主要包括：
- **動態篩選查詢** (GET_ALL 操作)
- **部分欄位更新** (UPDATE 操作) 
- **複雜條件查詢** (多表 JOIN 查詢)
- **交易操作** (需要原子性的操作)

## 檔案結構

```
internal/repository/sqlx/
├── repositories.go          # Repository 整合器
├── {domain}.go              # 各領域的 Repository 實作
├── booking.go
├── customer.go
├── service.go
├── store.go
├── stylist.go
├── time_slot.go
└── ...
```

## 基本結構模式

### Repository 結構定義
```go
type {Domain}Repository struct {
    db *sqlx.DB
}

func New{Domain}Repository(db *sqlx.DB) *{Domain}Repository {
    return &{Domain}Repository{
        db: db,
    }
}
```

### Repository 整合器
所有 SQLX Repository 都透過 `repositories.go` 統一管理：

```go
type Repositories struct {
    Booking       *BookingRepository
    Customer      *CustomerRepository
    Service       *ServiceRepository
    // ... 其他 Repository
}

func NewRepositories(db *sqlx.DB) *Repositories {
    return &Repositories{
        Booking:  NewBookingRepository(db),
        Customer: NewCustomerRepository(db),
        // ... 初始化其他 Repository
    }
}
```

## 開發模式

### 1. GET_ALL 動態查詢模式

用於需要多種篩選條件的查詢操作。

#### 參數結構定義
```go
type GetAll{Domain}sByFilterParams struct {
    // 篩選參數 (使用指標型別，支援 nil 檢查)
    Name      *string
    Status    *string
    IsActive  *bool
    StartDate *time.Time
    EndDate   *time.Time
    
    // 分頁參數
    Limit  *int
    Offset *int
    Sort   *[]string
}
```

#### 返回結構定義
```go
type GetAll{Domain}sByFilterItem struct {
    ID        int64              `db:"id"`
    Name      string             `db:"name"`
    Status    string             `db:"status"`
    IsActive  pgtype.Bool        `db:"is_active"`
    CreatedAt pgtype.Timestamptz `db:"created_at"`
    UpdatedAt pgtype.Timestamptz `db:"updated_at"`
}
```

#### 實作模式
```go
func (r *{Domain}Repository) GetAll{Domain}sByFilter(ctx context.Context, params GetAll{Domain}sByFilterParams) (int, []GetAll{Domain}sByFilterItem, error) {
    // 1. 動態 WHERE 條件組裝
    whereConditions := []string{}
    args := []interface{}{}
    
    if params.Name != nil && *params.Name != "" {
        whereConditions = append(whereConditions, fmt.Sprintf("name ILIKE $%d", len(args)+1))
        args = append(args, "%"+*params.Name+"%")
    }
    
    if params.IsActive != nil {
        whereConditions = append(whereConditions, fmt.Sprintf("is_active = $%d", len(args)+1))
        args = append(args, *params.IsActive)
    }
    
    whereClause := ""
    if len(whereConditions) > 0 {
        whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
    }
    
    // 2. 計數查詢 (先確認是否有資料)
    countQuery := fmt.Sprintf(`
        SELECT COUNT(*)
        FROM {table_name}
        %s
    `, whereClause)
    
    var total int
    if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
        return 0, nil, fmt.Errorf("count failed: %w", err)
    }
    if total == 0 {
        return 0, []GetAll{Domain}sByFilterItem{}, nil
    }
    
    // 3. 分頁與排序處理
    limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
    defaultSortArr := []string{"created_at DESC"}
    sort := utils.HandleSortByMap(map[string]string{
        "createdAt": "created_at",
        "updatedAt": "updated_at",
        "name":      "name",
    }, defaultSortArr, params.Sort)
    
    args = append(args, limit, offset)
    limitIndex := len(args) - 1
    offsetIndex := len(args)
    
    // 4. 資料查詢
    query := fmt.Sprintf(`
        SELECT id, name, status, is_active, created_at, updated_at
        FROM {table_name}
        %s
        ORDER BY %s
        LIMIT $%d OFFSET $%d
    `, whereClause, sort, limitIndex, offsetIndex)
    
    var results []GetAll{Domain}sByFilterItem
    if err := r.db.SelectContext(ctx, &results, query, args...); err != nil {
        return 0, nil, fmt.Errorf("query failed: %w", err)
    }
    
    return total, results, nil
}
```

### 2. UPDATE 部分欄位更新模式

用於只更新有變更的欄位，避免不必要的資料庫寫入。

#### 參數結構定義
```go
type Update{Domain}Params struct {
    // 所有欄位都使用指標型別
    Name      *string
    Status    *string
    IsActive  *bool
    Note      *string
}
```

#### 返回結構定義
```go
type Update{Domain}Response struct {
    ID        int64              `db:"id"`
    Name      string             `db:"name"`
    Status    string             `db:"status"`
    IsActive  pgtype.Bool        `db:"is_active"`
    Note      pgtype.Text        `db:"note"`
    CreatedAt pgtype.Timestamptz `db:"created_at"`
    UpdatedAt pgtype.Timestamptz `db:"updated_at"`
}
```

#### 實作模式
```go
func (r *{Domain}Repository) Update{Domain}(ctx context.Context, id int64, params Update{Domain}Params) (Update{Domain}Response, error) {
    // 1. 動態 SET 條件組裝
    setParts := []string{"updated_at = NOW()"}  // 總是更新 updated_at
    args := []interface{}{}
    
    if params.Name != nil && *params.Name != "" {
        setParts = append(setParts, fmt.Sprintf("name = $%d", len(args)+1))
        args = append(args, *params.Name)
    }
    
    if params.Status != nil {
        setParts = append(setParts, fmt.Sprintf("status = $%d", len(args)+1))
        args = append(args, *params.Status)
    }
    
    if params.IsActive != nil {
        setParts = append(setParts, fmt.Sprintf("is_active = $%d", len(args)+1))
        args = append(args, *params.IsActive)
    }
    
    if params.Note != nil {
        setParts = append(setParts, fmt.Sprintf("note = $%d", len(args)+1))
        args = append(args, *params.Note)
    }
    
    // 2. 檢查是否有欄位需要更新
    if len(setParts) == 1 {
        return Update{Domain}Response{}, fmt.Errorf("no fields to update")
    }
    
    // 3. 新增 WHERE 條件
    args = append(args, id)
    
    // 4. 執行更新並返回結果
    query := fmt.Sprintf(`
        UPDATE {table_name}
        SET %s
        WHERE id = $%d
        RETURNING id, name, status, is_active, note, created_at, updated_at
    `, strings.Join(setParts, ", "), len(args))
    
    var result Update{Domain}Response
    if err := r.db.GetContext(ctx, &result, query, args...); err != nil {
        return Update{Domain}Response{}, fmt.Errorf("update failed: %w", err)
    }
    
    return result, nil
}
```

### 3. 交易操作模式

用於需要原子性操作的場景，通常搭配 PGX Pool Transaction。

```go
func (r *{Domain}Repository) {Operation}Tx(ctx context.Context, tx *sqlx.Tx, params {Operation}TxParams) error {
    query := `
        UPDATE {table_name}
        SET status = $1,
            updated_at = NOW()
        WHERE id = $2
    `
    
    args := []interface{}{
        params.Status,
        params.ID,
    }
    
    _, err := tx.ExecContext(ctx, query, args...)
    if err != nil {
        return fmt.Errorf("{operation} failed: %w", err)
    }
    
    return nil
}
```

## 資料型別處理

### pgtype 型別對應

專案使用 `pgtype` 處理 PostgreSQL 特定型別：

| PostgreSQL 型別 | Go 型別 | pgtype 型別 |
|------------------|---------|-------------|
| `text` (nullable) | `*string` | `pgtype.Text` |
| `boolean` (nullable) | `*bool` | `pgtype.Bool` |
| `timestamptz` | `time.Time` | `pgtype.Timestamptz` |
| `date` | `time.Time` | `pgtype.Date` |
| `time` | `time.Time` | `pgtype.Time` |
| `numeric` | `decimal.Decimal` | `pgtype.Numeric` |
| `text[]` | `[]string` | `[]string` + `pgtype.NewMap().SQLScanner()` |

### 陣列型別處理

PostgreSQL 陣列欄位需要特殊處理：

```go
// 查詢時處理陣列
query := `
    SELECT 
        id,
        name,
        COALESCE(favorite_colors, '{}'::text[]) AS favorite_colors
    FROM customers
`

// 掃描陣列欄位
rows, err := r.db.QueryContext(ctx, query, args...)
defer rows.Close()

m := pgtype.NewMap()
var customer Customer
for rows.Next() {
    err := rows.Scan(
        &customer.ID,
        &customer.Name,
        m.SQLScanner(&customer.FavoriteColors),  // 使用 SQLScanner 處理陣列
    )
}
```

### 工具函式使用

專案提供了工具函式協助型別轉換：

```go
// 字串指標轉 pgtype.Text
args = append(args, utils.StringPtrToPgText(params.Address, false))

// 布林指標轉 pgtype.Bool  
args = append(args, utils.BoolPtrToPgBool(params.IsActive))

// 時間字串轉 pgtype.Time
startTime, err := utils.TimeStringToPgTime(*params.StartTime)
if err != nil {
    return UpdateTimeSlotResponse{}, fmt.Errorf("convert start time failed: %w", err)
}
args = append(args, startTime)
```

## 錯誤處理

### 統一錯誤格式
所有錯誤都使用 `fmt.Errorf` 包裝，提供清楚的錯誤訊息：

```go
if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
    return 0, nil, fmt.Errorf("count {domain} failed: %w", err)
}

if err := r.db.SelectContext(ctx, &results, query, args...); err != nil {
    return 0, nil, fmt.Errorf("query {domain} failed: %w", err)
}
```

### 資料驗證
```go
// 檢查更新參數
if len(setParts) == 1 {
    return Update{Domain}Response{}, fmt.Errorf("no fields to update")
}

// 檢查計數結果
if total == 0 {
    return 0, []GetAll{Domain}sByFilterItem{}, nil
}
```

## 效能考量

### 1. 先計數後查詢
避免不必要的複雜查詢：
```go
// 先執行簡單的計數查詢
if total == 0 {
    return 0, []Items{}, nil  // 直接返回空結果
}
// 只有在有資料時才執行複雜的 JOIN 查詢
```

### 2. 索引友善查詢
WHERE 條件順序要符合資料庫索引設計：
```go
// 按索引順序組織 WHERE 條件
whereParts := []string{"store_id = $1"}  // 主要索引欄位優先
if params.Status != nil {
    whereParts = append(whereParts, fmt.Sprintf("status = $%d", len(args)+1))
}
```

### 3. 分頁最佳化
使用工具函式設定合理的分頁預設值：
```go
limit, offset := utils.SetDefaultValuesOfPagination(params.Limit, params.Offset, 20, 0)
```

## 測試考量

### Repository 層測試重點
1. **動態查詢組裝**：驗證不同參數組合的 SQL 正確性
2. **部分更新邏輯**：確保只更新有變更的欄位
3. **分頁與排序**：測試邊界條件和預設值
4. **錯誤處理**：資料庫錯誤的正確傳播

### 測試資料準備
```go
func setupTestData(t *testing.T, db *sqlx.DB) {
    // 準備測試資料
    // 清理測試資料
}
```

## 開發檢查清單

- [ ] Repository 結構體使用 `*sqlx.DB`
- [ ] 函式命名遵循 `{Operation}{Domain}[sByFilter|Tx]` 格式
- [ ] 參數結構體所有欄位使用指標型別
- [ ] 實作動態 WHERE/SET 條件組裝
- [ ] 使用 `fmt.Sprintf` 和參數計數器避免 SQL 注入
- [ ] 先計數後查詢的效能模式
- [ ] 正確處理 pgtype 型別轉換
- [ ] 統一錯誤格式和包裝
- [ ] 交易函式使用 `*sqlx.Tx` 參數
- [ ] 使用工具函式處理分頁和型別轉換

---

**注意事項**：
- SQLX Repository 主要用於動態查詢，簡單 CRUD 操作建議使用 SQLC
- 所有查詢都必須使用參數化查詢避免 SQL 注入
- 複雜查詢建議先在資料庫工具中測試 SQL 語法
- 交易操作必須在 Service 層管理交易生命週期