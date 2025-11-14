# Medee Project

This project consists of two main parts:
- **Backend**: Go (Golang) API server
- **Frontend**: React (TypeScript) app using Vite

---

## Prerequisites

- **Node.js** (recommended v18+)
- **pnpm** (recommended v8+)
- **Go** (recommended v1.20+)

---

## Backend Setup (Go)

1. Open a terminal and navigate to the `backend` folder:
	```sh
	cd backend
	```
2. Run the Go server:
	```sh
	go run server.go
	```
	The backend server will start (default port is usually 8080 unless changed in code).

---

## Frontend Setup (React + Vite)

1. Open a terminal and navigate to the `frontend` folder:
	```sh
	cd frontend
	```
2. Install dependencies:
	```sh
	pnpm install
	```
3. Start the development server:
	```sh
	pnpm dev
	```
	The frontend will be available at [http://localhost:5173](http://localhost:5173) by default.

---

## Notes
- Make sure the backend server is running before using the frontend if the frontend fetches data from the backend.
- You may need to adjust CORS or API URLs in the frontend to match your backend server address.

---

## Project Structure

```
Medee/
├── backend/      # Go backend API
│   ├── go.mod
│   └── server.go
└── frontend/     # React frontend (Vite)
	 ├── src/
	 ├── public/
	 ├── package.json
	 └── ...
```