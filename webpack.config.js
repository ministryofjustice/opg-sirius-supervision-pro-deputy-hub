var path = require("path");
var CopyPlugin = require("copy-webpack-plugin");
var MiniCssExtractPlugin = require("mini-css-extract-plugin");

module.exports = {
    entry: {
        all: "./web/assets/main.js",
    },
    mode: "production",
    devtool: "source-map",
    module: {
        rules: [
            {
                test: /\.scss$/i,
                use: [
                    {
                        loader: MiniCssExtractPlugin.loader,
                    },
                    {
                        loader: "css-loader",
                        options: {
                            url: false,
                        },
                    },
                    {
                        loader: "sass-loader",
                    },
                ],
            },
            {
                test: /\.js$/,
                exclude: /node_modules/,
                loader: "babel-loader",
            },
        ],
    },
    output: {
        filename: "javascript/[name].js",
        path: path.resolve(__dirname, "web/static"),
    },
    plugins: [
        new CopyPlugin({
            patterns: [
                {
                    from: "node_modules/govuk-frontend/govuk/assets/images",
                    to: path.resolve(__dirname, "web/static/assets/images"),
                },
                {
                    from: "node_modules/@ministryofjustice/frontend/moj/assets/images",
                    to: path.resolve(__dirname, "web/static/assets/images"),
                },
                {
                    from: "node_modules/@fortawesome/fontawesome-free/webfonts",
                    to: path.resolve(__dirname, "web/static/assets/fonts"),
                },
            ],
        }),
        new MiniCssExtractPlugin({
            filename: "stylesheets/[name].css",
        }),
    ],
};
