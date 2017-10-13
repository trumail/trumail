$('.ui.dropdown').dropdown();

$(document).ready(function() {
    $('#testForm').on('submit', function(e) {
        e.preventDefault();
        var format = document.getElementsByName('test-format')[0].value;
        var email = document.getElementsByName('test-email')[0].value;

        // Verify the parameters were passed
        if (format === '' || email === '') {
            return;
        }

        var xhr = new XMLHttpRequest();
        xhr.open('GET', 'https://trumail.io/' + format + '/' + email, true);
        xhr.onload = function(e) {
            var results;
            if (format === 'json') {
                results = vkbeautify.json(xhr.responseText, 4);
            }
            if (format === 'xml') {
                results = vkbeautify.xml(xhr.responseText, 4);
            }
            document.getElementsByName('test-results')[0].textContent = results;
            $('.ui.modal').modal({
                closable: false,
                transition: 'flip vertical'
            }).modal('show');
        };
        xhr.onerror = function(e) {
            console.error(xhr.statusText);
        };
        xhr.send(null);
    });
});