let socket = new WebSocket("ws://" + document.location.host + "/ws")

socket.onmessage = (event) => {
    console.log(event.data)
}

function sendMessage(event) {
    if (event.key === "Enter") {
        const messageInput = event.target;
        text = messageInput.value.trim()
        
        if (text !== "" && socket.readyState === WebSocket.OPEN) {
            const roomID = messageInput.parentNode.id.slice(5)
            const userID = document.getElementsByClassName("user")[0].id.slice(5)
            msg = createMessageForSocket(roomID, userID, text)
            socket.send(msg)
            
            addMessageToFeed(text, roomID)
            messageInput.value = ""
        }
    }
}

function createMessageForSocket(roomID, userID, message) {
    data = {
        "type": "message",
        "room-id": roomID,
        "user-id": userID,
        "message": message,
        "timestamp": new Date()
    }
    return JSON.stringify(data)
}

function addMessageToFeed(message, roomID) {
    let feed = document.getElementById("feed-" + roomID)
    let newMessage = document.createElement("p")
    newMessage.id = "my-message"
    newMessage.innerHTML = message
    feed.append(newMessage)
}