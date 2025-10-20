API test helper files

Files:
- seed.sql: SQL to insert sample shop types, users, shops, vouchers, seckill record and one order.
- shop-type-list.http: GET /api/shop-type/list
- shop-list.http: shop list and shop detail tests
- voucher.http: voucher list, seckill add & get
- user-auth.http: request code, register, login and seckill purchase example

How to use:
1. Ensure your app is running (go run main.go) and Redis/MySQL are reachable.
2. Seed the database: run the SQL in seed.sql against the app database (after migrations).
3. Use VS Code REST Client extension or curl to run .http files. For the seckill purchase you need a valid JWT token obtained from /api/user/login (the server prints verification codes in logs in dev mode).

Notes:
- The seed.sql uses table names created by GORM AutoMigrate in main.go.
- If your DB has different prefixes or additional NOT NULL constraints, adjust seed.sql accordingly.
