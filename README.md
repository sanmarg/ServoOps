# ServoOps
![image](https://github.com/sanmarg/ServoOps/assets/50082154/ffb91c41-3738-4e61-a72f-0a215e2c8d79)

ServoOps is a Web API project written in Go using the Gin framework. The API provides endpoints to manage and monitor system services.
## Features
- List all running services with their PID, name, user, memory usage, CPU usage, and start time.
- Kill a service by its PID.
- Start a system service.
- Stop a system service.
- Health check endpoint to ensure the API is running.
- Basic authentication middleware.
- Logging middleware to log request details.

## Installation
1. Clone the repository: `git clone https://github.com/yourusername/yourrepository.git`
2. Navigate to the project directory: `cd yourrepository`
3. Run the application: `go run main.go`

## Endpoints
- `GET /services`: List all running services.
- `DELETE /services/:id`: Kill a service by its PID.
- `POST /services/start/:name`: Start a system service by its name.
- `POST /services/stop/:name`: Stop a system service by its name.
- `GET /health`: Health check endpoint.

## Authentication
- The API uses basic authentication with the secret token `c2FubWFyZwo=`.

## Logging
- Request details are logged using the logging middleware.

## Usage
1. Start the API.
2. Use the provided endpoints to manage and monitor system services.

## Example
Here is an example of how to use the API endpoints:

### List all running services
```bash
curl -H "Authorization: c2FubWFyZwo=" http://localhost:8082/
```
![image](https://github.com/sanmarg/ServoOps/assets/50082154/524f5480-e0dd-40cc-ad12-1aa10f75f9af)

### Kill a service by PID
```bash
curl -X DELETE -H "Authorization: c2FubWFyZwo=" http://localhost:8082/services/1234
```

### Start a system service
```bash
curl -X POST -H "Authorization: c2FubWFyZwo=" http://localhost:8082/services/start/servicename
```

### Stop a system service
```bash
curl -X POST -H "Authorization: c2FubWFyZwo=" http://localhost:8082/services/stop/servicename
```
## Contributing
- Fork the repository.
- Create a new branch: `git checkout -b feature`
- Make your changes and commit them: `git commit -m 'Add new feature'`
- Push to the branch: `git push origin feature`
- Submit a pull request.

## License
- This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
