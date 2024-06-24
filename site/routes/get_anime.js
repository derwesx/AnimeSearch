const express = require('express');
const router = express.Router();

/* Get anime by key. */
router.get('/:currentHash', async (req, res, next) => {
    const currentHash = req.params.currentHash;
    try {
        console.log("Trying to fetch info");
        fetch(`http://localhost:8080/api/anime/${currentHash}`)
            .then((response) => response.json())
            .then(data => {
                console.dir(data["current_hash"]);
            });
    } catch (error) {
        console.error('Error:', error.message);
    }
});

module.exports = router;
