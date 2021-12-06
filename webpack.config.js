const path = require('path');

module.exports = {
    mode: 'development',
    entry: './web/app.js',
    output: {
        path: path.join(__dirname, 'www'),
        filename: 'bundle.js'
    },
    module: {
        rules: [
            {
                test: /\.m?js$/,
                exclude: /node_modules/,
                use: {
                  loader: 'babel-loader',
                }
            },
            {
                test: /\.css$/i,
                loader: "css-loader"
            }
        ]
    },
    devtool: 'eval-cheap-module-source-map',
    devServer: {
        contentBase: path.join(__dirname, 'www')
    }
}