# EzKitchen

EzKitchen is a web application for managing kitchen remodeling estimates from initial survey/walkthrough to job completion.  
It is built around the day-to-day workflow of a small remodeling business, where estimates move through a defined set of states and multiple users interact with the same job over time.

The application handles authentication, estimate creation, line items, status transitions, and persistence using a PostgreSQL database. All pages are rendered server-side using Go templates, with standard HTTP form submissions.

This project was built as a portfolio piece with an emphasis on backend structure, data modeling, and security.

---

## What the Application Does

EzKitchen allows authenticated users (Administrators or Surveyors) to create, view, and manage customer estimates.

A typical workflow looks like this:

1. A surveyor logs in and creates a new estimate for a customer.
2. The estimate is populated with line items such as cabinetry, appliances, countertops, or labor.
3. Each estimate progresses through a series of statuses (for example: draft, awaiting agreement, in progress, completed).
4. Users can return to existing estimates to update details or advance the job to the next stage.
5. The system enforces validation rules and access controls to prevent invalid or unauthorized changes.

The goal is to provide a clear, structured way to track work from the initial on-site survey through completion.

## Features

- User authentication with server-side sessions
- Role-based access control
- Estimate creation and editing
- Line item management
- Status-based workflow tracking
- Server-side HTML rendering
- Form validation and error feedback
- Flash messages and redirects
- PostgreSQL-backed persistence

---

## Tech Stack

- **Go** - HTTP server, routing, handlers
- **PostgreSQL** - relational data storage
- **html/template** - server-side rendering
- **JavaScript** - client-side enhancements and interactions

---

## Project Structure

- cmd/web/
- main.go # application startup and configuration
- routes.go # route definitions
- middleware.go # HTTP middleware
- context.go # request-scoped context helpers
- template.go # template rendering helpers
- handlers\_\*.go # HTTP handlers grouped by domain

- internal/
- models/ # database models and queries
- storage/ # external storage integrations
- validator/ # input validation helpers
- mailer/ # email-related logic

- migrations/ # database schema migrations

- ui/html/ # server-rendered templates
- ui/static/ # CSS and static assets

Handlers are grouped by domain (users, estimates, products, etc.) to keep files readable as the application grows.  
Database access is separated in the internal/models package and dependencies are initialized once at startup and shared through the application struct.
