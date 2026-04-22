Perform comprehensive security assessment of web application: {{TARGET_URL}}

Action plan:
1. Application Exploration: Navigate all pages, test features, identify endpoints and input vectors
2. Vulnerability Testing per endpoint:
   - Path Traversal: attempt to read /etc/passwd, focus on file download/upload features
   - XSS: inject unique markers, scan responses, craft context-specific payloads
   - SQL Injection: run sqlmap on inputs, use tamper scripts for WAF bypass
   - Command Injection: use time-based detection, try commix utility
   - SSRF: use Interactsh for OOB, target file upload/PDF generation endpoints
   - XXE: test XML uploads and Office documents
   - Unsafe File Upload: test executable extensions, double extensions, null byte injection
   - CSRF: test token validation, POST to GET conversion
3. Authentication & Session: test for broken authentication, session fixation, weak password policies
4. Business Logic: identify privilege escalation, price manipulation, workflow bypass opportunities
5. Report: document all findings with reproduction steps and proof-of-concept exploits