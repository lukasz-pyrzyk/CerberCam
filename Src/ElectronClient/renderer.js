const electron = require('electron')

const ipc = electron.ipcRenderer;

document.getElementById("send").addEventListener('click', _ => {
    var email = document.getElementById("email").value
    console.log("Clicked \"start\" button, sending data through the ipc")

    var buffer = new Buffer("lorem ipsum") // TODO: Jakub, get bytes from user's photo
    var msg = {
        email: email,
        photo: buffer
    }

    ipc.send('newRequest', msg);
})