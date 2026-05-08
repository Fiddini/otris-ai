export default {
  content: ['./index.html', './src/**/*.{js,jsx,ts,tsx}'],
  theme: {
    extend: {
      boxShadow: {
        glow: '0 0 90px rgba(56, 189, 248, 0.18)',
      },
      colors: {
        midnight: '#020617',
      },
      backgroundImage: {
        'gradient-futuristic': 'radial-gradient(circle at top, rgba(56,189,248,0.16), transparent 28%), linear-gradient(180deg, #020617 0%, #06101e 100%)',
      },
    },
  },
  plugins: [],
}
