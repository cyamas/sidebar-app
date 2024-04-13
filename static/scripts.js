var socket;

function upgradeToWS() {
    socket = new WebSocket("ws://" + document.location.host + "/ws")

    socket.addEventListener('message', function(event) {
        data = JSON.parse(event.data)
        if (data["MsgType"] === "connected") {
            userID = document.getElementById("user-id").value
            msg = {
                "msg-type": "connected",
                "user-id": parseInt(userID)
            }
            socket.send(JSON.stringify(msg))
        }
        if (data["MsgType"] === "newroom") {
            var msg = data["HostName"] + " has created a new chatroom."
           var initialChatMsg = createMessageForSocket(data["RoomID"], data["HostID"], msg)
           console.log(initialChatMsg)
           socket.send(initialChatMsg)
        }
    })

}

document.addEventListener('htmx:afterRequest', upgradeToWS)

function createRoom() {
    userID = document.getElementById('user-id').value
    data = {
        "msg-type": "createroom",
        "user-id": parseInt(userID)
    }
    socket.send(JSON.stringify(data))
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
            
            messageInput.value = ""
            addMessageToFeed(text, roomID)
        }
    }
}

function addMessageToFeed(message, roomID) {
    let feed = document.getElementById("feed-" + roomID)
    let newMessage = document.createElement("p")
    newMessage.id = "my-message"
    newMessage.innerHTML = message
    feed.append(newMessage)
}

function createMessageForSocket(roomID, userID, message) {
    data = {
        "msg-type": "text",
        "room-id": roomID,
        "user-id": userID,
        "message": message,
    }
    return JSON.stringify(data)
}