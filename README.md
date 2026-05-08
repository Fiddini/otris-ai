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





