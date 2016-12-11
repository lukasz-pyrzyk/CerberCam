const electron = require('electron')
const {dialog} = require('electron').remote
const fs = require('fs');

const ipc = electron.ipcRenderer;
var fileName

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
    var email = document.getElementById('email').value
    console.log('Clicked \'start\' button, sending data through the ipc')

    var buffer = fs.readFileSync(fileName)
    var msg = {
        Email: email,
        Photo: buffer
    }

    ipc.send('newRequest', msg);
})

ipc.on("sendingFinished", (handler, data) => {
    console.log("Request has been send")
})