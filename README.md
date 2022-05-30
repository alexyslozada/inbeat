# INBEAT technical test

This is a simple proof-of-concept.

This project is an engagement rate calculator for instagram.

## Installation and use

To use this project you need Go ^1.16, clone, download the dependencies and run:

```bash
git clone git@github.com:alexyslozada/inbeat.git
cd inbeat
go mod tidy
go run .
```

Then you can open a browser and go to `http://localhost:8080`.

Now you can search an IG profile by the username. You can use `@username` or `username`.

## Dependencies

We are using `echo labstack` for routing and `jsonparser` for read the body response from IG. `jsonparser` is faster than `json.Unmarshal`.

## Issues

IG can block you when you are using a datacenter proxy or when you are making several requests. The firsts three request will not be blocked, but after it, you could be blocked. To fix it, you can use a residential proxy.
