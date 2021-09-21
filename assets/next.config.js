module.exports = {
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: "http://localhost:9080/api/:path*",
      },
    ];
  },
};
