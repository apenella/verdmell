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
		menuController.initialize(data);
		clusterlistController.initialize(data);
	});
})