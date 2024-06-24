const createError = require('http-errors');
const express = require('express');
const favicon = require('serve-favicon');
const logger = require('morgan');
const path = require('path');
const fs = require('fs');

const { createProxyMiddleware } = require('http-proxy-middleware');

// Routers
const favouriteRouter = require('./routes/favourite');
const mainRouter = require('./routes/main');
const nextAnimeRouter = require('./routes/next_anime');
const prevAnimeRouter = require('./routes/prev_anime');
const getAnimeRouter = require('./routes/get_anime');

// Initialize Express app
const app = express();

// Setup view engine
app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'jade');

// Middleware
app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({extended: false}));
app.use(express.static(path.join(__dirname, 'public')));
app.use(favicon(path.join(__dirname, 'public', 'icons', 'favicon.ico')));

// Routes
app.use('/', mainRouter);
app.use('/favourite', favouriteRouter);
app.use('/anime/get', getAnimeRouter)
app.use('/next/:currentHash', nextAnimeRouter)
app.use('/prev/:currentHash', prevAnimeRouter)

// Catch 404 and forward to error handler
app.use((req, res, next) => {
    next(createError(404));
});


// Error handler
app.use((err, req, res, next) => {
    // Set locals, only providing error in development
    res.locals.message = err.message;
    res.locals.error = req.app.get('env') === 'development' ? err : {};

    console.error("Unexpected error occurred:", err.message);

    // Render the error page
    res.status(err.status || 500);
    res.render('error', {errorCode: err.status || 500});
});

module.exports = app;
