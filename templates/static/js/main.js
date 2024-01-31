docUrl = window.location.origin + "/ordinals";

function createInscription() {
    document.getElementById('inscription_form').addEventListener('submit', function(event) {
        event.preventDefault();
    });
    hash = document.getElementById('hash').value;
    output = document.getElementById('output').value;
    key = document.getElementById('key').value;
    content = document.getElementById('content').value;
    address = document.getElementById('address').textContent;
    data = {
        hash: hash,
        output: parseInt(output, 10),
        key: key,
        content: content,
        address: address
    };
    result = $.ajax({
        url: docUrl + '/create_inscription',
        type: 'POST',
        data: JSON.stringify(data),
        dataType: 'text',
        async: true,
        success: function (data) {
            status_bar = document.getElementById('creation_result')
            status_bar.style.color = 'green'
            status_bar.innerHTML = '<b style="color: grey">Status: </b>Inscription was spread'
        },
        error: function (xhr, status, error) {
            status_bar = document.getElementById('creation_result')
            status_bar.style.color = 'red'
            status_bar.innerHTML = '<b style="color: grey">Status: </b>Creation failed'
        }
    });
    cleanForm();
}

function cleanForm() {
    document.getElementById('hash').value = '';
    document.getElementById('output').value = '';
    document.getElementById('key').value = '';
    document.getElementById('content').value = '';
}
