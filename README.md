## Dependencies
- **sqlc**
- **goose**
- **air**

## Running the Project
1. Open **two separate terminals**.
2. Copy `.env.example` to `.env` 
3. In the first terminal:
   - Run `make docker-run` to start the Docker containers. Wait for it to complete pulling images and setting up the database.
4. In the second terminal:
   - Run `make db-migrate` to apply migrations for the public schema.
   - After successful migration, run `make watch` to start the application.
5. Import the Postman collection from the `docs` directory.
6. Test the API using Postman.

> **Disclaimer:**  
> This code is designed only as a **proof of concept** for multi-tenancy and is not intended for production use. To simplify the implementation, many best practices and security measures have been intentionally overlooked. For example:
> - Errors are not handled properly.
> - User inputs are not validated.  
> Use this codebase only as a reference or learning material.

---

## Generating Migration Files
### Public Schema
- Command:  
  ```bash
  goose -s -dir=./db/migrations/public create <migration_name> sql
  ```
  **Example:**  
  ```bash
  goose -s -dir=./db/migrations/public create init_public_schema sql
  ```

### Tenant Schema
- Command:  
  ```bash
  goose -s -dir=./db/migrations/tenant create <migration_name> sql
  ```
  **Example:**  
  ```bash
  goose -s -dir=./db/migrations/tenant create init_tenant_schema sql
  ```

---

## Generating SQLC Files from Queries
- Command:  
  ```bash
  sqlc generate
  ```

---

## API Documentation
- The Postman collection is available in the `docs` directory at the project root.
