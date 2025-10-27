# catiSip - User Feedback System

对标南康功能的，基于sip的系统实现，集成freeswitch

## User Feedback Feature

This project includes a user feedback system that allows users to submit feedback with ratings and messages.

### Features

- **Feedback Form**: Clean, responsive UI for submitting feedback
- **Star Rating**: 1-5 star rating system
- **Validation**: Client and server-side validation
- **In-memory Storage**: Feedback stored in memory (can be extended to database)
- **CORS Support**: Backend configured for cross-origin requests

### Architecture

- **Backend**: Go (Golang) with Gorilla Mux router
- **Frontend**: React 18 with modern hooks
- **Containerization**: Docker and Docker Compose for easy deployment

### Getting Started

#### Prerequisites

- Docker and Docker Compose (recommended)
- OR Node.js 18+ and Go 1.21+ for local development

#### Quick Start with Docker

1. Clone the repository:
```bash
git clone https://github.com/phughe11/catiSip.git
cd catiSip
```

2. Start the application:
```bash
docker-compose up --build
```

3. Access the application:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

#### Local Development

**Backend:**
```bash
cd backend
go mod download
go run main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm start
```

### API Endpoints

#### Submit Feedback
```bash
POST /api/feedback
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "message": "Great service!",
  "rating": 5
}
```

**Response:**
```json
{
  "id": "feedback-1",
  "name": "John Doe",
  "email": "john@example.com",
  "message": "Great service!",
  "rating": 5,
  "created_at": "2025-10-27T12:00:00Z"
}
```

#### Get All Feedback
```bash
GET /api/feedback
```

**Response:**
```json
[
  {
    "id": "feedback-1",
    "name": "John Doe",
    "email": "john@example.com",
    "message": "Great service!",
    "rating": 5,
    "created_at": "2025-10-27T12:00:00Z"
  }
]
```

### Testing

Test the backend API:
```bash
curl -X POST http://localhost:8080/api/feedback \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "message": "This is a test feedback",
    "rating": 5
  }'
```

Get all feedback:
```bash
curl http://localhost:8080/api/feedback
```

### Project Structure

```
.
├── backend/
│   ├── Dockerfile
│   ├── go.mod
│   └── main.go
├── frontend/
│   ├── Dockerfile
│   ├── nginx.conf
│   ├── package.json
│   ├── public/
│   │   └── index.html
│   └── src/
│       ├── App.css
│       ├── App.js
│       ├── index.css
│       └── index.js
├── docker-compose.yml
└── README.md
```

### Future Enhancements

- [ ] Add database persistence (PostgreSQL/MongoDB)
- [ ] Add user authentication
- [ ] Add feedback analytics dashboard
- [ ] Add email notifications for new feedback
- [ ] Add feedback moderation features
- [ ] Add pagination for feedback list
- [ ] Add search and filter capabilities

### License

MIT

