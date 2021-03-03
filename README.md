# oauth-graphql-playground

An oauth2 protected [graphQL playground](https://github.com/graphql/graphql-playground)

## Features

- [x] Serves [GraphQL Playground](https://github.com/graphql/graphql-playground) user interface `/`
- [x] Login with [oauth authorization code grant](https://oauth.net/2/grant-types/authorization-code/)
    - automatically redirects the user to login if token is expired or cannot be refreshed
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

You will need to register an OAuth client application with an identity provider if you havent already.
Please note that your OAuth config should be setup as a "Web Application" with the "Authorization Code Grant" enabled.
You also may need to do additional configuration of your OAuth app depending on your configured scopes.


- [Google](https://support.google.com/googleapi/answer/6158849?hl=en)
    - token url: https://oauth2.googleapis.com/token
    - authorization url: https://accounts.google.com/o/oauth2/v2/auth

- [Microsoft Azure AD](https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-auth-code-flow)
    - token url: https://login.microsoftonline.com/${tenant}/oauth2/v2.0/token
    - authorization url: https://login.microsoftonline.com/${tenant}/oauth2/v2.0/authorize

- [Okta](https://developer.okta.com/docs/guides/implement-oauth-for-okta/create-oauth-app/)
    - token url: todo
    - authorization url: todo

- [Auth0](https://auth0.com/docs/applications/set-up-an-application)
    - token url: https://login.microsoftonline.com/${tenant}/oauth2/v2.0/token
    - authorization url: https://login.microsoftonline.com/${tenant}/oauth2/v2.0/authorize
    
- [Facebook](https://developers.facebook.com/docs/facebook-login/)
    - token url: https://graph.facebook.com/v3.2/oauth/access_token
    - authorization url: https://www.facebook.com/v3.2/dialog/oauth

- [Slack](https://api.slack.com/legacy/oauth)
    - token url: https://slack.com/api/oauth.access
    - authorization url: https://slack.com/oauth/authorize
    
- [Github](https://github.com/settings/applications/new)
    - token url: https://github.com/login/oauth/access_token
    - authorization url: https://github.com/login/oauth/authorize