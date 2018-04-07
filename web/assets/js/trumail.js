$('.ui.dropdown').dropdown();

function pollStats() {
    $.getJSON('/stats', function (data) {
        // Set all countups
        countup('dayDeliverable', data.daily.deliverable);
        countup('dayUndeliverable', data.daily.undeliverable);
        countup('daySuccessRate', data.daily.successRate, '', '%');
        countup('monthDeliverable', data.monthly.deliverable);
        countup('monthUndeliverable', data.monthly.undeliverable);
        countup('monthSuccessRate', data.monthly.successRate, '', '%');

        // Perform this action every 10 seconds
        setTimeout(pollStats, 10000); // Poll stats every 10 seconds and re-apply to UI
    });
}

// countup animates the passed id with a new counted up value
function countup(id, to, prefix, suffix) {
    var from = $('#' + id).text(); // Retrieve the starting value
    from = from.trim(); // Trim any whitespace
    from = from.replace(/,/g, ''); // Remove all commas
    if (from == '' || from > to) {
        from = to;
    }

    // Configure countup options
    var options = { 
        useEasing: false, 
        useGrouping: true, 
        separator: ',', 
        decimal: '.'
    };

    // Apply a prefix if one is passed
    if (prefix != undefined && prefix != '') {
        options.prefix = prefix;
    }

    // Apply a suffix if one is passed
    if (suffix != undefined && suffix != '') {
        options.suffix = suffix;
    }

    // Trigger the countup animation
    var count = new CountUp(id, from, to, 0, 10, options);
    if (!count.error) {
        count.start();
    } else {
        console.error(count.error);
    }
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