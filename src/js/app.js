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

define('app', ['ractive', 'logger', 'underscore', 'moment'], function(ractive, logger, _, moment) {
    var ws;

    var apiDateFormat = 'YYYY-MM-DDTHH:mm:ss.SSSSSSZ';

    function dateFormatStd(date) {
        return moment(date, apiDateFormat).format('MM/DD/YYYY HH:mm:ss Z');
    }

    function setLogsDate(logs) {
	logs.forEach(function(l) {
	    l.Date = dateFormatStd(l.Time);
	});
    }

    function startWs(loc) {
	ws = new WebSocket(loc);
	ws.onopen = function(event) {
	    logger.debug("WebSocket :: open");
	};
	ws.onmessage = function(event) {
	    var wsMessage=JSON.parse(event.data);
	    processMessage(wsMessage);
	};
	ws.onclose= function(event) {
	    map = {'Status':'disconnected','Disabled':true,'Message':"Disconnected, wait for reconnexion... or try reloading the page."};
	    ractive.set(map);
	    logger.debug("WebSocket :: closed");
	    ws = null;
	    // Try to reconnect in 5 seconds
	    setTimeout(function(){startWs(loc);}, 5000);
	};
	ws.onerror = function(event) {
	    logger.error("WebSocket :: error");
	    console.error("WebSocket error :", event);
	};
    }

    logger.debug("Defining Module :: app");

    function Init() {
        logger.debug("App :: Init");

	// Open WS
	startWs("ws://" + location.host + "/ws");

	// Connect events
	ractive.on({
	    'start': function( ctx ) {
		ractive.set('Disabled', true);
		cfg = ractive.get('ConfigId');
		msg={Cmd:'start',value:{ConfigId:cfg}};
		ws.send(JSON.stringify(msg));
	    }
	});

    }

    function processMessage(msg) {
	switch (msg.Cmd) {
	    case 'set':
		logger.debug("Msg :: set");
		setLogsDate(msg.Value.Logs);
		ractive.set(msg.Value);
		break;
	    case 'update':
		logger.debug("Msg :: update "+msg.Field);
		switch(msg.Field) {
		    case "Logs":
			setLogsDate(msg.Value);
			break;
		    default:
		}
		ractive.set(msg.Field, msg.Value);
		break;
	    default:
		logger.error("Msg :: Unknown command");
		console.log(msg);
	}
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


