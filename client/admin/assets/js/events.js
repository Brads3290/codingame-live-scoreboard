(function () {

    // When the document is ready, create the events datatable
    let $tblEvents;
    $(function () {

        $tblEvents = $('#tblEvents').DataTable({
            ajax: function (data, callback, settings) {
                $.ajax({
                    method: 'GET',
                    url: HTTP_BASE + '/event',
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
                    width: "25%",
                    data: "event_id",
                    render: function (data, type, row) {
                        return '<a href="event.html?event_id=' + data + '">' + data + "</a>"
                    }
                },
                {
                    width: "35%",
                    data: "name"
                },
                {
                    width: "20%",
                    data: "last_updated"
                },
                {
                    width: "20%",
                    data: "event_id",
                    render: function (data, type, row) {
                        return '<button class="btn-delete" data-event-id="' + data + '" data-event-name="' + row['name'] + '">Delete</button>'
                    }
                }
            ],
            rowCallback: function (row, data, displayNum, displayIndex, dataIndex) {
                let $btnDelete = $(row).find('.btn-delete');
                $btnDelete.on('click', function () {
                    if (window.confirm("Are you sure you want to delete this?\n" + $btnDelete.data('eventName'))) {
                        $.ajax({
                            method: 'DELETE',
                            url: HTTP_BASE + '/event/' + $btnDelete.data('eventId'),
                            success: function (res) {
                                if (!res.success) {
                                    alert('Failure: ' + res.error)
                                    console.error(res)
                                    return
                                }

                                $tblEvents.ajax.reload().draw();
                            },
                            error: function (e) {
                                alert('Fatal error: ' + e)
                                console.error(e)
                            }
                        });
                    }
                });
            }
        })
    });

    // Buttons
    let $btnNewEvent;
    $(function () {
        $btnNewEvent = $('#btnNewEvent');
        $btnNewEvent.on('click', function () {
            let eventName = window.prompt('What is the event called?', '');
            if (!eventName) {
                return;
            }

            $.ajax({
                url: HTTP_BASE + '/event',
                method: 'PUT',
                contentType: 'application/json',
                data: JSON.stringify({
                    name: eventName
                }),
                success: function (res) {
                    if (!res.success) {
                        alert('Failed.');
                        console.error(res);
                        return;
                    }

                    $tblEvents.ajax.reload().draw();
                },
                error: function (e) {
                    alert('Really failed.');
                    console.error(e);
                }
            })
        });
    })
}())