// next.config.js
/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    domains: ['images.unsplash.com', 'via.placeholder.com'],
  },
  experimental: {
    swcPlugins: [],
  },
  // Use Babel instead of SWC
  swcMinify: false,
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/api/:path*' || 'https://flight-dashboard-2.onrender.com/', // Proxy API requests to backend Go server
      },
    ];
  },
  /* For API routes, we need to handle proxying differently */
  env: {}, // This doesn't affect proxying
  // Enable trailing slash handling if needed
  trailingSlash: undefined,
  webpack(config) {
    // Further configure Webpack if needed
    return config;
  },
}

module.exports = nextConfig