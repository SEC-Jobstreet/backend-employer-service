# backend-employer-service
This repo is Employer Management Service of jobstreet application backend.

## Deploy

1. ```docker build -t thanhquy1105/backend-jobstreet-employer-service-prod:latest .```
2. ```docker push thanhquy1105/backend-jobstreet-employer-service-prod```
3. ```docker pull thanhquy1105/backend-jobstreet-employer-service-prod:latest```
4. ```docker run --name backend-jobstreet-employer-service-prod --network jobstreet-network -p 4001:4001 -e DB_SOURCE="postgresql://admin:admin@postgres:5432/employer_service_jobstreet?sslmode=disable" -d thanhquy1105/backend-jobstreet-employer-service-prod:latest```