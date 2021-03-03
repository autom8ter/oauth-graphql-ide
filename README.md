# oauth-graphql-playground

An oauth protected [graphQL playground](https://github.com/graphql/graphql-playground)


## Features

- [x] Serves [GraphQL Playground](https://github.com/graphql/graphql-playground) user interface `/`
- [x] Login with [oauth authorization code grant](https://oauth.net/2/grant-types/authorization-code/)
- [x] Serve local auth proxy to remote graphQL server/endpoint 
    - automatically adds authorization header with oauth bearer token to outbound request
- [x] Fully Configurable via environmental variables
- [x] Pluggable session management
    - [x] Cookie-based sessions
    - [ ] Redis-based sessions


## Installation

#### Using Containers

- Docker: `docker pull colemanword:oauth-graphql-playground:v0.0.0`

- [sample k8s manifest](k8s.yaml) 

- [sample docker compoase](docker-compose.yml)


## Environmental Variables

.env files are loaded if found in the same directory as oauth-graphql-playground

```
# enable debug logs
OAUTH_GRAPHQL_PLAYGROUND_DEBUG=true

# the port to serve on (default: 5000)
OAUTH_GRAPHQL_PLAYGROUND_PORT=5000

# the oauth2 client id
OAUTH_GRAPHQL_PLAYGROUND_CLIENT_ID=xxx-xxxx-xxxx-xxx

# the oauth2 client secret
OAUTH_GRAPHQL_PLAYGROUND_CLIENT_SECRET=xxx-xxxx-xxxx-xxx

# the redirect url the identity provider will send the user back to(this server)
OAUTH_GRAPHQL_PLAYGROUND_REDIRECT_URL=http://localhost:5000/oauth2/callback

# the oauth2 scopes to ask the user to consent to
OAUTH_GRAPHQL_PLAYGROUND_SCOPES=openid,email,profile

# the oauth2 authorization URL
OAUTH_GRAPHQL_PLAYGROUND_AUTHORIZATION_URL=https://accounts.google.com/o/oauth2/v2/auth

# the oauth2 token URL
OAUTH_GRAPHQL_PLAYGROUND_TOKEN_URL=https://oauth2.googleapis.com/token

# a JSON string used to configure the session manager. options: [cookies]
OAUTH_GRAPHQL_PLAYGROUND_SESSION_MANAGER={ "name": "cookies", "secret": "xxx-xxx-xxx" }

# use open id connect id token on outbound graphQL requests
OAUTH_GRAPHQL_PLAYGROUND_OPEN_ID=true

# the graphQL server to connect to (required)
OAUTH_GRAPHQL_PLAYGROUND_SERVER_ENDPOINT=http://localhost:8080/api/graphql
```

## OAuth Providers

- [Okta]()
- [Google]()
- [Auth0]()
- [Microsoft]()
- [Facebook]()