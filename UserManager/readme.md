# Workflow
## Register/Login
### Register:
    - User --> ServiceRouter --> UserManager -- /Register -- Database -- Returns Status --> Usermanager --> ServiceRouter --> User

### Login:
    - User --> ServiceRouter --> UserManager -- /Login -- Database -- Returns UUID & Role --> AuthManager -- Returns JWT -->  UserManager --> ServiceRouter --> User


