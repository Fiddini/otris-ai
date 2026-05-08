import { useState, useMemo } from 'react'

const modes = [
  { key: 'belajar', label: 'Belajar' },
  { key: 'ringkas', label: 'Ringkasan' },
  { key: 'soal', label: 'Quiz' },
  { key: 'contoh', label: 'Contoh' },
  { key: 'tanya', label: 'Tanya AI' },
]

const modeTitles = {
  belajar: 'Penjelasan Belajar',
  ringkas: 'Ringkasan Cepat',
  soal: 'Buat Quiz',
  contoh: 'Contoh Sederhana',
  tanya: 'Tanya AI',
}

function App() {
  const [activeMode, setActiveMode] = useState('belajar')
  const [question, setQuestion] = useState('')
  const [answer, setAnswer] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [fullscreen, setFullscreen] = useState(false)

  const modeLabel = useMemo(
    () => modeTitles[activeMode] || 'Belajar',
    [activeMode],
  )

  const handleClear = () => {
    setQuestion('')
    setAnswer('')
    setError('')
  }

  const toggleFullscreen = async () => {
    if (document.fullscreenElement) {
      await document.exitFullscreen()
      setFullscreen(false)
    } else {
      await document.documentElement.requestFullscreen()
      setFullscreen(true)
    }
  }

  const handleSubmit = async () => {
    if (!question.trim()) {
      setError('Silakan masukkan pertanyaan terlebih dahulu.')
      return
    }

    setLoading(true)
    setError('')
    setAnswer('')

    try {
      const response = await fetch('http://127.0.0.1:8002/api/chat', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          message: question,
          mode: activeMode,
        }),
      })

      if (!response.ok) {
        const result = await response.json().catch(() => null)
        throw new Error(result?.detail || 'Terjadi kesalahan server.')
      }

      const data = await response.json()
      setAnswer(data.reply || 'AI tidak memberikan jawaban.')
    } catch (err) {
      setError(err.message || 'Kesalahan tidak terduga.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-midnight text-slate-100">
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_top,_rgba(56,189,248,0.18),_transparent_38%),radial-gradient(circle_at_bottom,_rgba(59,130,246,0.16),_transparent_30%),linear-gradient(135deg,#020617_0%,#071026_100%)]" />
      <div className="relative flex min-h-screen flex-col px-8 py-6 lg:px-12">
        <header className="mb-6 flex flex-col gap-4 rounded-3xl border border-white/10 bg-white/5 p-6 shadow-[0_0_120px_rgba(59,130,246,0.08)] backdrop-blur-xl lg:flex-row lg:items-center lg:justify-between">
          <div>
            <p className="mb-2 text-sm uppercase tracking-[0.3em] text-sky-300/80">
              OTRIS AI Classroom
            </p>
            <h1 className="max-w-3xl text-4xl font-semibold tracking-tight text-slate-100 lg:text-5xl">
              Asisten belajar interaktif untuk layar kelas besar.
            </h1>
            <p className="mt-3 max-w-2xl text-base text-slate-300 sm:text-lg">
              Pilih menu, tulis pertanyaan, dan dapatkan jawaban AI dalam panel besar. Dirancang untuk touchscreen 86 inci.
            </p>
          </div>

          <div className="grid gap-3 sm:grid-cols-3">
            <button type="button" className="rounded-3xl border border-sky-500/20 bg-sky-500/10 px-6 py-4 text-left text-slate-100 transition hover:bg-sky-500/20">
              <span className="block text-sm text-sky-200">Mode active</span>
              <span className="mt-1 block text-xl font-semibold">{modeLabel}</span>
            </button>
            <button type="button" onClick={toggleFullscreen} className="rounded-3xl border border-slate-600/40 bg-slate-900/80 px-6 py-4 text-left text-slate-100 transition hover:bg-slate-800/90">
              <span className="block text-sm text-slate-400">Tombol Fullscreen</span>
              <span className="mt-1 block text-xl font-semibold">{fullscreen ? 'Keluar Fullscreen' : 'Fullscreen'}</span>
            </button>
            <button type="button" className="rounded-3xl border border-fuchsia-500/20 bg-fuchsia-500/10 px-6 py-4 text-left text-slate-100 transition hover:bg-fuchsia-500/20">
              <span className="block text-sm text-fuchsia-200">Voice Input</span>
              <span className="mt-1 block text-xl font-semibold">Placeholder</span>
            </button>
          </div>
        </header>

        <main className="grid flex-1 gap-6 lg:grid-cols-[360px_minmax(0,1fr)]">
          <aside className="space-y-6 rounded-[2rem] border border-white/10 bg-white/5 p-6 shadow-[0_0_80px_rgba(30,64,175,0.08)] backdrop-blur-xl">
            <div className="space-y-3">
              <h2 className="text-lg font-semibold tracking-wide text-sky-200">Menu Utama</h2>
              <p className="text-sm text-slate-400">Pilih kategori tugas yang akan dijalankan oleh AI.</p>
            </div>

            <div className="grid gap-3">
              {modes.map((mode) => (
                <button
                  key={mode.key}
                  type="button"
                  onClick={() => setActiveMode(mode.key)}
                  className={`rounded-3xl border px-5 py-4 text-left transition ${activeMode === mode.key ? 'border-sky-300/80 bg-sky-400/10 text-sky-100 shadow-[0_0_30px_rgba(56,189,248,0.18)]' : 'border-slate-700/80 bg-slate-950/80 text-slate-300 hover:border-sky-300/30 hover:bg-slate-900/90'}`}>
                  <p className="text-sm uppercase tracking-[0.3em] text-slate-400">{mode.label}</p>
                  <p className="mt-2 text-xl font-semibold">{modeTitles[mode.key]}</p>
                </button>
              ))}
            </div>

            <div className="rounded-3xl border border-slate-700/80 bg-slate-950/90 p-5">
              <h3 className="text-sm font-semibold uppercase tracking-[0.3em] text-sky-300/90">Voice Input</h3>
              <p className="mt-3 text-sm leading-6 text-slate-400">Placeholder untuk kontrol suara. Fungsi dimatikan untuk saat ini.</p>
            </div>

            <div className="grid gap-3">
              <button type="button" onClick={handleClear} className="rounded-3xl border border-rose-500/20 bg-rose-500/10 px-5 py-4 text-slate-100 transition hover:bg-rose-500/20">
                Clear Semua
              </button>
              <button type="button" onClick={toggleFullscreen} className="rounded-3xl border border-sky-500/20 bg-sky-500/10 px-5 py-4 text-slate-100 transition hover:bg-sky-500/20">
                {fullscreen ? 'Keluar Fullscreen' : 'Masuk Fullscreen'}
              </button>
            </div>
          </aside>

          <section className="space-y-6">
            <div className="rounded-[2rem] border border-white/10 bg-white/5 p-6 shadow-[0_0_80px_rgba(30,64,175,0.08)] backdrop-blur-xl">
              <div className="mb-5 flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
                <div>
                  <p className="text-sm uppercase tracking-[0.3em] text-sky-300/80">Pertanyaan</p>
                  <h2 className="mt-2 text-3xl font-semibold text-slate-100">Tulis pertanyaanmu untuk OTRIS AI</h2>
                </div>
                <div className="text-right text-sm text-slate-400">
                  Mode aktif <span className="font-semibold text-slate-100">{modeLabel}</span>
                </div>
              </div>

              <textarea
                value={question}
                onChange={(event) => setQuestion(event.target.value)}
                rows={6}
                placeholder="Ketik pertanyaan di sini..."
                className="w-full rounded-3xl border border-slate-700/80 bg-slate-950/90 px-5 py-5 text-lg leading-7 text-slate-100 outline-none transition focus:border-sky-400/80 focus:ring-2 focus:ring-sky-500/10"
              />

              <div className="mt-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <div className="space-y-2">
                  <p className="text-sm text-slate-400">Backend endpoint:</p>
                  <code className="block rounded-3xl border border-slate-700/80 bg-slate-950/80 px-4 py-2 text-sm text-slate-300">POST http://127.0.0.1:8002/api/chat</code>
                </div>
                <button
                  type="button"
                  onClick={handleSubmit}
                  disabled={loading}
                  className="inline-flex items-center justify-center rounded-3xl bg-sky-500 px-8 py-4 text-lg font-semibold text-slate-950 transition hover:bg-sky-400 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  {loading ? 'Mengirim...' : 'Kirim ke AI'}
                </button>
              </div>

              {error ? (
                <div className="mt-5 rounded-3xl border border-rose-500/20 bg-rose-500/10 px-5 py-4 text-rose-100">
                  {error}
                </div>
              ) : null}
            </div>

            <div className="rounded-[2rem] border border-white/10 bg-gradient-to-br from-slate-950/90 via-slate-900/80 to-sky-950/80 p-6 shadow-[0_0_90px_rgba(59,130,246,0.14)] backdrop-blur-xl">
              <div className="mb-4 flex items-center justify-between gap-4">
                <div>
                  <p className="text-sm uppercase tracking-[0.3em] text-sky-300/80">Jawaban AI</p>
                  <h2 className="mt-2 text-3xl font-semibold text-slate-100">Panel Tanggapan</h2>
                </div>
                <span className="rounded-full bg-slate-900/90 px-4 py-2 text-sm text-slate-300">{answer ? 'Selesai' : 'Menunggu input'}</span>
              </div>

              <div className="min-h-[320px] rounded-[2rem] border border-slate-700/70 bg-slate-950/90 p-6 text-lg leading-8 text-slate-200 shadow-inner shadow-slate-950/30">
                {answer ? (
                  <p className="whitespace-pre-line">{answer}</p>
                ) : (
                  <p className="text-slate-500">Jawaban AI akan muncul di sini setelah Anda mengirim pertanyaan.</p>
                )}
              </div>
            </div>
          </section>
        </main>
      </div>
    </div>
  )
}

export default App
