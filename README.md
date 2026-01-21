# AI Operations App

A full-stack web application for managing product orders, viewing mockups, and generating shipping labels.

This project demonstrates how an order moves through a simple fulfillment workflow — from creation to approval to shipping label generation.

---

## Features

- **View customer orders** - Browse and manage all orders in the system
- **View product mockups** - See generated product mockups with customer logos
- **Approve orders for fulfillment** - Move orders through the workflow
- **Enter shipping details** - Add customer shipping information
- **Generate printable shipping labels (PDF)** - Create shipping labels with barcodes
- **Barcode automatically added to each label** - Code128 barcodes for tracking

---

## Tech Stack

### Frontend
- React
- TypeScript
- Vite
- React Router

### Backend
- Go
- Gin Framework
- GORM (SQLite)
- gofpdf (PDF generation)
- Code128 barcode generator

---

## Project Structure

```
ai_operations_app/
├── backend/
│   ├── main.go                 # Application entry point
│   ├── go.mod                  # Go dependencies
│   ├── db/
│   │   └── db.go              # Database connection
│   ├── models/
│   │   └── order.go           # Data models
│   ├── handlers/
│   │   ├── order_handler.go   # Order API endpoints
│   │   └── upload_handler.go  # File upload handling
│   ├── services/
│   │   ├── mockup.go          # Mockup generation
│   │   ├── workflow.go        # Order workflow
│   │   └── label.go           # Shipping label generation
│   ├── assets/                # Static assets (templates)
│   ├── uploads/               # User uploaded files
│   ├── mockups/               # Generated mockups
│   └── labels/                # Generated shipping labels
├── frontend/
│   ├── src/
│   │   ├── App.tsx            # Main application component
│   │   ├── main.tsx           # Application entry point
│   │   └── pages/
│   │       ├── CreateOrder.tsx    # Order creation form
│   │       ├── Orders.tsx         # Orders list view
│   │       └── OrderDetail.tsx    # Order details and actions
│   ├── package.json           # Frontend dependencies
│   └── .env                   # Environment variables
└── README.md                  # This file
```

---

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- npm or yarn

### Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

3. Run the backend server:
   ```bash
   go run main.go
   ```

   The API will be available at `http://localhost:8080`

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Create a `.env` file:
   ```bash
   VITE_API_URL=http://localhost:8080
   ```

4. Start the development server:
   ```bash
   npm run dev
   ```

   The frontend will be available at `http://localhost:5173`

---

## API Endpoints

### Orders
- `POST /orders` - Create a new order
- `GET /orders` - List all orders
- `GET /orders/:id` - Get order details
- `POST /orders/:id/approve` - Approve order for fulfillment
- `POST /orders/:id/mockup` - Upload logo and generate mockup
- `POST /orders/:id/label` - Generate shipping label

### File Uploads
- `POST /upload/logo` - Upload logo file

### Static Files
- `/mockups/*` - Generated mockup images
- `/uploads/*` - Uploaded files
- `/labels/*` - Generated shipping labels
- `/assets/*` - Static assets

---

## Order Workflow

1. **CREATED** - Order is created with product details
2. **MOCKUP_GENERATED** - Logo is uploaded and mockup is generated
3. **APPROVED** - Order is approved for fulfillment
4. **READY_FOR_FULFILLMENT** - Shipping label is generated

---

## Usage Examples

### Creating an Order

```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "product": "Hoodie",
    "color": "Black",
    "size": "M",
    "logoUrl": "https://example.com/logo.png"
  }'
```

### Generating a Shipping Label

```bash
curl -X POST http://localhost:8080/orders/1/label \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "address": "123 Main St",
    "city": "Anytown",
    "state": "CA",
    "zip": "12345"
  }'
```

---

## Development

### Database

The application uses SQLite for development with automatic migrations. The database file (`minicronk.db`) is created automatically when the backend starts.

### File Storage

- **Uploads**: User-uploaded logos are stored in `backend/uploads/`
- **Mockups**: Generated product mockups are stored in `backend/mockups/`
- **Labels**: Generated shipping labels are stored in `backend/labels/`

### Environment Variables

#### Frontend (.env)
```
VITE_API_URL=http://localhost:8080
```

---

## Deployment

### Backend
The backend can be deployed as a standalone Go binary or using Docker:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

### Frontend
Build the frontend for production:

```bash
npm run build
```

The built files in `dist/` can be served by any static file server.

---

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

---

## License

This project is licensed under the MIT License.

