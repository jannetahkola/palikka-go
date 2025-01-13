# Palikka Go
Rewrite of Palikka in Golang. Not finished, only user management, authn & authz are done(-ish)

Project was last updated during 2024.

Tried a different project structure based on a repo I found which needs to be linked here (todo find it). 
Idea is to layer the project so that domain files are in the root, which any layer can easily depend on.
We can then have a structure like (example):

domain <- data layer <- HTTP server
domain <- data layer <- some other layer <- CLI

Any layer in-between is meant to provide an API so that the actual implementation can be switched, e.g. a 
different storage method for HTTP sessions.

Finally there is the cmd package that initialises and binds HTTP server and CLI layers together.

# Dependencies

## modernc.org/sqlite
SQLite driver - pure Go so trade in performance to avoid CGo dependency.

## sqlc
> Only needed for development

todo

todo remove .idea folder from git
