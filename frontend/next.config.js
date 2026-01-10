// next.config.js
/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    swcPlugins: [],
  },
  swcMinify: false,
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'https://flight-dashboard-2.onrender.com/api/:path*' , // Proxy API requests to backend Go server
      },
    ];
  },
  env: {}, 
  trailingSlash: undefined,
  webpack(config) {
    return config;
  },
}

module.exports = nextConfig