from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from openai import OpenAI
import os

app = FastAPI()
app.add_middleware(CORSMiddleware, allow_origins=["*"], allow_methods=["*"], allow_headers=["*"])

client = OpenAI(api_key=os.getenv("OPENAI_API_KEY"))

class ChatRequest(BaseModel):
    message: str
    mode: str = "belajar"

@app.post("/api/chat")
async def chat(request: ChatRequest):
    try:
        response = client.chat.completions.create(
            model="gpt-4o-mini",
            messages=[
                {"role": "system", "content": "Anda adalah OTRIS AI. Gunakan bahasa Indonesia sederhana untuk siswa sekolah. Format: 1. Penjelasan singkat 2. Contoh 3. Soal latihan. Mulai dengan 'OTRIS AI:'"},
                {"role": "user", "content": f"{request.mode}: {request.message}"}
            ]
        )
        return {"reply": response.choices[0].message.content}
    except:
        raise HTTPException(500, "Server error")

@app.get("/health")
async def health():
    return {"status": "OTRIS AI OK"}
