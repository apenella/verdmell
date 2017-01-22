//
//-- Ready --
//
// This function is run as soon the document is ready
$(document).ready(function() {

	var listenaddr = $(".datacontainer").attr('listenaddr');
	var proto = "http://";
	var baseurl = proto+listenaddr;
	var clusterurl = baseurl + "/api/cluster/";

	$.getJSON(clusterurl, function(data){
		// menuController.initialize(data);
		// clusterlistController.initialize(data);
		MainController.initialize(data);
	});

	//
	// SSE
	//
	if(typeof(EventSource) !== "undefined") {
		var source = new EventSource('/sse');
		source.onopen = function (event) {
			console.log("eventsource connection open");
		};
		source.onerror = function (event) {
			if (event.target.readyState === 0) {
				console.log("reconnecting to eventsource");
			} else {
				console.log("eventsource error");
			}
		};
		source.onmessage = function(event) {
			//clusterlistController.update($.parseJSON(event.data));
			MainController.update(event.data);
		};	
	} else {
		console.log("SSE not supported");
	}

});

var MainController = new Controller({
  
	nodeuri: "/api/node",

	initialize: function(data) {
		//console.log('MainController::initialize', data);
		menuController.initialize(data);
		clusterlistController.initialize(data);
		detailsController.initialize(data);
		locatorController.initialize(menuModel, clusterlistModel);
		checksController.initialize(data);
	},

	update: function(data) {
		//console.log('MainController::update', data);
		clusterlistController.update($.parseJSON(data));
		detailsController.update($.parseJSON(data));
		checksController.update($.parseJSON(data));
	}
});