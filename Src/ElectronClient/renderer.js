const electron = require('electron')

const ipc = electron.ipcRenderer;

document.getElementById('send').addEventListener('click', _ => {
    var email = document.getElementById('email').value
    console.log('Clicked \'start\' button, sending data through the ipc')

    var buffer = new Buffer('lorem ipsum') // TODO: Jakub, get bytes from user's photo
    var msg = {
        Email: email,
        Photo: buffer
    }

    ipc.send('newRequest', msg);
})

ipc.on("sendingFinished", (handler, data) => {
    console.log("Request has been send")
})