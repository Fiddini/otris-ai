import express from "express";
import cors from "cors";
import Groq from "groq-sdk";
import { createClient } from "@supabase/supabase-js";
import dotenv from "dotenv";
import slotsRouter from "./routes/slots.js";

dotenv.config();

// Validate environment variables
if (!process.env.GROQ_API_KEY) {
  console.error("❌ ERROR: GROQ_API_KEY tidak ditemukan di .env");
  process.exit(1);
}

if (!process.env.SUPABASE_URL) {
  console.error("❌ ERROR: SUPABASE_URL tidak ditemukan di .env");
  process.exit(1);
}

if (!process.env.SUPABASE_KEY) {
  console.error("❌ ERROR: SUPABASE_KEY tidak ditemukan di .env");
  process.exit(1);
}

if (!process.env.SUPABASE_URL.startsWith("https://")) {
  console.error("❌ ERROR: SUPABASE_URL harus dimulai dengan https://");
  process.exit(1);
}

const app = express();
app.use(cors({
  origin: true, // Allow all origins for development with ngrok
  credentials: true,
}));
app.use(express.json());

// Setup Groq client
const groq = new Groq({ apiKey: process.env.GROQ_API_KEY });

// Setup Supabase client
const supabase = createClient(
  process.env.SUPABASE_URL,
  process.env.SUPABASE_KEY
);

// Import routes
app.use("/api/slot", slotsRouter);

const SYSTEM_PROMPT =
  "Anda adalah LACITA AI EDU, asisten belajar untuk siswa SMA di Provinsi Riau. Jawab soal & jelaskan materi SMA dengan bahasa santai.";

app.post("/api/chat", async (req, res) => {
  try {
    console.log("🔥 CHAT ENDPOINT CALLED");
    const { messages } = req.body;

    if (!messages || !Array.isArray(messages) || messages.length === 0) {
      console.log("❌ ERROR: Messages tidak valid");
      return res.status(400).json({ error: "messages wajib diisi sebagai array" });
    }

    const userMessage = messages[messages.length - 1].content;
    const mode = "belajar";

    // Minimal logging for production
    console.log(`[${new Date().toISOString()}] Chat request: ${userMessage.substring(0, 50)}...`);
    const completion = await groq.chat.completions.create({
      model: "llama-3.3-70b-versatile",
      max_tokens: 200,
      temperature: 0.7,
      messages: [
        { role: "system", content: SYSTEM_PROMPT },
        ...messages,
      ],
    });

    const aiReply = completion.choices[0]?.message?.content || "Maaf, AI gak jawab.";

    // Simpan ke table chats
    const { error: insertError } = await supabase.from("chats").insert({
      message: userMessage,
      mode: mode,
      user_message: userMessage,
      ai_reply: aiReply,
    });

    if (insertError) {
      console.error("Error saving chat:", insertError.message);
    }

    res.json({ reply: aiReply });
  } catch (error) {
    console.error("API Error:", error.message);
    res.status(500).json({ error: "Server error" });
  }
});

app.get("/api/history", async (req, res) => {
  try {
    console.log("=== GET HISTORY ===");
    console.log("Timestamp:", new Date().toISOString());

    const { data, error } = await supabase
      .from("chats")
      .select("*")
      .order("created_at", { ascending: false })
      .limit(20);

    if (error) {
      console.error("History error:", error.message);
      return res.status(500).json({ error: error.message });
    }

    res.json({ history: data || [] });
  } catch (error) {
    console.error("History error:", error.message);
    res.status(500).json({ error: "Server error" });
  }
});

app.listen(3000, () => {
  console.log("OTRIS AI backend jalan di http://localhost:3000");
});
