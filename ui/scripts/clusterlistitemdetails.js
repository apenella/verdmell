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
// Model Object
//-----------------------------------------------------------
var detailsModel = new Model('detailsModel',{

});

//
// View Object
//-----------------------------------------------------------
var detailsView = new View({
	// place to append the items
	parent: ".clusterlistitemdetails",

	render: function() {
		console.log('detailsView::render');
		$(detailsView.parent).html(detailsView.parent);
	},

	// subscribe to model
	observe: function(model){
		// subscribe
		this.on(model.id, this.id, function(model){
			menuView.render(model);
		}.bind(this));
	}

});

//
// Controller Object
//-----------------------------------------------------------
var detailsController = new Controller({
	// menu data management
	model: detailsModel,
	// menu view management
	view: detailsView,

	//
	// menu initializes when document is ready
	initialize: function(data) {
		//console.log('detailsController::initialize');

		this.view.render();
	}

});