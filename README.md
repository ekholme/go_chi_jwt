# Go Chi JWT

A demo application using the [chi](https://github.com/go-chi/chi) router and a JWT middleware for authentication

# Jwt process

Ok so I'm pretty sure the jwt process is to:

- When a user logs in, generate a jwt and then save it in a cookie
- When a user attempts to access a page, read the cookie and then validate it in a middleware