# Project Status Board

- [x] Refactor: Combine TransactIncrementBy and TransactDecrementBy into TransactAddBy with float64 Value and named struct
- [x] Update all usages in codebase (client, mock, meterer, tests)
- [x] Update/add tests for both increment and decrement
- [x] Run and pass all relevant tests
- [ ] User review and confirmation

# Executor's Feedback or Assistance Requests

- All code and tests have been updated to use the new TransactAddBy method and TransactAddOp struct with float64 Value.
- Both increment (positive value) and decrement (negative value) are covered in code and tests.
- All tests pass after cleaning up Docker resources and rerunning.
- Please review the changes and confirm if the task is complete or if further adjustments are needed. 