# My Blog Posts Platform (Back-End)
In my Golang Back-End, I executed advanced SQL queries, ensured data normalization, db migrations, and established proper indexes. Developed a suite of middlewares for, secureHeaders and client-based rate limiting, to manage request traffic. I handled race conditions using mutual exclusion locks in the app and a versioning system in DB. Established a graceful shutdown process for the server, systematically formatted JSON responses, ensured proper error handling, and implemented automated API versioning using git.
Tests are written using Go and API is tested end-to-end using Postman. Application metrics are displayed and load-tested properly. The codebase is maintained on GitHub with dependency injection and third-party dependency security ensured. The CI/CD pipeline is implemented with GitHub workflow and Google Cloud Run. Mailtrap is used as the SMTP server for email dispatching.
Adopted multi-stage Docker file, and leveraged AWS ECR and RDS services. Implemented a stateful authentication process using tokens. Viper and cmd-flags are employed for app configuration. Load balancing is handled in Google Cloud Run, with support for HTTPS and TLS certificates.

## Technologies Used
Golang, Google Cloud, AWS, PostgreSQL, Docker, Postman

## Available Scripts
All the available commands are mentioned in Makefile in the root directory.

