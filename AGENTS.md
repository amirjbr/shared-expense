# Shared Expense development guidelines

## Architecture

This project is a Go modular monolith using Hexagonal Architecture and DDD.

Modules:
- account
- groups
- ledger

Dependency direction:
- handler -> service/application -> domain
- repository adapters implement core ports
- domain code must not depend on Gin, PostgreSQL, JWT or configuration

## HTTP boundaries

HTTP DTOs contain transport types.

UUID values coming from URLs or JSON must be parsed in handlers.

Inside services, commands, domain entities and repository ports, use uuid.UUID
instead of string IDs.

Handlers must not contain business logic.

## Cross-module communication

Modules must not query another module's tables directly.

Use narrow consumer-owned interfaces.

Examples:
- Groups depends on UserFinder implemented by Account.UserLookup.
- Ledger depends on GroupFinder implemented by Groups.GroupLookup.

Do not introduce direct cyclic module dependencies.

## Ledger rules

Money is stored as whole Toman using int64 / entity.Money.

Expense rules:
- total_amount > 0
- at least one payer
- payer total must equal expense total
- participants and payers are independent
- equal split remainder must be deterministic
- custom share total must equal expense total
- only group members can participate

Balance formula:

balance =
paid
- shares
+ settlements_sent
- settlements_received

The sum of all group balances must always equal zero.

Do not persist calculated balances.

## Database

Use PostgreSQL through database/sql and lib/pq.

Multi-table aggregate writes must use transactions.

Do not use SELECT *.

Always check rows.Err() after iterating rows.

Use explicit column lists in queries.

## MVP

Do not add non-MVP features unless explicitly requested.

Out of scope:
- guest users
- multiple currencies
- debt simplification
- notifications
- receipt OCR
- payment integrations
- settlement confirmation

## Verification

Before finishing code changes run:

go test ./...
go vet ./...