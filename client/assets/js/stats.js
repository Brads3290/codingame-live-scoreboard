let colorAssignments = {};

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
    let stat = query['stat'];

    if (!eventId) {
        throw Error('No event id');
    }

    $(function () {
        // Update the stats container color if we have a stats_color in query string
        if (query['stats_color']) {
            $('.stats-container').css('color', query['stats_color'])
        }

        // Get initial stats
        getStats(stat, eventId, function (data) {
            displayStats(stat, data);
        });

        // Start a timer to update stats
        let interval = window.setInterval(function () {
            getStats(stat, eventId, function (data) {
                displayStats(stat, data);
            });
        }, 7500);
    });


}());

function getStats(stat, eventId, callback) {
    $.ajax({
        method: 'GET',
        url: '/api/stats/' + eventId + '?fetch=' + stat,
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

function displayStats(stat, data) {
    switch (stat) {
        case "round":
            displayRoundsStat(data)
            break
        case "language":
            displayLanguageStat(data)
            break
        default:
            displayError('Invalid stat: ' + stat)
            break
    }
}

function displayRoundsStat(data) {
    $('.stats-container').html('Rounds: ' + data['rounds']['number_of_rounds'])
}

function displayLanguageStat(data) {
    let languageData = data['language'];

    // Clear stats container
    let $stats = $('.stats-container');
    $stats.html('');

    let popularityScores = languageData['popularity'];
    if (!popularityScores || !popularityScores.length) {

        // Display no data
        $stats.html('No data yet.')

        return;
    }

    // Sort the scores, highest score at the top
    popularityScores.sort(function (a, b) {
        return b['rank'] - a['rank'];
    });

    // Take the top 8
    popularityScores = popularityScores.slice(0, 8)

    // get the top score
    let topScore = popularityScores[0].rank;

    let maxWidthPercent = 75;

    let languages = [];
    for (let i = 0; i < popularityScores.length; i++) {
        let entry = $('.stats-templates > .stats-entry').clone();
        $stats.append(entry);

        let thisScore = popularityScores[i]['rank'];

        let scorePercentage;
        if (topScore === 0) {
            scorePercentage = 0
        } else {
            scorePercentage = thisScore / topScore;
        }

        let userColor = getStatBarColor(popularityScores[i]['name']);

        languages.push({
            name: popularityScores[i]['name'],
            width: maxWidthPercent * scorePercentage,
            score: thisScore,
            uses: popularityScores[i]['uses'],
            userColor: userColor,
            entry: entry,
        });
    }

    // Chop off the start of the leaderboard bars
    let redundantPart = languages[languages.length - 1].width;
    if (redundantPart >= 2) {
        redundantPart -= 2;
    }

    for (let i = 0; i < languages.length; i++) {
        languages[i].entry.find('.stats-entry-languagename').html(languages[i]['name'] + " (" + languages[i]['uses'].toString() + ")").css({
            'color': languages[i].userColor
        });

        let w = languages[i].width - redundantPart;

        // Set a minimum width
        if (w < 2) {
            w = 2;
        }

        languages[i].entry.find('.stats-entry-progressbar').css({
            width: w + '%',
            'background-color': languages[i].userColor
        });
    }
}

function displayError(msg) {
    $('.stats-container').html('<span class="error-display">' + msg + '</span>');
}

function getStatBarColor(name) {
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