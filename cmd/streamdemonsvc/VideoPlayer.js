function VideoPlayer(id) {
    var el = document.getElementById(id),
        canvas = document.querySelector("video"),
        playPauseBtn = el.querySelector(".playPauseBtn"),
        progressBar = el.querySelector(".progressBar"),
        volumeControl = el.querySelector(".volumeControl");

    var xhr = new XMLHttpRequest();
    xhr.open("GET", "/assets/video/1.mp4.meta", false);
    xhr.send();
    var videoMeta = JSON.parse(xhr.responseText);

    playPauseBtn.addEventListener("click", function () {
        if (canvas.paused) {
            canvas.play();
            playPauseBtn.classList.remove("paused")
        } else {
            canvas.pause();
            playPauseBtn.classList.add("paused")
        }
    });

    volumeControl.addEventListener("input", function (evt) {
        canvas.volume = evt.currentTarget.value / 100;
        document.cookie = "volume=" + evt.currentTarget.value;
    });
    volumeControl.value = getCookie("volume");
    canvas.volume = getCookie("volume") / 100;

    setInterval(function () {
        var loaded = Math.round(canvas.currentTime) / Math.round(videoMeta["duration"]);
        progressBar.style.width = (loaded * 100) + "%";
    }, 1000);

    return {};
}