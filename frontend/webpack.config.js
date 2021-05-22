const baseConfig = {
  mode: "development",
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
  target: ["web", "es5"],
};

const devServer = {
  contentBase: "dist",
  proxy: {
    "/api": "http://localhost:9080",
  },
  hot: true,
};

module.exports = [
  {
    ...baseConfig,
    ...{
      name: "public",
      entry: `./src/public/main.tsx`,
      output: {
        path: `${__dirname}/dist/public`,
        filename: "bundle.js",
      },
      devServer: {
        ...devServer,
        ...{
          contentBase: "dist/public",
          port: 8080,
        },
      },
    },
  },
  {
    ...baseConfig,
    ...{
      name: "private",
      entry: `./src/private/main.tsx`,
      output: {
        path: `${__dirname}/dist/private`,
        filename: "bundle.js",
      },
      devServer: {
        ...devServer,
        ...{
          contentBase: "dist/private",
          port: 8081,
        },
      },
    },
  },
];
