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

})

var MainController = new Controller({
	initialize: function(data) {
		menuController.initialize(data);	
		clusterlistController.initialize(data);
		detailsController.initialize(data);
		locatorController.initialize(menuModel, clusterlistModel);

	},

	update: function(data) {
		clusterlistController.update($.parseJSON(data));
	}
});