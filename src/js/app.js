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
    var ws;

    logger.debug("Defining Module :: app");

    function Init() {
        logger.debug("App :: Init");
	// testing ...
	//map = {'Status':'idle','Config':['One','Two','Three']};
	//map = {'Status':'idle','Config':['single']};
	//ractive.set(map);
	//ractive.set('Status', 'idle');

	// Open WS
	ws = new WebSocket("ws://" + location.host + "/ws");
	ws.onmessage = function(event) {
		var wsMessage=JSON.parse(event.data);
		processMessage(wsMessage);
	    };
	ws.onclose= function(event) {
	    ractive.set('Status', 'disconnected');
	    console.log("WebSocket is closed now.");
	};
	ws.onerror = function(event) {
	    console.error("WebSocket error observed:", event);
	};

	// Connect events
	ractive.on({
	    'start': function( ctx ) {
		ractive.set("Disabled", true);
		cfg = ractive.get("ConfigId");
		msg={Cmd:"start",value:{ConfigId:cfg}};
		ws.send(JSON.stringify(msg));
	    }
	});

    }

    function processMessage(msg) {
	switch (msg.Cmd) {
	    case 'set':
		logger.debug("Msg :: set");
		ractive.set(msg.Value);
		break;
	    case 'update':
		logger.debug("Msg :: update "+msg.Field);
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


