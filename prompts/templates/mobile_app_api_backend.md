Perform security testing of mobile application backend API: {{API_URL}}

Action plan:
1. Traffic Interception: analyze mobile app traffic, extract API endpoints and authentication
2. Authentication Mechanisms: test OAuth flows, JWT implementation, refresh token handling, certificate pinning bypass
3. API Endpoint Testing: test all discovered endpoints for BOLA/IDOR, broken function-level authorization
4. Data Validation: test for injection attacks in API parameters, test file upload endpoints
5. Business Logic: test premium feature bypass, subscription validation, in-app purchase verification
6. Session Management: test token expiration, concurrent session handling, session fixation
7. Sensitive Data: check for PII exposure, excessive data in responses, hardcoded secrets
8. Rate Limiting: test brute force protection on login, API rate limits, account lockout
9. Deep Linking: test for deep link hijacking, intent redirection (Android), URL scheme abuse (iOS)
10. Report: mobile-specific vulnerabilities with mitigation recommendations