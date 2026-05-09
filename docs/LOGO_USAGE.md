# Logo Usage Guide - OTRIS AI

## File Logo yang Tersedia

### 1. `logo.svg` - Full Logo dengan Text
- **Ukuran**: 200×250px
- **Penggunaan**: Login page, splash screen, dashboard header
- **Scalable**: Ya, bisa di-resize sesuai kebutuhan

### 2. `logo-icon.svg` - Icon Only
- **Ukuran**: 64×64px
- **Penggunaan**: Favicon, navbar icon, sidebar toggle
- **Scalable**: Ya, cocok untuk 16px hingga 128px

### 3. `logo-horizontal.svg` - Horizontal Layout
- **Ukuran**: 280×80px
- **Penggunaan**: Navbar, header branding, login header
- **Scalable**: Ya, responsive width

### 4. `logo-compact.svg` - Minimal/Compact Version
- **Ukuran**: 120×120px
- **Penggunaan**: Sidebar, mobile navbar, small displays
- **Scalable**: Ya, cocok untuk 48px hingga 200px

## Penggunaan di React

### Navbar
```jsx
<nav className="navbar">
  <img src="/logo-horizontal.svg" alt="OTRIS AI" style={{height: "50px"}} />
</nav>
```

### Login Page
```jsx
<div className="login-container">
  <img src="/logo.svg" alt="OTRIS AI" style={{width: "150px", height: "auto"}} />
  <h1>Welcome to OTRIS AI</h1>
</div>
```

### Favicon
```html
<!-- Di index.html -->
<link rel="icon" href="/logo-icon.svg" type="image/svg+xml" />
```

### Dashboard Header
```jsx
<header>
  <img src="/logo-icon.svg" alt="OTRIS" style={{width: "40px", height: "40px"}} />
  <span>Dashboard</span>
</header>
```

### Sidebar
```jsx
<aside className="sidebar">
  <img src="/logo-compact.svg" alt="OTRIS" style={{width: "80px"}} />
</aside>
```

## Fitur Design

✅ **Glow Effect** - Efek cahaya halus untuk modern look
✅ **Font Modern** - Menggunakan Poppins/Inter fallback
✅ **Scalable** - SVG responsive untuk semua ukuran
✅ **Clean Structure** - Code SVG teroptimasi untuk performa web
✅ **Gradient** - Gradien biru modern (#00D9FF → #0052CC)
✅ **Shadow** - Soft shadow untuk depth

## Rekomendasi Ukuran

| Lokasi | Logo | Ukuran Min | Ukuran Ideal |
|--------|------|-----------|------------|
| Favicon | logo-icon.svg | 16px | 32×32px |
| Navbar | logo-icon.svg | 32px | 40×40px |
| Navbar Branding | logo-horizontal.svg | 120px | 240px |
| Login Page | logo.svg | 150px | 200px |
| Dashboard Header | logo-icon.svg | 40px | 50×50px |
| Sidebar | logo-compact.svg | 60px | 80×80px |
| Mobile Menu | logo-icon.svg | 32px | 48×48px |

## Kustomisasi Warna (Optional)

Jika ingin mengubah warna, edit stop-color di dalam `<linearGradient>`:
- `#00D9FF` → Warna terang (cyan)
- `#0052CC` → Warna gelap (blue)

## Performa

- Semua file SVG sudah dioptimasi
- Ukuran file minimal (~1-2KB per file)
- Cocok untuk retina display dan responsive design
- Compatible dengan semua browser modern
