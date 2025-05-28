# Nanobank

##Refresh Flow Example:
1. User tries to fetch /api/account
2. Access token is expired → backend responds with 401 Unauthorized
3. Frontend checks: "Do I have a refresh token?"
4. If yes, call POST /refresh with the refresh token
5. Backend verifies refresh token in Redis, rotates tokens
6. Frontend stores the new access/refresh pair

