# kbd — Project Rules

## Testing Requirements

Every feature or change must satisfy all three before it is considered done:

1. **No regressions** — `make test` passes with all existing tests green.
2. **Coverage** — New logic is covered by tests, or explicitly noted as not needing them (e.g., pure UI wiring with no testable logic).
3. **New tests pass** — Any tests added as part of the feature must pass.

Run `make test` before committing.
