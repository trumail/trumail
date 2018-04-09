$('.ui.dropdown').dropdown();

var timeout = 10000;
var oldStats = {daily:{},monthly:{}};

function pollStats() {
    $.getJSON('/stats', function (newStats) {
        // Set all countups
        countup('dayDeliverable', 24, oldStats.daily.deliverable, newStats.daily.deliverable);
        countup('dayUndeliverable', 24, oldStats.daily.undeliverable, newStats.daily.undeliverable);
        countup('daySuccessRate', 24, oldStats.daily.successRate, newStats.daily.successRate, '', '%');
        countup('monthDeliverable', 720, oldStats.monthly.deliverable, newStats.monthly.deliverable);
        countup('monthUndeliverable', 720, oldStats.monthly.undeliverable, newStats.monthly.undeliverable);
        countup('monthSuccessRate', 720, oldStats.monthly.successRate, newStats.monthly.successRate, '', '%');
        oldStats = newStats; // Update the oldstats

        // Perform this action every 10 seconds
        setTimeout(pollStats, timeout); // Poll stats every 10 seconds and re-apply to UI
    });
}

// countup animates the passed id with a new counted up value
function countup(id, hours, from, to, prefix, suffix) {
    // if from isn't set yet, calculate a rough initial value 
    if (from == undefined) {
        var rps = ((to/hours)/60)/60; // Average requests per second
        from = to-(rps*(timeout/1000)); // Subtract 10 seconds worth
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