1. Buka: https://railway.app/new
2. Login (GitHub/Google)
3. "New Project" → "Deploy from GitHub repo"
4. Cari repo: otis-ai-backend → Select → Deploy
```
5. Tunggu build (1-2 menit)
6. Status 🟢 → Catat URL: https://otris-ai-xxx.up.railway.app
```
7. SET API KEY (KRITIS!)

```
Railway Dashboard → Variables tab → "+ New Variable"
Name: OPENAI_API_KEY  
Value: sk-proj-your-complete-key-here
Deploy → Redeploy
```

TASK 2: TEST LIVE (1 menit)


```
Browser: https://your-url/health
✅ {"status": "OTRIS AI OK"}

Postman/curl:
POST https://your-url/api/chat
{
  "message": "test OTRIS AI",
  "mode": "belajar"
}
✅ {"reply": "OTRIS AI: ..."}
Browser: https://your-url/health
✅ {"status": "OTRIS AI OK"}

```

TASK 3: ANDROID (Habib nanti)

```
Ganti 1 baris:
.url("https://your-railway-url/api/chat")
Build APK → DONE!
```
SUCCESS = 3 Tanda

```
1. Railway 🟢 Deployed
2. /health = OK  
3. /api/chat = OTRIS AI reply
```
DO NOW

```
Railway.app → Deploy → API Key → Test health
Kirim gw: 
- Railway URL
- Screenshot /health
```

