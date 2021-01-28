

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