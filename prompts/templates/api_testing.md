Perform comprehensive API security assessment: {{API_BASE_URL}}

Action plan:
1. API Discovery: identify all endpoints, HTTP methods, parameters
2. Authentication Testing: test broken authentication, token manipulation, JWT vulnerabilities
3. Authorization Testing: test broken object-level authorization (BOLA/IDOR), function-level authorization bypass
4. Input Validation: test injection attacks (SQL, NoSQL, Command, XXE), mass assignment vulnerabilities
5. Rate Limiting: test for absence of rate limiting, brute force protection
6. Business Logic: test for excessive data exposure, lack of resource limiting, unsafe consumption of APIs
7. Security Misconfiguration: check CORS policy, security headers, verbose error messages
8. GraphQL Specific (if applicable): test introspection, query depth limits, batching attacks
9. Report: document API vulnerabilities with curl/Postman proof-of-concepts