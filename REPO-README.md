# Monorepo Overview / 仓库总览

This repository can be organized as a simple monorepo with two top-level folders:

- backend/ — Go (Gin + GORM + Redis) service
- frontend/ — React + Vite app

If your current code hasn’t been moved yet, see the steps below to restructure. After moving, both folders build and run independently but work together via the Vite dev proxy.

## Structure / 目录结构

```
repo-root/
├── backend/              # Go backend (current code)
│   ├── config/
│   ├── dao/
│   ├── handler/
│   ├── models/
│   ├── router/
│   ├── script/
│   ├── service/
│   ├── utils/
│   ├── main.go
│   └── README.md
└── frontend/             # React + Vite frontend
    ├── src/
    ├── index.html
    ├── package.json
    ├── vite.config.ts
    └── README.md
```

## How To Restructure / 如何重构为双目录

From repo root:

- Create backend folder and move tracked files (preserves history):
  - PowerShell
    - `mkdir backend`
    - `git ls-files -z | % { $_ -replace "`0","" } | % { $dest = Join-Path backend $_; New-Item (Split-Path $dest) -ItemType Directory -Force | Out-Null; git mv $_ $dest }`
- Move untracked files into backend (if any):
  - `Get-ChildItem -Force | Where-Object { $_.Name -notin @('backend','.git','.gitignore','REPO-README.md') } | ForEach-Object { $dest = Join-Path 'backend' $_.Name; Move-Item $_.FullName $dest -Force }`
- Bring the existing frontend into the repo:
  - If it exists as sibling `../dianping-frontend`: `Move-Item ..\dianping-frontend .\frontend`
- Commit:
  - `git add -A`
  - `git commit -m "Restructure: backend/ + frontend/ monorepo"`

## Backend / 后端

- Language: Go 1.21+
- Stack: Gin, GORM, Redis (Bloom/GEO/Stream)
- Run / 运行：
  - `cd backend`
  - `go mod tidy`
  - `go run main.go`
- Default server: `http://localhost:8080`

## Frontend / 前端

- Stack: React 18 + Vite + React Router + Axios
- Dev / 开发：
  - `cd frontend`
  - `npm install`
  - `npm run dev`
- Dev server: `http://localhost:5173` (proxies `/api` → `http://localhost:8080`)
  - If backend port differs, edit `frontend/vite.config.ts` proxy target.

## Git Ignore / 忽略建议

A combined `.gitignore` is included for Go + Vite:

- Go: `bin/`, `tmp/`, `coverage.out`, `*.test`
- Node/Vite: `node_modules/`, `dist/`, `.vite/`, `*.tsbuildinfo`
- General: `.env`, `.vscode/`, `.idea/`, `*.log`

## CI (Optional) / 可选 CI

- Add separate jobs for `backend/` (Go build/test) and `frontend/` (Node build)
- Cache Go module and Node dependencies to speed up pipelines

## Notes / 备注

- Frontend talks to backend via `/api` paths; the Vite proxy avoids CORS in dev.
- Keep backend CORS middleware enabled for production cross-origin needs.
- Ensure `JWT secret` and DB credentials are set via env in production.

