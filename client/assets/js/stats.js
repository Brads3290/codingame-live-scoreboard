(function () {

    $(function () {
        // Error handling
        window.onerror = function () {
            console.error('Unhandled error. Details=', arguments);
            displayError('Scoreboard: unhandled error.');
        };
    });

    let query = globals.getUrlVars();
    let eventId = query['event_id'];

    if (!eventId) {
        throw Error('No event id');
    }

    $(function () {
        // Update the stats container color if we have a stats_color in query string
        if (query['stats_color']) {
            $('.stats-container').css('color', query['stats_color'])
        }

        // Get initial stats
        getStats(eventId, function (data) {
            displayStats(data);
        });

        // Start a timer to update stats
        let interval = window.setInterval(function () {
            getStats(eventId, function (data) {
                displayStats(data);
            });
        }, 7500);
    });


}());

function getStats(eventId, callback) {
    $.ajax({
        method: 'GET',
        url: '/api/stats/' + eventId,
        success: function (res) {
            if (!res.success) {
                console.log('/api/stats returned an error: ', res);
                displayError('Error getting stat data.');
                return;
            }

            callback(res.data);
        },
        error: function (e) {
            console.log('Failed to get scoreboard data: ', e);
            displayError('Stat: Unknown error.');
        }
    });
}

function displayStats(data) {
    $('.stats-container').html('Rounds: ' + data['number_of_rounds'])
}

function displayError(msg) {
    $('.stats-container').html('<span class="error-display">' + msg + '</span>');
}