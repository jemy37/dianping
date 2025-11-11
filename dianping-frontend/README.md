## Dianping Frontend (React + Vite)

This is a lightweight frontend to interact with the backend under `dianping/`.

### Prerequisites
- Node.js 18+

### Install & Run
```
npm install
npm run dev
```

The dev server runs on `http://localhost:5173` and proxies API requests to the backend at `http://localhost:8080`.

Update the proxy target in `vite.config.ts` if your backend runs on a different port.

### Scripts
- `npm run dev` — start dev server on port 5173
- `npm run build` — production build
- `npm run preview` — preview production build

