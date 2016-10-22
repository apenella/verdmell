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
var locatorModel = new Model('locatorModel',{
	_location: "",
	
	// set location
	set: function(location) {
		this._location = location;
		locatorModel.set(this);
	},

	// get location
	get: function() {
		return this._location;
	},

	// subscribe to model
	observe: function(model) {
		// subscribe
		locatorModel.on(model.id, locatorModel.id, function(){
			selectedItem = model.attributes.getSelectedItem();
			if (selectedItem.length > 0) {
				locatorModel.attributes.set(selectedItem[0].locator);				
			}
		}.bind(this));
	}

});

//
// View Object
//-----------------------------------------------------------
var locatorView = new View({

	// place to append the items
	parent: ".locator",

	// draw model to view
	render: function(model){
		$('.locator').text(model.attributes.get());
	},

	// define the actions to do when locator changes
	observe: function(model) {
		// subscribe
		this.on(model.id, this.id, function(){
			locatorView.render(model);
		}.bind(this));
	}
});

//
// Controller Object
//-----------------------------------------------------------
var locatorController = new Controller({
	//console.log('locatorController::initialize');

	// menu data management
	model: locatorModel,
	// menu view management
	view: locatorView,

	//
	// locator initializes when menu is ready
	initialize: function(menuModel,clusterlistModel) {
		//
		// set listeners
		//
		// own model: new elements
		this.view.observe(this.model);

		// observe changes to menu or cluster list
		this.model.attributes.observe(menuModel);
		this.model.attributes.observe(clusterlistModel);

		// set the select menu item
		this.model.attributes.set(menuModel.attributes.getSelectedItem()[0].locator);
	}


});