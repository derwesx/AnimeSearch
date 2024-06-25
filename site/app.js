var createError = require('http-errors');
var express = require('express');
var favicon = require('serve-favicon');
var logger = require('morgan');
var path = require('path');

// Routers
var mainRouter = require('./routes/main');
var favouriteRouter = require('./routes/favourite');

// Initialize Express app
var app = express();

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

// Catch 404 and forward to error handler
app.use((req, res, next) => {
    next(createError(404));
});


// Error handler
app.use((err, req, res, next) => {
    // Set locals, only providing error in development
    res.locals.message = err.message;
    res.locals.error = req.app.get('env') === 'development' ? err : {};

    console.error("Unexpected error occurred:", err);

    // Render the error page
    res.status(err.status || 500);
    res.render('error', {errorCode: err.status || 500});
});

module.exports = app;
