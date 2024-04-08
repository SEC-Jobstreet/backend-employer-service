# backend-employer-service
This repo is Employer Management Service of jobstreet application backend.

## Build app
1. chỉnh sửa schema trong db/migration để phù hợp với repo hoặc tạo 1 migration mới bằng cách ```make new_migrate {tên migration}``` rồi thêm SQL QUERY vào file up để tạo DB schema. file down để migration down trong db/migration/...

2. Chạy docker DB lên bằng cách ```make run_postgres``` nếu lần đầu chạy hoặc ```make start postgres```

3. Tạo DB name trong DB vừa chạy ```make createdb```

4. migration code ở db/migration vào DB vừa chạy ở docker ```make migrate```

5. Viết query trong thư mục db/query (xem application_service làm mẫu)

6. tự động generate golang code từ db/query sang go code ```make sqlc```. Code sinh ra ở db/sqlc. Sài thôi

7. Chỉnh sửa code employer handler

8. ```go run main.go```

9. Lỗi gì sai đâu k biết thì hỏi