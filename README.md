# Cerebra – Backend (Go)

Máy chủ trung gian của **Cerebra**, đồ án môn Học máy dự đoán nguy cơ đột quỵ. Backend nhận dữ liệu bệnh nhân từ frontend, kiểm tra tính hợp lệ, chuyển tiếp sang dịch vụ dự đoán Python (FastAPI) rồi trả kết quả về frontend.

```
React ──POST /api/v1/predict──► Go (Gin) ──POST PYURL──► Python (FastAPI)
React ◄────────── JSON ───────── Go ◄────────── JSON ───── Python
```

Python không cần biết địa chỉ của Go: Go gọi sang và nhận kết quả ngay trên cùng kết nối.

## Công nghệ

- Go, framework Gin
- godotenv (đọc file `.env`)

## Cấu trúc thư mục

```
DoAnBackend/
├── config/
│   └── config.go        # Đọc biến môi trường
├── dataset/
│   └── healthcare-dataset-stroke-data.csv
├── handler/
│   └── handler.go       # Kiểm tra dữ liệu, gọi Python, trả kết quả
├── router/
│   └── router.go        # Khai báo route và CORS
├── .env                 # Cấu hình trên máy (không đưa lên Git)
├── .env.example         # Mẫu cấu hình
├── .gitignore
├── go.mod
├── go.sum
├── main.go              # Điểm vào: nạp config, chạy server
└── README.md
```

Ghi chú: thư mục `dataset/` chứa dữ liệu gốc, backend Go không đọc đến nó (dùng cho phần huấn luyện bên repo Python).

## Cấu hình

Sao chép `.env.example` thành `.env` và điền giá trị:

| Biến | Ý nghĩa | Ví dụ |
|---|---|---|
| `PORT` | Cổng chạy của Go | `8080` |
| `PYURL` | Địa chỉ **đầy đủ** của dịch vụ Python, gồm cả `/predict` | `http://localhost:8000/predict` |
| `ALLOWED_ORIGINS` | Địa chỉ frontend được phép gọi (CORS) | `http://localhost:5173` |

Lưu ý: thiếu `http://` ở đầu `PYURL` sẽ gặp lỗi `unsupported protocol scheme`; thiếu `/predict` ở cuối sẽ nhận 404 từ Python.

## Chạy

```bash
go mod download
go run .
```

Mặc định server chạy tại `http://localhost:8080`. Cần bật dịch vụ Python trước.

## API

### `POST /api/v1/predict`

Yêu cầu:

```json
{
  "gender": "Male",
  "age": 86,
  "hypertension": 1,
  "heart_disease": 1,
  "ever_married": "No",
  "work_type": "Private",
  "Residence_type": "Urban",
  "avg_glucose_level": 150,
  "bmi": 50,
  "smoking_status": "never smoked",
}
```

| Trường | Giá trị hợp lệ |
|---|---|
| `gender` | `Male`, `Female`, `Other` |
| `age` | lớn hơn 0 và nhỏ hơn 120 |
| `hypertension`, `heart_disease` | `0` hoặc `1` |
| `ever_married` | `Yes`, `No` |
| `work_type` | `Private`, `Self-employed`, `Govt_job`, `children`, `Never_worked` |
| `Residence_type` | `Urban`, `Rural` |
| `avg_glucose_level`, `bmi` | lớn hơn 0 |
| `smoking_status` | `never smoked`, `formerly smoked`, `smokes`, `Unknown` |

Phản hồi thành công (`200`):

```json
{
  "probability": 50.33,
  "high_risk": true,
  "risk_level": "RẤT CAO",
  "advice": "Nguy cơ rất cao. Cần đến bệnh viện kiểm tra càng sớm càng tốt.",
  "factors": ["Tuổi cao (86 tuổi)", "Bị cao huyết áp", "Có bệnh tim"],
  "warning": "Tuổi vượt quá dữ liệu huấn luyện (tối đa 82), kết quả chỉ mang tính ước lượng."
}
```

`probability` là tỉ lệ ước lượng theo phần trăm (đã hiệu chỉnh), không phải xác suất y khoa thật. `warning` rỗng khi tuổi nằm trong vùng dữ liệu đã học.

### Mã lỗi

| Mã | Khi nào |
|---|---|
| `400` | Dữ liệu gửi lên sai định dạng hoặc giá trị không hợp lệ |
| `502` | Không kết nối được tới Python, hoặc Python trả lỗi |
| `500` | Lỗi xử lý nội bộ (đóng gói hoặc đọc kết quả) |

## Lỗi thường gặp

- **`Post "": unsupported protocol scheme`:** `PYURL` đang rỗng. Kiểm tra tên biến trong `.env` là `PYURL`, file `.env` nằm đúng thư mục chạy `go run .`, và đã khởi động lại Go sau khi sửa.
- **`502 Không kết nối được tới dịch vụ dự đoán`:** dịch vụ Python chưa chạy hoặc `PYURL` sai cổng.
- **Trình duyệt báo CORS:** `ALLOWED_ORIGINS` chưa khớp địa chỉ frontend.

## Các repo liên quan

| Phần | Repo | Cổng mặc định |
|---|---|---|
| Frontend (React) | `<link-repo-frontend>` | 5173 |
| Backend Go (Gin) | `<link-repo-backend>` | 8080 |
| Huấn luyện và API Python | `<link-repo-train>` | 8000 |

## Tác giả

- `Nhóm 12 - Trưởng nhóm: Nguyễn Minh Tiến - 2001240494`
- Môn: `Học máy` – Giảng viên: `TS.Phùng Thế Bảo`