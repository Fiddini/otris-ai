from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from openai import OpenAI
import os
import uvicorn
from dotenv import load_dotenv

load_dotenv()

app = FastAPI(title="OTRIS AI - Asisten Pembelajaran Indonesia")

# CORS Middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Inisialisasi OpenAI Client menggunakan API key dari environment variable
client = OpenAI(api_key=os.getenv("OPENAI_API_KEY"))

SYSTEM_PROMPT = """
Anda adalah OTRIS AI — Asisten Pembelajaran untuk Sekolah Indonesia.
Gunakan bahasa Indonesia sederhana dan jelas.
Selalu mulai dengan: "OTRIS AI – Asisten Pembelajaran:"
Jawab sesuai permintaan siswa dengan ramah dan sabar.
"""

class ChatRequest(BaseModel):
    message: str
    mode: str = "belajar"

@app.get("/")
async def root():
    return {"message": "OTRIS AI is running! 🚀", "api": "/api/chat"}

@app.get("/api/health")
async def health():
    return {"status": "OK", "service": "OTRIS AI"}

@app.post("/api/chat")
async def chat(request: ChatRequest):
    try:
        if not request.message or request.message.strip() == "":
            raise HTTPException(status_code=400, detail="Pesan tidak boleh kosong!")
        
        mode_instruction = {
            "belajar": "Jelaskan dengan sederhana dan mudah dipahami.",
            "soal": "Buat 3-5 soal pilihan ganda + jawaban kunci.",
            "ringkas": "Ringkas poin-poin penting secara singkat.",
            "contoh": "Berikan contoh sederhana + penjelasan."
        }.get(request.mode.lower(), "Jelaskan materi dengan baik.")
        
        full_prompt = f"{mode_instruction}\n\nTopik/Pertanyaan: {request.message}"
        
        response = client.chat.completions.create(
            model="gpt-4o-mini",
            messages=[
                {"role": "system", "content": SYSTEM_PROMPT},
                {"role": "user", "content": full_prompt}
            ],
            max_tokens=1000,
            temperature=0.7
        )
        
        return {
            "reply": response.choices[0].message.content,
            "mode": request.mode,
            "success": True
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error: {str(e)}")

if __name__ == "__main__":
    uvicorn.run("ai_tutor:app", host="0.0.0.0", port=8002, reload=True)
