const express = require('express');
const router = express.Router();

/* Get next anime. */
router.get('/', function (req, res, next) {
    res.render('main');
});

module.exports = router;
