# Pokerer simple frontend

Plain HTML + CSS + JavaScript frontend for the existing Pokerer Go API.

No React, Vite, TypeScript, or additional frontend server is required.

## Run

Start the Pokerer Go API first.

Then serve this directory with any static HTTP server. For example:

```bash
python3 -m http.server 5173
```

Open:

```text
http://localhost:5173
```

The API URL is currently:

```text
http://localhost:8080
```

Change `API` at the top of `app.js` if your Go server uses another address.

## API calls

This frontend uses:

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/me`
- `GET /api/v1/wallet`
- `GET /api/v1/tables`
- `GET /api/v1/tables/{id}`
- `POST /api/v1/tables/{id}/join`
- `POST /api/v1/tables/{id}/leave`

The JWT is stored in localStorage and sent as a Bearer token.

If the browser blocks requests because the frontend and API are on different origins, enable CORS in the Go API for `http://localhost:5173`.

The poker table itself does not fake live cards/actions. It displays the real persisted table/player data available from the current API.
