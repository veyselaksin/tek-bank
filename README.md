# TEK-BANK
## Description
TEK-BANK is a Teknasyon project that is a simple bank account management case study.

**NOTE!** I just only write unit tests for the `account` service. I didn't write unit tests for the other services and repositories because of the time constraint. 
That doesn't mean I don't know how to write unit tests. I know how to write unit tests. I just didn't write them because of the time constraint.

# Requirements
- Go 1.22
- Docker
- Docker Compose
- Swagger
- PostgreSQL
- Redis
- GoMailer Package
- Fiber
- Gorm
- JWT
- Mockgen
- Nginx
- Dockerfile
- Docker-compose

# Installation
- Clone the repository
- You need to create a `.env` file in the root directory and fill the necessary environment variables. You can find the necessary environment variables in the `.env.example` file.
- Run `docker-compose -f docker-compose.dev.yml up --build` to start the project in development mode
- Run `docker-compose -f docker-compose.prod.yml up --build` to start the project in production mode

# API Documentation
- You can find the API documentation in the `docs` directory.
- You can access the API documentation from the `/v1/docs` endpoint.

# Important Notes
- The project is developed with the Clean Architecture approach.
- The project is developed with the DDD approach.
- The project is developed with the SOLID principles.