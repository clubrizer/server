# User Service

## Auth flow

POST `/login` -> return access & refresh token
access token: memory
refresh token: http-only cookie

request: sends refresh token automatically (because cookie) & access token in a header -> backend
verifies access token

refresh page: refresh token still there, access token gone

refresh token cookie does nothing by default, when sending refreshtoken to /refresh, it gets a new access token 
which an attacker can't read

everytime clubrizer is loaded or when the access toekn is expired, a new access token is queried