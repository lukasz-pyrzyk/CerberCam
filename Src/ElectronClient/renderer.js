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

$("#capture").on('click', function() {
  alert('Capture button was clicked!');
  var vid = document.getElementById("camera");
  console.log(vid.videoWidth)
  var array = [];
  var canvas = document.getElementById('canvas');
  canvas.width = vid.videoWidth
  canvas.height = vid.videoHeight
  var ctx = canvas.getContext('2d');
  vid.pause();
  ctx.drawImage(vid, 0, 0, vid.videoHeight, vid.videoHeight);

  var imageData = ctx.getImageData(0, 0, vid.videoHeight, vid.videoWidth);
  var arrayBuffer = imageData.data.buffer;  // ArrayBuffer
  var byteArray = new Uint8Array(arrayBuffer)

  var msg = prepareMessage(byteArray, "capture.jpg");
  
  console.log("Sending captured frame through the ipc");
  ipc.send('newRequest', msg);
  vid.play();
})

document.getElementById('send').addEventListener('click', _ => {
    console.log('Clicked \'start\' button, sending data through the ipc')

    var buffer = fs.readFileSync(fileName);    
    lastMessage = prepareMessage(buffer, fileName);

    ipc.send('newRequest', msg);
})

function prepareMessage(buffer, fileName) {
    var email = electron.remote.getGlobal('user').email;

    var msg = {
        Email: email,
        Photo: buffer,
        Filename: fileName
    }

    return msg;
}

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
