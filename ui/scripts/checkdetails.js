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
var checkDetailsModel = new Model('checkDetailsModel',{
	
	_check: {},

	set: function(data) {
		this._check =  {};

		checkDetailsModel.set(this._check);
	},

	//
	//
	// subscriptions
	observe: function(model, f) {
		clusterlistModel.on(model.id, this.id, function(){f(model);}.bind(this));
	},

	// observe check for changes
	observeChecks: function(model) {
		console.log('checkDetailsModel::observeChecks',model);
	}


});

//
// View Object
//-----------------------------------------------------------
var checkDetailsView = new View({

	render:function(model) {
		 console.log('checkDetailsView::render');
	},

	observe: function(model) {
		checksModel.on(model.id, this.id, function(model){
			// console.log('checkDetailsView::observe', model);
			checksView.render(model);
		});
	}

});

//
// Controller Object
//-----------------------------------------------------------
var checkDetailsController = new Controller({

	// menu data management
	model: checkDetailsModel,
	// menu view management
	view: checkDetailsView,

	events: {
		".close::click": "close"
	},

	initialize: function() {

		//
		// set listeners
		//
		// subscribe to detailsModel
		this.model.attributes.observe(checksModel,this.model.attributes.observeChecks);
		// subscribe for changes to checksModel
		this.view.observe(this.model);
	},

});