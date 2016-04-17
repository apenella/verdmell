/*
	Package 'ui' 
	-server
	-handler
	-router
	-routes

	-html/
	-images/
	-pages/
	-scripts/
	-style/

*/

//
//# Object to handle the verdmell page 
//# page
function page(l,p) {	
	this.listenaddr = l;
	this.proto = p;
	
	var baseurl = this.proto+this.listenaddr;
	var clusterurl = baseurl + "/api/cluster/";
	
	this.cluster = new cluster(clusterurl);
	this.generateObjects();

// end page object
}
//#
//# page prototypes
//
//# generateObjects
page.prototype.generateObjects = function() {
  //getJSON
  var default_item = 'nodes';
  $.getJSON(this.cluster.url,function(data){
		// Create menu object
		createContainerMenu(data);
		// Create node object
		createContainerClusterList(data, default_item);
	// end page::generateObjects getJSON for clusterurl
	});
// end page::generateObjects
};

//
//# createMessageForUser
function  createMessageForUser(parent, type, message){
	$('<div/>', {
		class: 'message',
		type: type,
		text: message 
	}).appendTo(parent);
}

//
//# createContainerMenu
function createContainerMenu(data){
	m = new pageObject(".menu", function(parent){
		$.each(data,function(i){
			createMenuItem(data, parent, i)
		});

		$('<div/>',{
			style: 'clear:both;'
		}).appendTo(parent);
	});
}
//
//# createMenuItem
function createMenuItem(data, parent, item){
	i = new pageObject(parent, function(parent){
		$('<div/>', {
			class: "menuitem",
			id: item,
			text: item
		}).appendTo(parent);		
	});

	i.setOnClickAction('.menuitem#'+item,function(){
		//fill container clusterlist
		createContainerClusterList(data, item);
	});
}
//
//# clearContainerClusterList(){
function  clearContainerClusterList(){
	$('.clusterlist').empty();
}
//
//# createContainerClusterList
function createContainerClusterList(data, item) {
	$('.clusterlist').empty();

	n = new pageObject('.clusterlist', function(parent) {
		// Create a new node for each node defined on cluster
		if (data[item] != null) {
			switch(item){
				case "nodes":
						$.each(data[item], function(nodename, clusternode){
							// Create a node
							createNodeForClusterList(parent,nodename,clusternode);
						});
					break;
				case "services":
					//TODO
					createMessageForUser(parent,"info","Unknown item selected");
					break;
				default:
					//TODO
					createMessageForUser(parent,"info","Unknown item selected");
			}
		} else {
			createMessageForUser(parent,"info","No data available for "+item);
		}
	// end rendeClusterNodes
	});
}
//
//# updateContainerClusterList
function updateContainerClusterList(data, item) {
	
	if (data.data[item] != null ) {
		switch(item) {
			case "nodes":
				$.each(data.data[item], function(nodename, clusternode){
					if ($('.clusternode#'+nodename).attr('status') != null) {
						console.log(nodename+": "+clusternode.services[nodename].status+" - "+clusternode.status);
 						$('.clusternode#'+nodename).attr('status',clusternode.status)
 					}
				});

				break;
		}
	}
}


//
//# createNodeForClusterList
function createNodeForClusterList(parent, nodename, clusternode) {
	node = new pageObject(parent,function(parent){
		nodeurl = clusternode.URL + "/api/node";
		if (clusternode.services[nodename] != null ) {
			$('<div/>', {
				class: 'clusternode',
				id: nodename,
				url: nodeurl,
				status: clusternode.services[nodename].status,
				text: nodename
			}).appendTo(parent);
		} else {
			$('<div/>', {
				class: 'clusternode',
				id: nodename,
				url: nodeurl,
				text: nodename 
			}).appendTo(parent);
		//end check clusternode.service 
		}
	});
	setActionToNode(node,nodename,nodeurl);
}

//
//# setActionToNode()
function setActionToNode(node, nodename, nodeurl) {
	node.setOnClickAction('.clusternode#'+nodename, function(object){

		if ( $('.servicesdetails').attr('id') == nodename ) {
			$('.servicesdetails#'+nodename).toggle(200);
		} else {
			// it generates servicesdetails div
			createContainerServicesDetails('.clusternodedetails',nodename);
			fillContainerServicesDetails('.servicesdetails',nodeurl);
		}
	// end page::generateObjects --> node setOnClickAction
	});
}

//
//# createContainerServicesDetails
function createContainerServicesDetails(parent,nodename) {
	
	/*
	for each service is nested a servicedetail structure into servicedetails
	 ___________________________________________
	| servicesdetails
	|	 _______________________________________
	|	|	servicesdetail									  		|
	|	|  _________________ 						  			|
	|	|	| servicename 		| 				  				|
	|	|	|_________________| 			  					|
	|	|	 ___________________________________	|
	|	|	| servicechecks 				  					|	| 
	|	|	|	 _____________ 	   _____________ 	| |
	|	|	|	|servicecheck |...|servicecheck |	| |
	|	|	|	|_____________|	  |_____________|	| |
	|	|	|___________________________________|	|
	|	|_______________________________________|
	|	...
	|	 _______________________________________
	|	|	servicesdetail									  		|
	|	|  _________________ 						  			|
	|	|	| servicename 		| 				  				|
	|	|	|_________________| 			  					|
	|	|	 ___________________________________	|
	|	|	| servicechecks 				  					|	| 
	|	|	|	 _____________ 	   _____________ 	| |
	|	|	|	|servicecheck |...|servicecheck |	| |
	|	|	|	|_____________|	  |_____________|	| |
	|	|	|___________________________________|	|
	|	|_______________________________________|
	|___________________________________________
	*/

	clearContainerServicesDetails();

	servicesdetails = new pageObject(parent,function(parent){
		// Create the new servicesdatails container
		$('<div/>',{
			class: 'servicesdetails',
			id: nodename
		}).appendTo(parent);
	});
}
//
//#clearContainerServicesDetails
function clearContainerServicesDetails(){
	// Clear servicesdetails container's content before to update it
	$('.servicesdetails').remove();
}

//
//#fillContainerServicesDetails
function fillContainerServicesDetails(parent,url) {
	// Fill th the content
	d = new pageObject(parent, function(parent){
		//get JSON from /api/node
		$.getJSON(url, function(data){

			if (data.services.servicesroot.services != null) {
				// path to get services from /api/node		 				
				$.each(data.services.servicesroot.services, function(servicename, servicedetail){
					//create servicedetail for node
					createContainerServiceDetail(data,'.servicesdetails',servicename,servicedetail);
				// end each data.services
				});
			// if data.services.servicesroot.services != null
			} else {
				// if 
				createMessageForUser(parent,"info","No services available");
			// end if data.services.servicesroot.services != null
			}
		// end page::generateObjects --> getJSON /api/node	
		});
	// end page::generateObjects::pageObject definition clusternodedetails
	});
}
//
//# createContainerServiceDetail
function createContainerServiceDetail(data, parent, servicename, servicedetail) {
	detail = new pageObject(parent,function(parent){
		$('<div/>', {
				class: 'servicedetail',
				status: servicedetail.status,
				id: servicename 
			}).appendTo(parent);
	});
	createContentServiceDetail(data, servicename, servicedetail);
// end createContainerServiceDetail
}
//
//# createContentServiceDetail
function createContentServiceDetail(data, servicename, servicedetail){
	content = new pageObject('.servicedetail', function(parent){
		createContainerServiceName(parent, servicename, servicedetail.status);
		createContainerServiceChecks(data, parent, servicename, servicedetail.checks);
	// end pageObject detail
	});
//end createContentServiceDetail
}
//
//# createContainerServiceChecks
function createContainerServiceName(parent, servicename, status) {
	//
	//create servicename for servicedetail	
	servicename = new pageObject(parent, function(parent){
		$('<div/>', {
			class: 'servicename',
			status: status,
			id: servicename,
			text: servicename 
		}).appendTo(parent+'#'+servicename);
	});
}
//
//# createContainerServiceChecks
function createContainerServiceChecks(data, parent, servicename, checks) {
	//
	//create servicechecks
	servicechecks = new pageObject(parent, function(parent){
		$('<div/>',{
			class: 'servicechecks',
			id: servicename
		}).appendTo(parent+'#'+servicename);
	});

	//				 					
	//create servicecheck for servicechecks
	$.each(checks,function(it,check){

		if ( data.samples.Samples[check] != null ){
			createContainerCheck(data, servicename, check, data.samples.Samples[check].Sample.exit);			
		} else {
			//console.log("No sample for "+check);
			createContainerCheck(data, servicename, check, "-1");
		}

		// end for each servicedetail.checks
	});
		
	//
	// add a clear for relative position
	clear = new pageObject('.servicechecks#'+servicename,function(parent) {
		$('<div/>',{
			style: 'clear:both;'
		}).appendTo(parent);
	// end pageObject clear
	});

}

//
//# createContainerCheck
function createContainerCheck(data, servicename, checkname, status){

	check = new pageObject('.servicechecks#'+servicename, function(parent){
		$('<div/>',{
				class: 'servicecheck',
				status: status,
				id: servicename+'_'+checkname,
				text: checkname 
		}).appendTo(parent);
	// en pageObject check
	});

	check.setOnClickAction('.servicecheck#'+servicename+'_'+checkname, function(parent){
		createContainerCheckAllDetails(data, servicename, checkname);
	});
}
//
//# clearContainerCheckDetails
function clearContainerCheckAllDetails() {
	$('.servicecheckalldetails').remove();
}
//
//# createContainerCheckAllDetails
function createContainerCheckAllDetails(data, servicename, checkname){
	if ( $('.servicecheckalldetails').attr('id') == servicename+'_'+checkname ) {
		// check if it's clicked the same checkname
		$('.servicecheckalldetails#'+servicename+'_'+checkname).toggle(200);
	} else {
		clearContainerCheckAllDetails();

		detail = new pageObject('.servicechecks#'+servicename, function(parent){
			$('<div/>',{
				class: 'servicecheckalldetails',
				id: servicename+'_'+checkname
			}).appendTo(parent);
		});

		createContainerCheckADetails(data.checks.checks.checks[checkname], servicename, checkname);
		createContainerCheckSampleDetails(data.samples.Samples[checkname], servicename, checkname);		
		
		clear = new pageObject('.servicechecks#'+servicename,function(parent) {
			$('<div/>',{
				style: 'clear:both;'
			}).appendTo(parent);
		// end pageObject clear
		});
	}
}
//
//# createContainerCheckADetails
function createContainerCheckADetails(data, servicename, checkname){

	details = new pageObject('.servicecheckalldetails#'+servicename+'_'+checkname, function(parent){
		$('<div/>',{
			class: 'servicecheckdetails',
			id: servicename+"_"+checkname
		}).appendTo(parent);
	});

	if (data != null) {
		detailcontent = new pageObject('.servicecheckdetails#'+servicename+'_'+checkname, function(parent){
			// Name
			$('<div/>',{
				class: 'checkinfotitle',
				id: servicename+"_"+checkname,
				text: 'Check'
			}).appendTo(parent);
			$('<div/>',{
				class: 'checkinfo',
				id: servicename+"_"+checkname,
				text: data.name
			}).appendTo(parent);
			// Description
			$('<div/>',{
				class: 'checkinfotitle',
				id: servicename+"_"+checkname,
				text: 'Description'
			}).appendTo(parent);
			$('<div/>',{
				class: 'checkinfo',
				id: servicename+"_"+checkname,
				text: data.description
			}).appendTo(parent);
			// Command
			$('<div/>',{
				class: 'checkinfotitle',
				id: servicename+"_"+checkname,
				text: 'Command'
			}).appendTo(parent);
			$('<div/>',{
				class: 'checkinfo',
				id: servicename+"_"+checkname,
				text: data.command
			}).appendTo(parent);
		});
	} else {
		createMessageForUser('.servicecheckdetails#'+servicename+'_'+checkname,"info","No information for '"+checkname+"'")
		clear = new pageObject('.servicecheckalldetails#'+servicename+'_'+checkname,function(parent) {
			$('<div/>',{
				style: 'clear:both;'
			}).appendTo(parent);
			// end pageObject clear
		});
	}

	clear = new pageObject('.servicecheckdetails#'+servicename+'_'+checkname,function(parent) {
		$('<div/>',{
			style: 'clear:both;'
		}).appendTo(parent);
	// end pageObject clear
	});
}
//
//# createContainerCheckSampleDetails
function createContainerCheckSampleDetails(data, servicename, checkname){

	details = new pageObject('.servicecheckalldetails#'+servicename+'_'+checkname, function(parent){
		$('<div/>',{
			class: 'servicechecksampledetails',
			id: servicename+"_"+checkname
		}).appendTo(parent);
	});

	if (data != null) {
		detailcontent = new pageObject('.servicechecksampledetails#'+servicename+'_'+checkname, function(parent){
			// Sample date
			$('<div/>',{
				class: 'sampleinfotitle',
				id: servicename+"_"+checkname,
				text: 'Sample time'
			}).appendTo(parent);
			$('<div/>',{
				class: 'sampleinfo',
				id: servicename+"_"+checkname,
				text: data.Sample.sampletime
			}).appendTo(parent);
			// Exit value
			$('<div/>',{
				class: 'sampleinfotitle',
				id: servicename+"_"+checkname,
				text: 'Exit value'
			}).appendTo(parent);
			$('<div/>',{
				class: 'sampleinfo',
				id: servicename+"_"+checkname,
				text: data.Sample.exit
			}).appendTo(parent);
			// Exit value
			$('<div/>',{
				class: 'sampleinfotitle',
				id: servicename+"_"+checkname,
				text: 'Elapsed time (ns)'
			}).appendTo(parent);
			$('<div/>',{
				class: 'sampleinfo',
				id: servicename+"_"+checkname,
				text: data.Sample.elapsedtime
			}).appendTo(parent);
			// Output
			$('<div/>',{
				class: 'sampleinfotitle',
				id: servicename+"_"+checkname,
				text: 'Output'
			}).appendTo(parent);
			$('<div/>',{
				class: 'sampleinfo',
				id: servicename+"_"+checkname,
				text: data.Sample.output
			}).appendTo(parent);
			// Timestamp
			$('<div/>',{
				class: 'sampleinfotitle',
				id: servicename+"_"+checkname,
				text: 'Timestamp'
			}).appendTo(parent);
			$('<div/>',{
				class: 'sampleinfo',
				id: servicename+"_"+checkname,
				text: data.Sample.timestamp
			}).appendTo(parent);
		});
	} else {
		createMessageForUser('.servicecheckdetails#'+servicename+'_'+checkname,"info","No sample for '"+checkname+"'")
		clear = new pageObject('.servicecheckalldetails#'+servicename+'_'+checkname,function(parent) {
			$('<div/>',{
				style: 'clear:both;'
			}).appendTo(parent);
			// end pageObject clear
		});
	}

	clear = new pageObject('.servicechecksampledetails#'+servicename+'_'+checkname,function(parent) {
		$('<div/>',{
			style: 'clear:both;'
		}).appendTo(parent);
	// end pageObject clear
	});

	// clear the servicecheckalldetail once all data have been written into content
	clear = new pageObject('.servicecheckalldetails#'+servicename+'_'+checkname,function(parent) {
		$('<div/>',{
			style: 'clear:both;'
		}).appendTo(parent);
	// end pageObject clear
	});
		
}

//
//#clearNodesDetailsContainer
function clearNodesDetailsContainer(){
	// Clear servicesdetails container's content before to update it
	$('.nodesdetails').remove();
}
//
//# createNodesDetailsContainer
function createNodesDetailsContainer(parent) {
	clearNodesDetailsContainer();
	nodesdetails = new pageObject(parent,function(parent){
		// Create the new servicesdatails container
		$('<div/>',{
			class: 'nodesdetails'
		}).appendTo(parent);
	});
}


//
//
//# Object to handle the verdmell page 
//# page
function pageObject(parent,render) {
	this.renderObject(parent,render);
}
//#
//# menu prototypes
//
//# cleanElement
pageObject.prototype.clearObject = function(item) {
	$(item).removeData();
	// end clearElement
};
//# appendElement
pageObject.prototype.appendObject = function(item, parent) {
	$(parent).appendTo(item)
	// en renderElement
};
//# renderElement
pageObject.prototype.renderObject = function(parent, render) {
	this.clearObject(parent);
	render(parent);
// end renderElement
};
pageObject.prototype.setOnClickAction = function(object, action) {
	$(object).click( function(){
		action($(this));
	});
// end setOnClickAction
};

//
//# Object to handle the elements from cluster 
//# cluster
function cluster(url) {
	this.url = url;
//end cluster object	
}

//
//-- Ready --
//
// This function is run as soon the document is ready
$(document).ready(function(){
	listenaddr = $(".datacontainer").attr('listenaddr');

	//define the ui's page
	var p = new page(listenaddr,"http://");

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
			console.log(event)

			updateContainerClusterList($.parseJSON(event.data),"nodes");
		};	
	} else {
		console.log("SSE not supported");
	} 

});