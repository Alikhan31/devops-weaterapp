# Weather Application

A modern weather application built with Go, featuring a clean frontend interface, Redis caching, and Docker containerization.

## Features

- Real-time weather data from OpenWeatherMap API
- Redis caching for improved performance
- Docker containerization for easy deployment
- CI/CD pipeline with GitHub Actions
- Responsive frontend design

## Prerequisites

- Docker and Docker Compose
- OpenWeatherMap API key

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/Alikhan31/devops-weaterapp.git
cd devops-weaterapp
```

2. Create a `.env` file in the root directory with your OpenWeatherMap API key:
```
OPENWEATHER_API_KEY=your_api_key_here
```

3. Start the application:
```bash
docker-compose up -d
```

4. Access the application:
- Frontend: http://localhost:80
- Backend API: http://localhost:8080/weather

## Project Structure

```
.
├── backend/           # Go backend service
│   ├── main.go       # Main application code
│   └── Dockerfile    # Backend container configuration
├── frontend/         # Frontend service
│   ├── index.html    # Main HTML file
│   ├── styles.css    # CSS styles
│   ├── script.js     # Frontend JavaScript
│   └── Dockerfile    # Frontend container configuration
├── docker-compose.yml # Docker services configuration
└── .github/          # GitHub Actions workflows
    └── workflows/    # CI/CD pipeline configuration
```

## API Endpoints

### GET /weather
Returns weather information for a specific location.

Query Parameters:
- `lat`: Latitude
- `lon`: Longitude

Example:
```
GET http://localhost:8080/weather?lat=51.5074&lon=-0.1278
```

Response:
```json
{
    "temperature": 17.08,
    "description": "clear sky",
    "city": "London",
    "humidity": 38,
    "windSpeed": 6.69
}
```

## Development

### Running Tests
```bash
cd backend
go test ./...
```

### Building Locally
```bash
# Backend
cd backend
go build

# Frontend
cd frontend
# Serve with any static file server
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 