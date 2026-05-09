<<<<<<< HEAD
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
=======
```
1. Isi file README.md (jika belum ada)

Buat file bernama README.md di repo kamu dengan isi singkat berikut sebagai petunjuk penggunaan:
```
# OTRIS AI Backend

### Cara menjalankan lokal:

```bash
pip install -r requirements.txt
uvicorn main:app --reload --port 8000
```
Endpoint:

GET /health → cek status server
POST /api/chat → kirim JSON dengan properti:
message: string (pertanyaan)
mode: string (belajar, soal, ringkas)
Contoh POST request:
```
{
  "message": "Jelaskan fotosintesis",
  "mode": "belajar"
}
```
2. Siapkan Deployment ke Railway

Login ke Railway.app
Buat project baru → hubungkan dengan repo GitHub kamu
Set environment variable:
OPENAI_API_KEY dengan API key OpenAI-mu
3. Deploy dan Testing

Railway akan build dan deploy otomatis berdasarkan requirements.txt dan main.py
Setelah deploy selesai, buka URL yang diberikan Railway, misal:

```
https://namaproject.up.railway.app/health
```
Pastikan health endpoint merespon dengan {"status":"OTRIS AI OK"}
4. Tes API chat

Gunakan Postman / curl untuk coba endpoint chat misal:

```
curl -X POST https://namaproject.up.railway.app/api/chat \
    -H "Content-Type: application/json" \
    -d '{"message":"Jelaskan rumus luas lingkaran","mode":"belajar"}'
```
5. Update URL di aplikasi Android

Ganti URL backend di aplikasi kamu ke URL Railway agar aplikasi bisa berkomunikasi dengan backend.

Jika sudah ada repo dan filenya lengkap, langkah utama sekarang adalah:

Deploy ke Railway
Set API KEY
Test endpoint /health dan /api/chat
Update URL backend di aplikasi
Kalau mau, saya bisa bantu buatkan README.md lengkap untuk repo kamu.

Mau?
```





>>>>>>> origin/main
