module.exports = {
  mode: "development",
  entry: "./src/internal/main.tsx",
  output: {
    path: `${__dirname}/dist`,
    filename: "bundle.js",
  },
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: "ts-loader",
      },
    ],
  },
  resolve: {
    extensions: [".ts", ".tsx", ".js", ".json"],
  },
  devServer: {
    contentBase: "dist",
    proxy: {
      "/api": "http://localhost:9080",
    },
    hot: true,
  },
  target: ["web", "es5"],
};
