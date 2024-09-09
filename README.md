# สามารถปรับการตั้งค่าที่ไฟล์ 
```
./configs/configs.yaml

server:
    port: 8080
    host: "0.0.0.0"

logger:
    production: false

restaurant:
    name: "OneSiam Fine Dining"
    maxTables: 100 # จำนวนโต๊ะสูงสุดที่ init ได้
    seatsPerTable: 4 # จำนวนที่นั่งต่อโต๊ะ
    code:
        charset: "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" # ตัวอักษร Gen โค๊ดคูปอง
        length: 6 # ความยาวของคูปอง **แนะนำห้ามตั้งต่ำเกินไปจะเกิดปัญหาคูปองซ้ำ

```

# Run Service
```
docker compose up
# OR
docker-compose up
```

# Call api
```
POST : http://localhost:3001/api/v1/initialize
BODY : { "tables": 100 }

POST : http://localhost:3001/api/v1/reserve
BODY : { "customers": 100 }

POST : http://localhost:3001/api/v1/cancel
BODY : { "bookingID": "30OTOI" }
```

# หากต้องการ docker build img
```
docker buildx build --platform linux/amd64,linux/arm64 -t {SERVER}/{REPO}:{VERSION} -t {SERVER}/{REPO}:latest --push .

```