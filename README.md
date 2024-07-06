.env (placed in project dir) format:
```
───────┬────────────────────────────────────────────
   1   │ # Config
   2   │ SERVER_ADDRESS=000.000.000.000:0000
   3   │ 
   4   │ # Credentials for USER1
   5   │ USER1_API_KEY=...
   6   │ USER1_WORKSPACE_ID=...
   7   │ USER1_USER_NAME=...
   8   │ USER1_PAY_PER_HOUR=...
   9   │ 
  10   │ # Credentials for USER2
  11   │ USER2_API_KEY=...
  12   │ USER2_WORKSPACE_ID=...
  13   │ USER2_USER_NAME=...
  14   │ USER2_PAY_PER_HOUR=...
───────┴────────────────────────────────────────────
```

Report output:
```
| ID | User  | Duration | Sum   | Client | Task | Vikunja link |
|----|------------------|----------------|------|--------------|
| 1  | USER1 | 00:00:00 | 00.00 | ...    | ...  | ...          |
| 2  | USER1 | 00:00:00 | 00.00 | ...    | ...  | ...          |
| 3  | USER2 | 00:00:00 | 00.00 | ...    | ...  | ...          |
```
