# TinyR URL Shortener – Backend

TinyR là hệ thống rút gọn URL được phát triển bằng Golang, hướng đến hiệu năng cao, bảo mật và khả năng mở rộng.  
Dự án áp dụng Service-Oriented Architecture kết hợp kiến trúc nhiều lớp (Clean Architecture) để dễ bảo trì và mở rộng.

## Công nghệ chính

| Lĩnh vực        | Công nghệ        |
|----------------|------------------|
| Ngôn ngữ        | Go (Golang)      |
| Web Framework   | Gin              |
| Database        | PostgreSQL       |
| ORM             | GORM             |
| Caching         | Redis            |
| Auth            | JWT              |
| Email Service   | Gomail + SMTP    |

## Kiến trúc thư mục
```
├── internal/
│ ├── config/ # Đọc .env
│ ├── handler/ # HTTP Controller
│ ├── middleware/ # Auth, Verification
│ ├── model/ # GORM Models
│ ├── repository/ # PostgreSQL + Redis
│ └── service/ # Business Logic
├── router/ # Routes & Middleware
└── main.go # App entrypoint
```
## Tính năng nổi bật

- Caching theo kiến trúc Cache-Aside với Redis để giảm truy vấn DB.
- Sliding TTL giữ các URL phổ biến trong cache lâu hơn.
- Xác thực hai lớp: JWT + trạng thái tài khoản `is_verified`.
- 2 loại token: đăng ký tài khoản và đặt lại mật khẩu.
- Middleware phân tầng truy cập: chỉ đăng nhập / chỉ verified / public.
- Dependency Injection cho Email Service, dễ thay đổi SMTP provider.
- Sử dụng biến môi trường (.env) để bảo mật thông tin hệ thống.
## Hướng dẫn chạy

### 1. Clone và cài dependencies
```bash
git clone https://github.com/sung2708/shorten_url
shorten_url
go mod tiny
```
### 2. Tạo file `.env
```
PORT="
JWT_SECRET="
REDIS_HOST=""
REDIS_USER="
REDIS_PASSWORD="
POSTGRES_DB=""
SENDER_EMAIL=""
SMTP_HOST=""
SMTP_PORT=""
SMTP_USER=""
SMTP_PASS=""
```
### 3. Khởi chạy DB (khuyến nghị Docker)
[Link](https://hub.docker.com/_/postgres/)

### 4. Chạy ứng dụng

```bash
go run ./cmd/api/main.go
# or
air
```
Ứng dụng mặc định chạy tại: `http://localhost:8080`

Project by **sung2708**
