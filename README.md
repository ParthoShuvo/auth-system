# Auth-System

An User **Auth**entication **System** that will allow users to Login via Email/Phone and Password combination and receive a short lived Access Token that will allow them to access some authenticated routes in other services (these services are out of scope of this project but think about the use cases).

## Outline

- [Auth-System](#auth-system)
  - [Outline](#outline)
    - [Features to be implemented](#features-to-be-implemented)

### Features to be implemented

1. **Registration** - optional, can go with already populated Users in DB. If you decide to go for it, mock any verification process ex. Email/Phone Verification
2. **Login**
3. **Access & Refresh Tokens** - upon successful login User will receive an Access Token(short lived) and a Refresh Token(relatively long lived, can be used to avoid forcing the user to login each time an Access Token expires).
4. **JWT tokens** are preferable.
5. **Verify Token** - endpoint to verify an Access Token. Verified Access token will return the User's profile, role, permission etc.
6. **New Access Token** -  endpoint to acquire a new Access Token using the Refresh Token generated upon Login.
