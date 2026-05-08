# OTRIS AI - Asisten Pembelajaran Indonesia

OTRIS AI adalah platform asisten pembelajaran berbasis AI untuk sekolah Indonesia.

## Teknologi

- Backend: FastAPI
- Frontend: React + Vite
- AI Engine: OpenAI API

## Struktur Project

```
.
├── backend/             # FastAPI backend
├── frontend/            # React/Vite frontend
├── docs/                # Dokumentasi proyek
├── README.md            # Dokumentasi utama
├── .gitignore           # Aturan file yang diabaikan
└── .env                 # Environment variables (gitignored)
```

## Backend

Backend menggunakan FastAPI untuk menyediakan endpoint chat AI.

### Jalankan Backend

1. Masuk ke folder backend:

```bash
cd backend
```

2. Install dependency:

```bash
pip install -r requirements.txt
```

3. Siapkan file `.env` di root project atau di folder backend dengan isi:

```env
OPENAI_API_KEY=sk-proj-your_actual_openai_api_key_here
```

4. Jalankan server:

```bash
uvicorn ai_tutor:app --reload --port 8002
```

Backend akan berjalan di `http://127.0.0.1:8002`

## Frontend

Frontend menggunakan React dan Vite untuk UI interaktif layar besar.

### Jalankan Frontend

1. Masuk ke folder frontend:

```bash
cd frontend
```

2. Install dependency:

```bash
npm install
```

3. Jalankan dev server:

```bash
npm run dev -- --host
```

Frontend akan tersedia di `http://localhost:5174` atau port lain jika 5173 sudah digunakan.

## OpenAI API

OTRIS AI menggunakan OpenAI API untuk merespons pertanyaan pembelajaran.

Pastikan `OPENAI_API_KEY` sudah diatur di file `.env`.

## Catatan

- File `.env` tidak boleh ditambahkan ke repository.
- `node_modules` dan `dist` juga harus diabaikan oleh Git.

## Dokumentasi

Lihat folder `docs/` untuk dokumentasi tambahan.
