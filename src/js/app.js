define('logger', function() {

    function isEnable() {
        return localStorage ? (localStorage.getItem("flashrom.logging") === 'true') : false;
    }

    function log(severity, message) {
        try {
            if (!isEnable() || typeof console === 'undefined') {
                return;
            }
            var text = "[" + new Date().toLocaleString().replace(",", "") + "] " + severity + " | " + message;
            if (severity === 'error') {
                console.error(text);
            } else {
                console.log(text);
            }
        }
        catch(err) {
            if (typeof console !== 'undefined') {
                console.log("Ignored logging error");
            }
        }
    }

    function info(message) {
        log(' INFO', message);
    }

    function debug(message) {
        log('DEBUG', message);
    }

    function error(message) {
        log('ERROR', message);
    }

    return {
        log: log,
        info: info,
        debug: debug,
        error: error,
        isEnable: isEnable,
    };
});

define('app', ['ractive', 'logger', 'underscore'], function(ractive, logger, _) {

    logger.debug("Defining Module :: app");

    function Init() {
        logger.debug("App :: Init");
	ractive.set('Status', 'idle');
    }

    return {
        Init: Init,
    };

});

define('ractive', ['Ractive', 'logger', 'partials'], function(Ractive, logger, partials){

    logger.debug("Defining Module :: ractive");

    Ractive.DEBUG = logger.isEnable();

    return new Ractive({
        el: 'wrapper',
        template: partials.template,
        partials: partials.list,
        data: {}
    });
});

define('partials', [
    'text!../pages/main.html',
], function(index) {

    return {
        template: index,
        list: {}
    };
});

require(['app', 'logger', 'ractive'], function(app, logger, ractive) {
    logger.info("Initialize application");

    app.Init();
});


