Great question! Here's what to reuse and what to skip when building your product-service:

## 🔄 REUSE (Copy & Adapt)
□ Copy jsonlog/, validator/ folders as-is
□ Copy errors.go, healthcheck.go, helpers.go, server.go as-is
□ Copy main.go → remove SMTP, change port, change DB DSN env var
□ Copy middleware.go → adapt authenticate() to call user-service
□ Copy routes.go → replace with service-specific routes
□ Copy Dockerfile → change binary name
□ Copy docker-compose.yaml → change service name, port, DB name
□ Copy go.mod → change module path, remove mail/crypto
□ Copy .gitignore, Makefile, .env.example → adapt names
□ Create new models in internal/data/
□ Create new handlers in cmd/api/
□ Create new migrations/
□ Write new README.md
### 1. **Entire Package Structure** ✅
```
product-service/
├── cmd/api/
│   ├── context.go          ✅ REUSE (change User → Product where needed)
│   ├── errors.go           ✅ REUSE (identical)
│   ├── healthcheck.go      ✅ REUSE (identical)
│   ├── helpers.go          ✅ REUSE (identical)
│   ├── main.go             ✅ REUSE & ADAPT
│   ├── middleware.go       ✅ REUSE & ADAPT
│   ├── routes.go           ✅ REUSE & ADAPT
│   └── server.go           ✅ REUSE (identical)
├── internal/
│   ├── data/
│   │   ├── models.go       ✅ REUSE & ADAPT
│   │   └── products.go     ⚠️ NEW (your product logic)
│   ├── jsonlog/            ✅ REUSE (identical - copy entire folder)
│   └── validator/          ✅ REUSE (identical - copy entire folder)
├── migrations/             ⚠️ NEW (product-specific schemas)
├── Dockerfile              ✅ REUSE & ADAPT
├── docker-compose.yaml     ✅ REUSE & ADAPT
├── go.mod                  ✅ REUSE & ADAPT
├── .gitignore              ✅ REUSE (identical)
├── .env.example            ✅ REUSE & ADAPT
├── Makefile                ✅ REUSE & ADAPT
└── README.md               ✅ REUSE & ADAPT
```

### 2. **Copy 100% As-Is (No Changes)**
- `internal/jsonlog/` - entire folder
- `internal/validator/` - entire folder
- `cmd/api/errors.go` - all error handlers
- `cmd/api/healthcheck.go` - health check logic
- `cmd/api/helpers.go` - JSON read/write, envelope, background
- `cmd/api/server.go` - graceful shutdown logic
- `.gitignore`

### 3. **Copy & Adapt (Simple Find/Replace)**

#### **Dockerfile** - Just change binary name
```dockerfile
# Find: user-service
# Replace: product-service
RUN CGO_ENABLED=0 go build -o product-service ./cmd/api
COPY --from=builder /app/product-service .
```

#### **go.mod** - Change module path & remove unused deps
```go
module github.com/PaulBabatuyi/FashionMarket-Backend/product-service

// REMOVE these (user-service specific):
// - github.com/go-mail/mail/v2 (no emails in product service)
// - golang.org/x/crypto (no password hashing)

// KEEP these:
// - github.com/go-chi/chi/v5
// - github.com/lib/pq
// - golang.org/x/time (rate limiting)
// - github.com/felixge/httpsnoop (metrics)
```

#### **cmd/api/main.go** - Keep structure, remove SMTP config
```go
type config struct {
    port int
    env  string
    db   struct {
        dsn          string
        maxOpenConns int
        maxIdleConns int
        maxIdleTime  string
    }
    limiter struct {
        rps     float64
        burst   int
        enabled bool
    }
    // ❌ REMOVE smtp struct - no emails needed
    cors struct {
        trustedOrigins []string
    }
}

type application struct {
    config config
    logger *jsonlog.Logger
    models data.Models
    // ❌ REMOVE mailer - no emails
    wg     sync.WaitGroup
}

func main() {
    var cfg config
    flag.IntVar(&cfg.port, "port", 5000, "API server port") // Different port!
    flag.StringVar(&cfg.env, "env", "development", "Environment")
    flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("FASHIONPRODUCTS_DB_DSN"), "PostgreSQL DSN")
    
    // ... keep db pool flags
    // ... keep rate limiter flags
    // ... keep CORS flags
    // ❌ REMOVE all smtp flags
    
    flag.Parse()
    
    logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
    db, err := openDB(cfg)
    // ... same DB setup
    
    // ... same expvar metrics
    
    app := &application{
        config: cfg,
        logger: logger,
        models: data.NewModels(db),
        // ❌ REMOVE mailer initialization
    }
    
    err = app.serve()
    // ... same shutdown
}

// openDB function - identical, keep as-is
```

#### **cmd/api/middleware.go** - Keep most, simplify authenticate
```go
// ✅ KEEP: recoverPanic (identical)
// ✅ KEEP: rateLimit (identical)
// ✅ KEEP: enableCORS (identical)
// ✅ KEEP: metrics (identical)

// ⚠️ ADAPT: authenticate - call user-service instead
func (app *application) authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Vary", "Authorization")
        
        authorizationHeader := r.Header.Get("Authorization")
        if authorizationHeader == "" {
            r = app.contextSetUser(r, data.AnonymousUser)
            next.ServeHTTP(w, r)
            return
        }
        
        headerParts := strings.Split(authorizationHeader, " ")
        if len(headerParts) != 2 || headerParts[0] != "Bearer" {
            app.invalidAuthenticationTokenResponse(w, r)
            return
        }
        
        token := headerParts[1]
        
        // 🔥 INSTEAD OF: checking tokens table locally
        // DO THIS: Call user-service to validate token
        user, err := app.getUserFromUserService(token) // HTTP call to user-service
        if err != nil {
            app.invalidAuthenticationTokenResponse(w, r)
            return
        }
        
        r = app.contextSetUser(r, user)
        next.ServeHTTP(w, r)
    })
}

// ⚠️ NEW: Add HTTP client method
func (app *application) getUserFromUserService(token string) (*data.User, error) {
    // Make HTTP call: GET http://user-service:4000/v1/auth/validate
    // With Authorization: Bearer {token}
    // Return user info
}

// ✅ KEEP: requireAuthenticatedUser (identical)
// ✅ KEEP: requireActivatedUser (identical)
// ⚠️ ADAPT: requirePermission (check "products:write" instead of "movies:read")
```

#### **cmd/api/routes.go** - Completely different routes
```go
func (app *application) routes() http.Handler {
    router := chi.NewRouter()
    
    router.NotFound(http.HandlerFunc(app.notFoundResponse))
    router.MethodNotAllowed(http.HandlerFunc(app.methodNotAllowedResponse))
    
    router.MethodFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
    
    // 🔥 PRODUCT-SPECIFIC ROUTES
    router.MethodFunc(http.MethodGet, "/v1/products", app.listProductsHandler)
    router.MethodFunc(http.MethodPost, "/v1/products", 
        app.requirePermission("products:write", app.createProductHandler))
    router.MethodFunc(http.MethodGet, "/v1/products/{id}", app.getProductHandler)
    router.MethodFunc(http.MethodPatch, "/v1/products/{id}", 
        app.requirePermission("products:write", app.updateProductHandler))
    router.MethodFunc(http.MethodDelete, "/v1/products/{id}", 
        app.requirePermission("products:write", app.deleteProductHandler))
    
    router.Method(http.MethodGet, "/debug/vars", expvar.Handler())
    
    return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
```

#### **internal/data/models.go** - Adapt model names
```go
package data

import (
    "database/sql"
    "errors"
)

var (
    ErrRecordNotFound = errors.New("record not found")
    ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
    Products ProductModel  // Changed from Users, Tokens, Permissions
    // Add more if needed: Categories, Reviews, etc.
}

func NewModels(db *sql.DB) Models {
    return Models{
        Products: ProductModel{DB: db},
    }
}
```

#### **docker-compose.yaml** - Change service name & port
```yaml
version: '3.8'

services:
  product-service:  # Changed from user-service
    build: .
    ports:
      - "5000:5000"  # Different port
    environment:
      FASHIONPRODUCTS_DB_DSN: postgres://fashionuser:fashionpass@postgres:5432/fashionproducts?sslmode=disable
    depends_on:
      postgres:
        condition: service_healthy

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: fashionproducts  # Different DB name
    # ... rest same
```

---

## ❌ DON'T COPY (User-Service Specific)

### 1. **Entire Files to Skip**
- `cmd/api/tokens-handler.go` - no token generation
- `cmd/api/users-handler.go` - no user management
- `internal/data/users.go` - no user data
- `internal/data/tokens.go` - no token data
- `internal/data/permissions.go` - depends on your design*
- `internal/mailer/` - entire folder (no emails)

*Note: You might want permissions, but call user-service to check them

### 2. **Removed Dependencies**
```go
// DON'T import in product-service:
"github.com/go-mail/mail/v2"
"golang.org/x/crypto" // unless you need bcrypt for something else
```

### 3. **Removed Code Sections**
- SMTP configuration in main.go
- Mailer initialization
- Token generation/validation logic (delegate to user-service)
- Password hashing logic
- Email sending in background goroutines

---

## 🆕 NEW CODE (Product-Specific)

### 1. **internal/data/products.go** - Your core logic
```go
package data

import (
    "context"
    "database/sql"
    "errors"
    "time"
)

type Product struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    Stock       int       `json:"stock"`
    ImageURL    string    `json:"image_url"`
    CategoryID  int64     `json:"category_id"`
    SellerID    int64     `json:"seller_id"` // Reference to user from user-service
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Version     int       `json:"version"`
}

type ProductModel struct {
    DB *sql.DB
}

func (m ProductModel) Insert(product *Product) error {
    query := `
        INSERT INTO products (name, description, price, stock, image_url, category_id, seller_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at, version`
    
    args := []interface{}{
        product.Name,
        product.Description,
        product.Price,
        product.Stock,
        product.ImageURL,
        product.CategoryID,
        product.SellerID,
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    return m.DB.QueryRowContext(ctx, query, args...).Scan(
        &product.ID,
        &product.CreatedAt,
        &product.UpdatedAt,
        &product.Version,
    )
}

func (m ProductModel) Get(id int64) (*Product, error) {
    if id < 1 {
        return nil, ErrRecordNotFound
    }
    
    query := `
        SELECT id, name, description, price, stock, image_url, category_id, seller_id, created_at, updated_at, version
        FROM products
        WHERE id = $1`
    
    var product Product
    
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    err := m.DB.QueryRowContext(ctx, query, id).Scan(
        &product.ID,
        &product.Name,
        &product.Description,
        &product.Price,
        &product.Stock,
        &product.ImageURL,
        &product.CategoryID,
        &product.SellerID,
        &product.CreatedAt,
        &product.UpdatedAt,
        &product.Version,
    )
    
    if err != nil {
        switch {
        case errors.Is(err, sql.ErrNoRows):
            return nil, ErrRecordNotFound
        default:
            return nil, err
        }
    }
    
    return &product, nil
}

// Add: Update, Delete, GetAll with filters, etc.
```

### 2. **cmd/api/products-handler.go** - Your handlers
```go
package main

import (
    "errors"
    "net/http"
    
    "github.com/PaulBabatuyi/FashionMarket-Backend/product-service/internal/data"
)

func (app *application) listProductsHandler(w http.ResponseWriter, r *http.Request) {
    // Parse query params: ?page=1&page_size=20&sort=-created_at&name=shirt
    // Call products.GetAll()
    // Return JSON
}

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Name        string  `json:"name"`
        Description string  `json:"description"`
        Price       float64 `json:"price"`
        Stock       int     `json:"stock"`
        ImageURL    string  `json:"image_url"`
        CategoryID  int64   `json:"category_id"`
    }
    
    err := app.readJSON(w, r, &input)
    if err != nil {
        app.badRequestResponse(w, r, err)
        return
    }
    
    // Get authenticated user from context
    user := app.contextGetUser(r)
    
    product := &data.Product{
        Name:        input.Name,
        Description: input.Description,
        Price:       input.Price,
        Stock:       input.Stock,
        ImageURL:    input.ImageURL,
        CategoryID:  input.CategoryID,
        SellerID:    user.ID, // Current user is the seller
    }
    
    // Validate
    // Insert
    // Return JSON
}

func (app *application) getProductHandler(w http.ResponseWriter, r *http.Request) {
    id, err := app.readIDParam(r)
    // Get product
    // Return JSON
}

func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
    id, err := app.readIDParam(r)
    // Get product
    // Check if user is owner (product.SellerID == user.ID)
    // Update
    // Return JSON
}

func (app *application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
    id, err := app.readIDParam(r)
    // Get product
    // Check if user is owner
    // Delete
    // Return JSON
}
```

### 3. **migrations/** - Product schema
```sql
-- 000001_create_products_table.up.sql
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    image_url TEXT NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    
    -- Array of categories: ['men', 'women', 'unisex']
    category TEXT[] NOT NULL DEFAULT '{}',
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1,

    -- Full-text search
    tsv tsvector GENERATED ALWAYS AS (
        setweight(to_tsvector('english', name), 'A') ||
        setweight(to_tsvector('english', coalesce(description, '')), 'B')
    ) STORED
);

-- Indexes
CREATE INDEX idx_products_user_id ON products(user_id);
CREATE INDEX idx_products_category ON products USING GIN(category);  -- GIN for array
CREATE INDEX idx_products_created_at ON products(created_at);
CREATE INDEX idx_products_tsv ON products USING GIN(tsv);

-- 000001_create_products_table.down.sql
DROP TABLE IF EXISTS products;
```

---

## 📋 Quick Checklist for New Services

```
□ Copy jsonlog/, validator/ folders as-is
□ Copy errors.go, healthcheck.go, helpers.go, server.go as-is
□ Copy main.go → remove SMTP, change port, change DB DSN env var
□ Copy middleware.go → adapt authenticate() to call user-service
□ Copy routes.go → replace with service-specific routes
□ Copy Dockerfile → change binary name
□ Copy docker-compose.yaml → change service name, port, DB name
□ Copy go.mod → change module path, remove mail/crypto
□ Copy .gitignore, Makefile, .env.example → adapt names
□ Create new models in internal/data/
□ Create new handlers in cmd/api/
□ Create new migrations/
□ Write new README.md
```

---

## 🔑 Key Architecture Decision

**DON'T copy user validation logic to product-service!**

Instead, product-service should:
1. Extract token from Authorization header
2. Call user-service API: `GET http://user-service:4000/v1/auth/validate`
3. User-service validates token and returns user info
4. Product-service uses that info for authorization

This keeps services truly independent and avoids duplicating auth logic.