const HTTP_BASE_API = 'https://cg-leaderboard-api.codezone.brads3290.com';
const HTTP_BASE = 'https://cg-leaderboard.codezone.brads3290.com';

let globals = {
    // Read a page's GET URL variables and return them as an associative array.
    'getUrlVars': function () {
        var vars = [], hash;
        var hashes = window.location.href.slice(window.location.href.indexOf('?') + 1).split('&');
        for (var i = 0; i < hashes.length; i++) {
            hash = hashes[i].split('=');
            vars.push(hash[0]);
            vars[hash[0]] = hash[1];
        }
        return vars;
    }
}