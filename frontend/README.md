# OTRIS AI Frontend

Frontend React + Tailwind untuk OTRIS AI, dirancang sebagai UI fullscreen untuk layar kelas besar.

## Fitur

- React + Vite + Tailwind CSS
- Fullscreen touch-friendly layout
- Tema futuristik dark blue
- Menu utama: Belajar, Ringkasan, Quiz, Contoh, Tanya AI
- Koneksi langsung ke backend: `POST http://127.0.0.1:8002/api/chat`
- Tombol besar untuk send, clear, fullscreen, dan voice input placeholder
- Panel jawaban AI yang luas dan mudah dibaca

## Setup

1. Install dependencies:

```bash
cd ai-tutor-ui
npm install
```

2. Jalankan development server:

```bash
npm run dev -- --host
```

3. Buka browser:

```text
http://127.0.0.1:5173
```

## Struktur Penting

- `src/App.jsx`: antarmuka utama OTRIS AI
- `src/index.css`: gaya Tailwind dan tema fullscreen
- `tailwind.config.js`: konfigurasi Tailwind
- `postcss.config.js`: konfigurasi PostCSS

## Backend

Pastikan backend OTRIS AI berjalan di:

```text
http://127.0.0.1:8002/api/chat
```

## Perhatian

- UI ini menggunakan endpoint backend lokal, jadi backend harus aktif saat testing.
- Jika ingin memindahkan backend, perbarui URL pada `src/App.jsx`.
