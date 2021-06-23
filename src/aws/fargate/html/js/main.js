// Get UUID from the URL's query string
const uuid = new URLSearchParams(window.location.search).get("uuid");

// Locate the HTML elements to show data in
const requestId = document.querySelector("#request-id");
const status = document.querySelector("#data #status");
const sentiment = document.querySelector("#data #sentiment");
const lastUpdate = document.querySelector("#last-update span");
const spinner = document.querySelectorAll(".spinner");

requestId.textContent = uuid;

// Get data
const interval = setInterval(() => {
    fetch("/api/result/" + uuid)
    .then(response => response.json())
    .then(data => {
        const datetime = new Date();
        status.textContent = data.status;
        sentiment.textContent = data.sentiment;
        lastUpdate.textContent = [
            datetime.toDateString(), " ",
            datetime.getHours(), ":",
            datetime.getMinutes(), ":",
            datetime.getSeconds()]
        .join("");

        spinner.forEach(s => s.style.display = data.status == "processing" ? "inline-block" : "none");

        if (data.status == "finished") {
            clearInterval(interval);
        }
    })
}, 3000);