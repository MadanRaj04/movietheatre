Converted Go backend (Gin + GORM MySQL + JWT)

How to run:
1. Install Go 1.20+
2. Set environment variables:
   - MYSQL_DSN (example: root:pass@tcp(127.0.0.1:3306)/cinema?parseTime=true)
   - SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASS (for email)
   - PORT (optional, default 3000)
3. go mod download
4. go run main.go

Notes:
- This project is a complete port of your Node app's routes and models.
- It uses GORM for DB access and auto-migrates tables.
- Authentication uses JWT; change jwtSecret in middleware/auth.go to a secure secret.
- QR codes are saved to uploads/qrcodes and served at /uploads/qrcodes/<file>.
