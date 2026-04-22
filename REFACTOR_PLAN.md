# Homeserver Refactor Plan: "The Great Shearing"

## Phase 0: The "Structural Debt" Audit (DONE)
**Goal:** Dismantle over-engineered abstractions and fix bad naming.

- [x] **Package Consolidation:**
    - [x] **Abolish `db/query_builder/`**: Delete all sub-packages. Raw SQL is the new standard.
    - [x] **Flatten `constants/`**: Merge `uconsts`, `fconsts`, `tconsts`, etc. into their respective repository files or a single `constants/constants.go`.
    - [x] **Flatten `pkg/utils/`**: Merge `errresp`, `generics`, and `formatters` into a single `pkg/utils` package.
- [x] **Idiomatic Naming & Structure:**
    - [x] Rename `InUserRepo` -> `UserRepo` (Interface).
    - [x] Rename implementation `UserRepo` -> `pgUserRepo`.
    - [x] Move `db/repositories/[repo]_repo/schema/*` directly into `db/[repo]_repo.go`.
- [x] **Remove "Middleman" Packages:** Delete packages that only contain one or two small functions.

## Phase 1: Bare Database & Repository Overhaul (DONE)
**Goal:** Switch to a direct, low-ceremony database layer.

- [x] **The "God" Database Object:** 
    - [x] Rename the `Repositories` struct to `Database` and move it to the `db` package.
    - [x] Add a `Raw` field to the `Database` struct that exposes the `QueryExecutor` (the pool) directly.
- [x] **Ditch `db.Database[T]`**: Remove the generic struct and the reflection-based `ScanRows`.
- [x] **QueryExecutor Standard:** Repository implementations now depend directly on `interfaces.QueryExecutor`.
- [x] **Explicit Scanning:** Raw SQL with direct `.Scan()` calls.
- [x] **Consolidated Repositories:** User, File, Folder, Token are each a single file in `db/`.

## Phase 2: Testing & Storage Overhaul (PENDING)
**Goal:** Fix the "Enterprise C++" test ceremony and add file flexibility.

- [ ] **Abolish `RunScan` and `MockRow`**: Remove all reflection-based mocking helpers.
- [ ] **Mock the Repo, Not the Driver:** Web handlers will mock `UserRepo`, not `pgx`.
- [ ] **Storage Interface:** Create `interfaces.FileStorage` to abstract `os` calls.
- [ ] **Standard Logging:** Replace the empty `internal/logger` with standard `log/slog`.

## Phase 3: Modernization (PARTIAL)
- [x] **Update JWT:** Migrated to `golang-jwt/jwt/v5`.
- [x] **Middleware Fix:** Pre-parse templates in `AuthMiddleware`.
- [ ] **Abolish redundant interfaces:** Clean up `interfaces/database.go` if it still exists and is redundant.
