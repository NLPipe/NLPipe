// Get UUID from the URL's query string
const uuid = new URLSearchParams(window.location.search).get("uuid");

// Locate the HTML elements to show data in
const dataContainer = document.querySelector("#data");
const requestId = document.querySelector("#request-id");
const status = document.querySelector("#data #status");
const sentiment = document.querySelector("#data #sentiment");
const lastUpdate = document.querySelector("#last-update span");
const spinner = document.querySelectorAll(".spinner");

requestId.textContent = uuid;

// Get data
const interval = setInterval(() => {
    fetch("/api/result/" + uuid)
    .then(response => {
        if (response.status != 200) {
            dataContainer.textContent = response.status + ": " + response.statusText;
            lastUpdate.parentNode.style.display = "none";
            throw new Error();
        }
        return response.json();
    })
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
            stopPolling();
            return;
        }
    })
    .catch(() => {
        stopPolling();
    })
}, 3000);

function stopPolling() {
    clearInterval(interval);
}