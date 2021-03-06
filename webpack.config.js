const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    mode: 'development',
    entry: {
        'sqlite_bridge': './node_modules/go-sqlite-js/js/bridge.js',
        // 'go_http_bridge': './node_modules/go-http-js-libp2p/bridge.js',
        'sql_wasm': './node_modules/sql.js/dist/sql-wasm.wasm',
    },
    devtool: 'inline-source-map',
    devServer: {
        contentBase: './',
        publicPath: '/',
        liveReload: false,
    },
    module: {
        rules: [
            {
                test: /\.wasm$/,
                loader: "file-loader",
                type: "javascript/auto", // https://github.com/webpack/webpack/issues/6725
                options: {
                    name: '[name].[ext]',
                    outputPath: './',
                },
            },
        ],
    },
    output: {
        filename: "bundles/[name].js",
        chunkFilename: "bundles/[name].js",
        path: path.resolve(__dirname, 'dist'),
        libraryTarget: 'var',
        library: '[name]',
    },
    node: {
        fs: 'empty'
    },
};
