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
var clusterlistModel = new Model({
	_elements: [],

	add: function(type, item, data) {
		this._elements.push({type: type, name: item, data: data});			
		menuModel.set(this._elements);
	}

});

//
// View Object
//-----------------------------------------------------------
var clusterlistView = new View({

	observe: function(model) {
		console.log(model);

		this.on(this.id+'_update', function(model){
			console.log(model);
		}.bind(this));
	}

});

//
// Model Object
//-----------------------------------------------------------
var clusterlistController = new Controller({
	model: clusterlistModel,
	view: clusterlistView,

	events: {
		".clusternode::click": "selectListItem"
	},

	initialize: function(data) {
		_.each(data, function(contentType, type){
			_.each(contentType, function(contentItem, item){
				clusterlistController._initializeWorker(type, item, contentItem);
			});
		});
		
		//
		// set listeners
		//
		// listener to menu
		this.view.observe(this.model.attributes._elements);


	},

	_initializeWorker: function(type, item, content) {
		this.model.attributes.add(type, item, content);
	},

	selectListItem: function() {

	}


});
