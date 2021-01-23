(function () {

    // Get event_id from querystring
    let query = globals.getUrlVars();
    let EventId = query['event_id'];

    // When the document is ready, create the events datatable
    let $tblRounds;
    $(function () {

        $tblRounds = $('#tblRounds').DataTable({
            ajax: function (data, callback, settings) {
                $.ajax({
                    method: 'GET',
                    url: '/round/' + EventId,
                    success: function (res) {
                        if (!res.success) {
                            alert('Failure: ' + res.error)
                            console.error(res)
                            return
                        }

                        callback(res);
                    },
                    error: function (e) {
                        alert('Fatal error: ' + e)
                        console.error(e)
                    }
                })
            },
            autoSize: false,
            order: [[2, 'desc']],
            columns: [
                {
                    width: "60%",
                    data: "round_id",
                    render: function (data, type, row) {
                        return '<a target="_blank" href="https://www.codingame.com/clashofcode/clash/report/' + data + '">' + data + "</a>"
                    }
                },
                {
                    width: "20%",
                    data: "is_active"
                },
                {
                    width: "20%",
                    data: "round_id",
                    render: function (data, type, row) {
                        return '<button class="btn-delete" data-round-id="' + data + '">Remove</button>'
                    }
                }
            ],
            rowCallback: function (row) {
                let $btnDelete = $(row).find('.btn-delete');
                $btnDelete.on('click', function () {
                    if (window.confirm("Are you sure you want to remove this round?\n" + $btnDelete.data('roundId'))) {
                        $.ajax({
                            method: 'DELETE',
                            url: '/round/' + EventId + '/' + $btnDelete.data('roundId'),
                            success: function (res) {
                                if (!res.success) {
                                    alert('Failure: ' + res.error)
                                    console.error(res)
                                    return
                                }

                                $tblRounds.ajax.reload().draw();
                            },
                            error: function (e) {
                                alert('Fatal error: ' + e)
                                console.error(e)
                            }
                        });
                    }
                });
            }
        });

        // Scoreboard link
        let sbLink = '../scoreboard.html?event_id=' + EventId;
        $('#lnkScoreboard').attr('href', sbLink).html('/scoreboard.html?event_id=' + EventId)
    });

    // Buttons
    let btnNewRound;
    let btnUpdateEvent;
    $(function () {
        btnNewRound = $('#btnNewRound');
        btnNewRound.on('click', function () {
            let roundId = window.prompt('Paste the unique ID from the CodinGame URL here', '');
            if (!roundId) {
                return;
            }

            $.ajax({
                url: '/round/' + EventId,
                method: 'PUT',
                contentType: 'application/json',
                data: JSON.stringify({
                    round_id: roundId
                }),
                success: function (res) {
                    if (!res.success) {
                        alert('Failed.');
                        console.error(res);
                        return;
                    }

                    $tblRounds.ajax.reload().draw();
                },
                error: function (e) {
                    alert('Really failed.');
                    console.error(e);
                }
            })
        });

        btnUpdateEvent = $('#btnUpdateEvent');
        btnUpdateEvent.on('click', function () {
            $.ajax({
                url: '/update/' + EventId,
                method: 'GET',
                success: function (res) {
                    if (!res.success) {
                        alert('Failed.');
                        console.error(res);
                        return;
                    }

                    $tblRounds.ajax.reload().draw();
                },
                error: function (e) {
                    alert('Really failed.');
                    console.error(e);
                }
            })
        });
    })
}());