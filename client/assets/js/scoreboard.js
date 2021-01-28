

const progressBarColors = [
    '#ffff00',
    '#00ff00',
    '#00ffff',
    '#7fffd4',
    '#a9a9a9',
    '#e9967a',
    '#adff2f',
    '#f0e68c',
    '#f08080',
    '#90ee90',
    '#87cefa',
    '#dda0dd',
    '#ff6347',
];

let colorAssignments = {};

(function () {
    let query = globals.getUrlVars()
    let eventId = query['event_id'];

    if (!eventId) {
        throw Error('No event id');
    }

    window.onerror = function () {
        console.error('Unhandled error. Details=', arguments);
        displayError('Scoreboard: unhandled error.');
    }

    // Wait till document is ready
    $(function () {

        // Fetch initial data
        getScoreboardDataInitial(eventId, function (data) {
            displayScoreboard(data);
        });

        // update timer
        let updateTimer = window.setInterval(function () {
            console.log('Fetching scoreboard update..')
            getScoreboardUpdate(eventId, function (data) {
                console.log('Displaying scoreboard update..')
                displayScoreboard(data);
            });
        }, 7500);
    });
}());

function getScoreboardDataInitial(eventId, callback) {
    $.ajax({
        method: 'GET',
        url: '/api/scoreboard/' + eventId,
        success: function (res) {
            if (!res.success) {
                console.log('/api/scoreboard returned an error: ', res);
                displayError('Error getting scoreboard data.')
            }

            callback(res.data);
        },
        error: function (e) {
            console.log('Failed to get scoreboard data: ', e);
            displayError('Scoreboard: Unknown error.');
        }
    })

}

function getScoreboardUpdate(eventId, callback) {
    $.ajax({
        method: 'GET',
        url: '/api/update/' + eventId,
        success: function (res) {
            if (!res.success) {
                console.log('/api/scoreboard returned an error: ', res);
                displayError('Error getting scoreboard data.')
            }

            callback(res.data);
        },
        error: function (e) {
            console.log('Failed to get scoreboard data: ', e);
            displayError('Scoreboard: Unknown error.');
        }
    });
}

function displayScoreboard(scoreData) {
    let scores = scoreData['scores'];
    if (!scores || !scores.length) {
        return;
    }

    scores.sort(function (a, b) {
        return b.score.event_points - a.score.event_points;
    });

    // Take the top 5
    scores = scores.slice(0, 8)

    // get the top score
    let topScore = scores[0].score.event_points

    // Clear scoreboard container
    let $scoreboard = $('.scoreboard-container');
    $scoreboard.html('');

    let maxWidthPercent = 75;

    let players = [];
    for (let i = 0; i < scores.length; i++) {
        let entry = $('.scoreboard-templates > .scoreboard-entry').clone();
        $scoreboard.append(entry);

        let thisScore = scores[i].score.event_points;
        let scorePercentage = thisScore / topScore;

        let userColor = getProgressBarColor(scores[i].player.name);

        players.push({
            name: scores[i].player.name,
            width: maxWidthPercent * scorePercentage,
            score: thisScore,
            userColor: userColor,
            entry: entry,
        });
    }

    // Min width is the length of the last player
    let minWidth = players[players.length - 1].width;
    if (minWidth >= 2) {
        minWidth -= 2;
    }

    for (let i = 0; i < players.length; i++) {
        players[i].entry.find('.scoreboard-entry-playername').html(scores[i].player.name + " (" + players[i].score.toString() + ")").css({
            'color': players[i].userColor
        });

        players[i].entry.find('.scoreboard-entry-progressbar').css({
            width: players[i].width - minWidth + '%',
            'background-color': players[i].userColor
        });
    }
}

function displayError(msg) {
    $('.scoreboard-container').html('<span class="error-display">' + msg + '</span>');
}

function getProgressBarColor(name) {
    if (colorAssignments[name]) {
        return colorAssignments[name];
    }

    let colorAssignmentCounts = {};
    Object.keys(colorAssignments).forEach(function (key) {
        if (!colorAssignmentCounts[colorAssignments[key]]) {
            colorAssignmentCounts[colorAssignments[key]] = 0;
        }

        colorAssignmentCounts[colorAssignments[key]] += 1;
    });

    // Fill colorAssignmentCounts with any missing colors
    for (let i = 0; i < progressBarColors.length; i++) {
        if (colorAssignmentCounts[progressBarColors[i]] === undefined) {
            colorAssignmentCounts[progressBarColors[i]] = 0
        }
    }

    // Get the lowest count
    let lowestCount = Object.keys(colorAssignments).length;
    Object.keys(colorAssignmentCounts).forEach(function (color) {
        if (colorAssignmentCounts[color] < lowestCount) {
            lowestCount = colorAssignmentCounts[color];
        }
    });

    // Get all the possible colors for that count
    let possibleColors = [];
    Object.keys(colorAssignmentCounts).forEach(function (color) {
        if (colorAssignmentCounts[color] === lowestCount) {
            possibleColors.push(color);
        }
    });

    let idx = Math.floor(Math.random() * possibleColors.length);

    colorAssignments[name] = possibleColors[idx];
    return possibleColors[idx];
}