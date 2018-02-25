var webpack = require('webpack');
var path = require('path');

var BUILD_DIR = path.resolve(__dirname, 'dist/js');
var APP_DIR = path.resolve(__dirname, 'src');

var config = {
  entry: {
      index: APP_DIR + '/public/index.jsx',
      page1: APP_DIR + '/public/page1.jsx',
      page2: APP_DIR + '/private/page2.jsx',
  },
  output: {
    path: BUILD_DIR,
    filename: '[name].bundle.js'
  },
  module : {
     loaders : [
        {
            test : /\.jsx?/,
            include : APP_DIR,
            loader : 'babel-loader'
        }
     ]
  }
};

module.exports = config;

