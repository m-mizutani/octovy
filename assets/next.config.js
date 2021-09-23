module.exports = {
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: "http://localhost:9080/api/:path*",
      },
      {
        source: "/auth/:path*",
        destination: "http://localhost:9080/auth/:path*",
      },
    ];
  },
};
