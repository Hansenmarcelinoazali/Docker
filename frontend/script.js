document.addEventListener("DOMContentLoaded", () => {
    fetch("http://localhost:8080/messages")
        .then(response => response.json())
        .then(data => {
            const messageList = document.getElementById("messageList");
            data.forEach(msg => {
                const li = document.createElement("li");
                li.textContent = `${msg.id}: ${msg.text}`;
                messageList.appendChild(li);
            });
        })
        .catch(error => console.error("Error fetching messages:", error));
});
