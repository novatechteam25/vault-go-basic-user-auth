┌─────────────────────────────────────────────────────────────────────────┐
│                         CLIENT REQUEST                                   │
│                                                                           │
│  curl -H "Authorization: Bearer TOKEN" http://localhost:8080/route      │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │
                                 ▼
                    ┌────────────────────────┐
                    │  Gin Router Receives   │
                    │      Request           │
                    └────────────┬───────────┘
                                 │
                                 ▼
                    ┌────────────────────────────────┐
                    │  vaultAuthMiddleware()         │
                    │  Checks "Authorization" header │
                    └────────────┬───────────────────┘
                                 │
                    ┌────────────┴──────────────┐
                    │                           │
            Header Missing?              Header Found
                    │                           │
                    ▼                           ▼
            Return 401 Error         Extract Token (remove "Bearer ")
                    │                           │
                    │                           ▼
                    │            ┌──────────────────────────────┐
                    │            │  Send Token to Vault         │
                    │            │  client.SetToken(token)      │
                    │            │  LookupSelf()                │
                    │            └──────────┬───────────────────┘
                    │                       │`
                    │        ┌──────────────┴──────────────┐
                    │        │                             │
                    │    Token Valid?              Token Invalid/Expired
                    │        │                             │
                    │        ▼                             ▼
                    │    ✓ Continue              Return 401 Error
                    │        │                             │
                    │        ▼                             │
                    │   ┌──────────────────────┐           │
                    │   │  Route Handler       │           │
                    │   │  /my-protected-route │           │
                    │   │                      │           │
                    │   │  Response:           │           │
                    │   │  {                   │           │
                    │   │    "message":        │           │
                    │   │    "Access granted"  │           │
                    │   │  }                   │           │
                    │   └──────────┬───────────┘           │
                    │              │                        │
                    └──────────────┼────────────────────────┘
                                   │
                                   ▼
                        ┌──────────────────────┐
                        │   Client Response    │
                        │   200 OK or 401      │
                        └──────────────────────┘