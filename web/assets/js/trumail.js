$('.ui.dropdown').dropdown();

$(document).ready(function() {
    $('#test-form').on('submit', function(e) {
        e.preventDefault();
        var format = document.getElementsByName('test-format')[0].value;
        var email = document.getElementsByName('test-email')[0].value;

        // Verify the parameters were passed
        if (format === '' || email === '') {
            return;
        }

        $('#test-button').addClass('loading');
        var xhr = new XMLHttpRequest();
        var url = '/' + format + '/' + email;
        if (format === 'jsonp') {
            url = url + '?callback=myCallback';
        }
        xhr.open('GET', url, true);
        xhr.onload = function(e) {
            var results = xhr.responseText;
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
            $('#test-button').removeClass('loading');
        };
        xhr.onerror = function(e) {
            console.error(xhr.statusText);
        };
        xhr.send(null);
    });
});