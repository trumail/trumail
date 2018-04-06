$('.ui.dropdown').dropdown();

var deliverable = 0;
var undeliverable = 0;
var successRate = 0;

function pollStats() {
    $.getJSON('/stats', function (data) {
        // Set the initial values
        if (deliverable == 0) {
            deliverable = data.deliverable - 35;
        }
        if (undeliverable == 0) {
            undeliverable = data.undeliverable - 5;
        }
        if (successRate == 0) {
            successRate = data.successRate;
        }

        var options = { useEasing: false, useGrouping: true, separator: ',', decimal: '.'};
        var count = new CountUp('deliverable', deliverable, data.deliverable, 0, 10, options);
        if (!count.error) {
            count.start();
        } else {
            console.error(count.error);
        }
        var count = new CountUp('undeliverable', undeliverable, data.undeliverable, 0, 10, options);
        if (!count.error) {
            count.start();
        } else {
            console.error(count.error);
        }
        options.suffix = '%';
        var count = new CountUp('successRate', successRate, data.successRate, 0, 10, options);
        if (!count.error) {
            count.start();
        } else {
            console.error(count.error);
        }

        // Set the global variables
        deliverable = data.deliverable;
        undeliverable = data.undeliverable;
        successRate = data.successRate;

        // Perform this action every 10 seconds
        setTimeout(pollStats, 10000); // Poll stats every 10 seconds and re-apply to UI
    });
}

$(document).ready(function () {
    pollStats();
    $('#test-form').on('submit', function (e) {
        e.preventDefault();
        var format = document.getElementsByName('test-format')[0].value;
        var email = document.getElementsByName('test-email')[0].value;

        // Verify the parameters were passed
        if (format === '' || email === '') {
            return;
        }

        // Set the loading button
        $('#test-button').addClass('loading');

        // Build the request URL
        var url = '/' + format + '/' + email;
        if (format === 'jsonp') {
            url = url + '?callback=myCallback';
        }

        // Perform the get request
        $.get(url, function (data) {
            if (format === 'json') {
                data = vkbeautify.json(data, 4);
            }
            if (format === 'xml') {
                data = vkbeautify.xml(data, 4);
            }
            document.getElementsByName('test-results')[0].textContent = data;
            $('.ui.modal').modal({
                closable: false,
                transition: 'flip vertical'
            }).modal('show');
            $('#test-button').removeClass('loading');
        }, 'text');
    });
});