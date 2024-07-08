function getRandomString() {
    const num = Math.floor(Math.random() * 100000000);
    return `${num}`
}

let videoContainer, videoPlayer, videoSource;
let miscHandler, descHandler, rateHandler;
let englishNameHandler, japaneseNameHandler;

let currentHash;

const SetAnime = function (data, animeID) {
    // Just in case
    animeID = animeID % data["episodes"].length;
    // Parse Data from Map
    let videoURL = "http://localhost:8080" + data["episodes"][animeID]
    // Put video
    videoPlayer.pause()
    videoSource.setAttribute('src', videoURL);
    videoPlayer.load()
    videoPlayer.play()
}

const SetAnimeBlock = function (data) {
    // Get name, description and rating
    let englishName = data["name"]
    let japaneseName = data["origin_name"]
    englishNameHandler.textContent = englishName;
    // japaneseNameHandler.textContent = japaneseName;
    let description = data["description"]
    let rating = data["rating"]
    descHandler.textContent = description
    rateHandler.textContent = "⭐".repeat(rating)

    // Set current video params
    SetAnime(data, Math.floor(Math.random() * 10))
}

// Open-up handler

let isDescOpen = false;
let leftButton, rightButton;
let expandMask;
let detachedVideoHoverTag;

function OpenDescription() {
    leftButton.style.transform = "translate(-100%)";
    rightButton.style.transform = "translate(+100%)";
    expandMask.style.display = "none";

    videoPlayer.style.transform = "translate(-20%)";
    videoContainer.style.width = "90%";
    videoContainer.classList.remove("no-hover");
    videoContainer.style.visibility = "hidden";
    miscHandler.style.visibility = "visible";
    miscHandler.style.transform = "translate(200%)";
    miscHandler.style.fontSize = "1.4vw";
}

function CloseDescription() {
    leftButton.style.transform = "translate(0)";
    rightButton.style.transform = "translate(0)";
    expandMask.style.display = "flex";

    videoPlayer.style.transform = "translate(0)";
    videoContainer.style.width = "80%";
    videoContainer.classList.add("no-hover");
    videoContainer.style.visibility = "visible";
    miscHandler.style.visibility = "hidden";
    miscHandler.style.transform = "translate(0)";
    miscHandler.style.fontSize = "0";
}

// Next/Previous anime Handler

function NextAnime() {
    fetch(`http://localhost:8080/api/anime/next/${currentHash}`)
        .then(response => response.json())
        .then(data => {
            currentHash = data["current_hash"];
            localStorage.setItem("current_hash", currentHash);
            SetAnimeBlock(data)
        })
        .catch(error => {
            console.error('Error fetching next anime video URL:', error);
        });
}

function PreviousAnime() {
    fetch(`http://localhost:8080/api/anime/prev/${currentHash}`)
        .then(response => response.json())
        .then(data => {
            currentHash = data["current_hash"];
            localStorage.setItem("current_hash", currentHash);
            SetAnimeBlock(data)
        })
        .catch(error => {
            console.error('Error fetching next anime video URL:', error);
        });
}

// Onload handler

window.addEventListener("load", function () {
    currentHash = localStorage.getItem('current_hash');
    if (currentHash == null) {
        currentHash = getRandomString()
    }
    // Buttons
    leftButton = document.querySelector(".left-button-container")
    rightButton = document.querySelector(".right-button-container")
    expandMask = document.querySelector(".expand-mask")
    // Video
    videoPlayer = document.getElementById('video_player');
    videoSource = document.getElementById('video_source');
    videoContainer = document.querySelector('.video-body--video');
    // Text, etc.
    descHandler = document.getElementById('description');
    rateHandler = document.getElementById('rating');
    miscHandler = document.querySelector(".group-up")
    // Name
    englishNameHandler = document.getElementById('video-name-english');
    japaneseNameHandler = document.getElementById('video-name-japanese');
    // Hover
    videoContainer.classList.add("no-hover");

    try {
        fetch(`http://localhost:8080/api/anime/${currentHash}`)
            .then((response) => response.json())
            .then(data => {
                localStorage.setItem("current_hash", data["current_hash"]);
                SetAnimeBlock(data)
            });
    } catch (error) {
        console.error('Error:', error.message);
    }
});