const electron = require('electron')
const {dialog} = require('electron').remote
const fs = require('fs');

const ipc = electron.ipcRenderer;
var fileName
var lastMessage

function loadPage(page) {
    $("#main-content").load(page)
}

loadPage("welcome.html")

$("#close-button").on('click', function() {
    var window = electron.remote.getCurrentWindow();
    window.close();
})