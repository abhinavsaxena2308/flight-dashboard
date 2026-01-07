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
}

module.exports = nextConfig