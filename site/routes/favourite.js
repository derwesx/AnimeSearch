const express = require('express');
const router = express.Router();

/* GET favourite page. */
router.get('/', function (req, res, next) {
    res.render('favourite');
});

module.exports = router;
