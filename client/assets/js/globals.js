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

let globals = {

    // Retrieves the URL parameters as an associative array.
    'getUrlVars': function () {
        var vars = [], urlParam;
        var urlParamList = window.location.href.slice(window.location.href.indexOf('?') + 1).split('&');
        for (var i = 0; i < urlParamList.length; i++) {
            urlParam = urlParamList[i].split('=');
            vars.push(urlParam[0]);
            vars[urlParam[0]] = urlParam[1];
        }

        return vars;
    }
}