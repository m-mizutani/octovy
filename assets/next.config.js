module.exports = {
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: "http://127.0.0.1:9080/api/:path*",
      },
      {
        source: "/auth/:path*",
        destination: "http://127.0.0.1:9080/auth/:path*",
      },
    ];
  },
};
