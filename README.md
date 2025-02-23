# Loan Service

A Buffalo-based loan management system with Docker support.

## Prerequisites

- Docker
- Docker Compose

## Quick Start

1. Clone the repository:
```bash
git clone <repository-url>
cd loan_service
```

2. Start the application using Docker:
```bash
docker-compose up -d
```

The application will be available at http://localhost:3000

## Environment Variables

The following environment variables can be configured in docker-compose.yml:

- `GO_ENV`: Application environment (development/production)
- `DB_HOST`: Database host
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `SESSION_SECRET`: Secret key for session encryption
- `PORT`: Application port (default: 3000)

## API Documentation

API documentation is available at http://localhost:3000/docs

## Database Migrations

Migrations are automatically run when the container starts. To manually run migrations:

```bash
docker-compose exec app ./app migrate
```

## Development

To rebuild the application after changes:
```bash
docker-compose down
docker-compose up -d --build
```

## Project Structure

```
loan_service/
├── actions/         # HTTP handlers
├── models/         # Database models
├── migrations/     # Database migrations
├── templates/      # HTML templates
├── public/         # Static files
└── docker/         # Docker configuration
```

## License

[License Type]

## Problem Statement

The system manages loans through multiple states with specific rules:

### Loan States
1. **Proposed** (Initial State)
   - Created when loan application is submitted
   - Requires: borrower ID and loan details

2. **Approved**
   - Requires validation proof:
     - Field validator's visit photo
     - Field validator's employee ID
     - Approval date
   - Cannot revert to proposed state
   - Makes loan available for investment

3. **Invested**
   - Achieved when investment equals loan principal
   - Can have multiple investors
   - Total investment cannot exceed principal
   - Triggers agreement letter generation and email to investors

4. **Disbursed**
   - Final state when loan is given to borrower
   - Requires:
     - Signed loan agreement (PDF/JPEG)
     - Field officer's employee ID
     - Disbursement date

### Required Services (APIs)

#### Loan Management
1. `POST /api/v1/loans`
   - Create new loan proposal
   - Initial state: PROPOSED

2. `GET /api/v1/loans`
   - List all loans
   - Supports filtering by status

3. `GET /api/v1/loans/{id}`
   - Get loan details

4. `PUT /api/v1/loans/{id}/approve`
   - Approve loan
   - Upload validation proof
   - Update field validator details

5. `PUT /api/v1/loans/{id}/disburse`
   - Disburse loan
   - Upload signed agreement
   - Update field officer details

#### Investment Management
1. `POST /api/v1/loans/{id}/investments`
   - Add investment to loan
   - Validate against principal amount

2. `GET /api/v1/loans/{id}/investments`
   - List all investments for a loan

3. `GET /api/v1/investments/{id}`
   - Get investment details

#### Document Management
1. `POST /api/v1/loans/{id}/documents`
   - Upload loan related documents
   - Supports validation proof and agreements

2. `GET /api/v1/loans/{id}/documents`
   - Get loan documents

#### Notification Service
1. `POST /api/v1/notifications/email`
   - Send agreement letters to investors

### Core Entities
- Loan
- Investment
- Document
- Notification

### Tech Stack
- Backend: Go with Buffalo Framework
- Database: MySQL
- ORM: POP (Buffalo's ORM)
- Frontend: Buffalo Templates (Server-side rendering)

### Project Structure
Following MVC (Model-View-Controller) pattern

## Database Setup

It looks like you chose to set up your application using a database! Fantastic!

The first thing you need to do is open up the "database.yml" file and edit it to use the correct usernames, passwords, hosts, etc... that are appropriate for your environment.

You will also need to make sure that **you** start/install the database of your choice. Buffalo **won't** install and start it for you.

### Create Your Databases

Ok, so you've edited the "database.yml" file and started your database, now Buffalo can create the databases in that file for you:

```console
buffalo pop create -a
```

## Starting the Application

Buffalo ships with a command that will watch your application and automatically rebuild the Go binary and any assets for you. To do that run the "buffalo dev" command:

```console
buffalo dev
```

If you point your browser to [http://127.0.0.1:3000](http://127.0.0.1:3000) you should see a "Welcome to Buffalo!" page.

**Congratulations!** You now have your Buffalo application up and running.

## What Next?

We recommend you heading over to [http://gobuffalo.io](http://gobuffalo.io) and reviewing all of the great documentation there.

Good luck!

[Powered by Buffalo](http://gobuffalo.io)
