var socket;

function upgradeToWS() {
  socket = new WebSocket("ws://" + document.location.host + "/ws");
  console.log("Websocket created");
  socket.addEventListener("message", function (event) {
    data = JSON.parse(event.data);
    if (data["MsgType"] === "connected") {
      userID = document.getElementById("user-id").value;
      msg = {
        "msg-type": "connected",
        "user-id": parseInt(userID),
      };
      console.log(msg);
      socket.send(JSON.stringify(msg));
    }
    if (data["MsgType"] === "newroom") {
      renderNewChatroom(data);
      var msg = data["HostName"] + " has created a new chatroom.";
      var initialChatMsg = createMessageForSocket(
        data["RoomID"],
        data["HostID"],
        msg,
      );

      socket.send(initialChatMsg);
    }
    if (data["MsgType"] === "allusers") {
      var allUsers = data["Users"];
      var selectuserForm = document.getElementById("select-users-form");
      var userID = document.getElementById("user-id").value;
      Object.keys(allUsers).forEach((key) => {
        if (key != userID) {
          let option = document.createElement("input");
          option.id = "user-" + key;
          option.type = "checkbox";
          option.name = allUsers[key];
          option.value = key;

          let label = document.createElement("label");
          label.htmlFor = option.id;
          label.innerHTML = allUsers[key];
          label.appendChild(option);
          selectuserForm.appendChild(label);
        }
      });
      selectuserForm.style.display = "flex";
      selectuserForm.style.flexDirection = "column";
    }
    if (data["MsgType"] === "text") {
      if (document.getElementById("room-" + data["RoomID"]) === null) {
        renderNewChatroom(data);
      }
      let msg = document.createElement("p");
      msg.innerHTML = data["SenderName"] + ": " + data["Msg"];
      let feed = document.getElementById("feed-" + data["RoomID"]);
      feed.appendChild(msg);
    }
  });
}

document.addEventListener("htmx:afterRequest", upgradeToWS);

function createRoom() {
  let memberFormData = new FormData(
    document.getElementById("select-users-form"),
  );

  hostID = document.getElementById("user-id").value;
  var memberIDs = [];
  memberIDs.push(parseInt(hostID));
  for (var tup of memberFormData.entries()) {
    let id = tup[1];
    memberIDs.push(parseInt(id));
  }
  console.log("Members of chatroom: ");
  console.log(memberIDs);
  data = {
    "msg-type": "createroom",
    "user-id": parseInt(hostID),
    "member-ids": memberIDs,
  };
  socket.send(JSON.stringify(data));
}

function sendMessage(event) {
  if (event.key === "Enter") {
    const messageInput = event.target;
    text = messageInput.value.trim();

    if (text !== "" && socket.readyState === WebSocket.OPEN) {
      const roomID = messageInput.parentNode.id.slice(5);
      const userID = document.getElementById("user-id").value;
      msg = createMessageForSocket(roomID, userID, text);
      socket.send(msg);

      messageInput.value = "";
      addMessageToFeed(text, roomID);
    }
  }
}

function addMessageToFeed(message, roomID) {
  let feed = document.getElementById("feed-" + roomID);
  let newMessage = document.createElement("p");
  newMessage.id = "my-message";
  newMessage.innerHTML = message;
  feed.append(newMessage);
}

function createMessageForSocket(roomID, userID, message) {
  data = {
    "msg-type": "text",
    "room-id": parseInt(roomID),
    "user-id": parseInt(userID),
    message: message,
  };
  console.log("from createMessageForSocket: ", data);
  return JSON.stringify(data);
}

function displayForm() {
  let form = document.getElementById("host-room-form");
  if (form.style.display === "none") {
    form.style.display = "flex";
  } else {
    form.style.display = "none";
  }
}

function getAllUsers() {
  userID = document.getElementById("user-id").value;
  let activeUserMsg = {
    "msg-type": "activeusers",
    "user-id": parseInt(userID),
  };
  socket.send(JSON.stringify(activeUserMsg));
}

function renderNewChatroom(data) {
  let chatroom = document.createElement("div");
  chatroom.id = "room-" + data["RoomID"];
  chatroom.classList.add("chatroom");

  let messageFeed = document.createElement("div");
  messageFeed.id = "feed-" + data["RoomID"];
  messageFeed.classList.add("message-log");

  let messageInput = document.createElement("input");
  messageInput.type = "text";
  messageInput.classList.add("input-box");
  messageInput.id = "input-" + data["RoomID"];
  messageInput.addEventListener("keypress", sendMessage);

  chatroom.appendChild(messageFeed);
  chatroom.appendChild(messageInput);
  document.body.appendChild(chatroom);
}
