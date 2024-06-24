const express = require('express');
const router = express.Router();

/* GET main page. */
router.get('/', function (req, res, next) {
    res.render('main');
});

module.exports = router;
