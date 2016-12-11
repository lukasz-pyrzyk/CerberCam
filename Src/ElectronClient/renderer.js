const electron = require('electron')
const {dialog} = require('electron').remote
const fs = require('fs');

const ipc = electron.ipcRenderer;
var fileName
var lastMessage

function loadPage(page) {
    $("#main-content").load(page)
}

document.getElementById('selectFile').addEventListener('click', _ => {
    dialogOptions = {
        filters: [
            {
                name: 'Images', extensions: ['jpg', 'png', 'gif']
            }
        ],
        properties: [
            'openFile'
        ]
    }
    dialog.showOpenDialog(dialogOptions, function(data){
        fileName = data[0]
        console.log(fileName)
    })

})

document.getElementById('send').addEventListener('click', _ => {
    var email = electron.remote.getGlobal('user').email;
    console.log('Clicked \'start\' button, sending data through the ipc')

    var buffer = fs.readFileSync(fileName)
    var msg = {
        Email: email,
        Photo: buffer,
        Filename: fileName
    }

    lastMessage = msg

    ipc.send('newRequest', msg);
})

document.getElementById('logout').addEventListener('click', _ => {
    console.log('Logging out...')
    electron.remote.getGlobal('user').email = '';
    loadPage("welcome.html")
})

ipc.on("sendingFinished", (handler, data) => {
    console.log("Request has been send")
    if(lastMessage != null)
        new Notification("Image has been sent!", { body: "to email : " + lastMessage.Email, icon : lastMessage.Filename });
})

document.ondragover = document.ondrop = (ev) => {
  ev.preventDefault()
}

document.body.ondrop = (ev) => {
  console.log(ev.dataTransfer.files[0].path)
  fileName = ev.dataTransfer.files[0].path
  ev.preventDefault()
}
