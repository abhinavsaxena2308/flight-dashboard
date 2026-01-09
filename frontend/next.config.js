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
        destination: 'https://flight-dashboard-2.onrender.com/api/:path*' , // Proxy API requests to backend Go server
      },
    ];
  },
  env: {}, // This doesn't affect proxying
  // Enable trailing slash handling if needed
  trailingSlash: undefined,
  webpack(config) {
    // Further configure Webpack if needed
    return config;
  },
}

module.exports = nextConfig