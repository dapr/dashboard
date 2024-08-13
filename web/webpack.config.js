const TerserPlugin = require("terser-webpack-plugin");

module.exports = {
  optimization: {
    minimize: true,
    minimizer: [
      new TerserPlugin({}),
    ],
    runtimeChunk: "single",
    splitChunks: {
      chunks: "all",
      minChunks: 1,
      maxInitialRequests: Infinity,
      minSize: 30000,
      cacheGroups: {
        vendor: {
          test: /[\\/]node_modules[\\/]/,
          name(module) {
            const packageName = module.context.match(
              /[\\/]node_modules[\\/](.*?)([\\/]|$)/
            )[1];
            return `dapr.${packageName.replace("@", "")}`;
          },
        }
      }
    }
  }
};
