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
var clusterlistModel = new Model('clusterlistModel',{
	
	_elements: [],

	// set data to clusterlist
	set: function(data) {
		this._elements = [];
		// achive the current locator
		this._locator = locatorModel.attributes.get();
		this._selected = false;

		//
		// iterate all the items from data to generate the clusterlist
		// achive types
		_.each(data, function(contentType, type){
			// achieve item
			_.each(contentType, function(contentItem, item){

				// generate locator for current type-item
				locator = clusterlistModel.attributes.generateLocator(type,item);
				// add object
				clusterlistModel.attributes.add(type, item, locator, this._selected, false, contentItem);
			});
		});

		// if there was any selected clusterlist item, this is setted again
		if ( _.size(this._locator) > 0 && _.findWhere(this._elements, {locator: this._locator}) != undefined ) {
		 	clusterlistModel.attributes.setSelected(this._locator);
		}

		// set the element to show
		clusterlistModel.attributes.show(menuModel.attributes.getSelected()[0].name.toLowerCase());
	},

	// add a new item
	add: function(type, item, locator, selected, show, data) {
		this._elements.push(this.create(type, item, locator, selected, show, data));			
	},

	// create clusterlist item
	create: function(type, item, locator, selected, show, data) {
		return {type: type, name: item, locator: locator, selected: selected, show: show, data: data};
	},

	// select items to be showed
	show: function(type) {
		 //console.log('clusterlistModel::show', 'type', type);
		_.each(this._elements, function(item){
			if (item.type == type) {
				item.show = true;
				//console.log('clusterlistModel::show', item);
			} else {
				// console.log('clusterlistModel::show', 'show false', item);
				item.show = false
			}
		});

		// set the array of elements
		clusterlistModel.set(this._elements);
	},

	// set a selected item
	select: function(locator) {
		// review each item
		_.each(this._elements, function(item){
			// set selected
			if ( item.locator == locator) {
				item.selected = true;	
			} else {
				// if any other item is set as selected, unselect it
				if (item.selected) {
					item.selected = false;
				}
			}
		});
		// set the changes an notify subscribers
		clusterlistModel.set(this._elements);
	},

	// return elements marked as showed
	getShowed: function() {
		return _.where(this._elements, {show: true});
	},

	// getSelected
	getSelected: function(){
		return _.where(this._elements, {selected: true});
	},

	// generateLocator
	generateLocator: function (type, item) {
		return "/"+ type + "/" + item;
	},

	//
	//
	// subscriptions
	observe: function(model, f) {
		clusterlistModel.on(model.id, this.id, function(){f(model);}.bind(this));
	},

	observeMenu: function(model) {
		// model contains menu item.
		// get menu's selected item
		
		//console.log('clusterlistModel::observeMenu', menuModel.attributes.getSelected()[0]);
		clusterlistModel.attributes.show(menuModel.attributes.getSelected()[0].name.toLowerCase())
	},

});

//
// View Object
//-----------------------------------------------------------
var clusterlistView = new View({
	// main class
	parent: '.clusterlist',

	// subscriptions
	observe: function(model, f) {
		clusterlistModel.on(model.id, this.id, function(){f(model);}.bind(this));
	},

	render: function(model) {
		// console.log('clusterlistModel::render',clusterlistModel.attributes.getShowed());

		// clear container for not keeping older items
		$(clusterlistView.parent).empty();
		
		// render items
		_.each(clusterlistModel.attributes.getShowed(), function(item) {
				$('<div/>', {
					class: 'clusterlistitem',
					id: item.name,
					type: item.type,
					status: item.data.status,
					text: item.name
				}).appendTo(clusterlistView.parent);
		});
	
		// set events for clusterlist
		clusterlistController.setEvents();
	}
});

//
// Model Object
//-----------------------------------------------------------
var clusterlistController = new Controller({
	model: clusterlistModel,
	view: clusterlistView,

	events: {
		".clusterlistitem::click": "select"
	},
	// initialize the clusterlistCOntroller
	initialize: function(data) {	
		//
		// set listeners
		//

		//
		// subscribe model to menu
		this.model.attributes.observe(menuController.model,this.model.attributes.observeMenu);
		// subscribe view to model
		this.view.observe(this.model, this.view.render);

		// initialitze data for clusterlist
		// set data to model
		clusterlistModel.attributes.set(data);

		// set events for clusterlist
		clusterlistController.setEvents();
	},

	_initializeWorker: function(type, item, content) {
		this.model.attributes.add(type, item, content);
	},

	update: function(data) {
		clusterlistModel.attributes.set(data);
	},

	setEvents: function() {
		var parts, selector, eventType;
		if(this.events){
			_.each(this.events, function(method, eventName){
				parts = eventName.split('::');
				// get item from events property
				selector = parts[0];
				// get method from events property
				eventType = parts[1];
				// hold the event to item 
				$(selector).on(eventType, this[method]);
			}.bind(this));
		}
	},

	select: function() {
		// console.log('clusterlistController.select',this.getAttribute('type'),this.getAttribute('id'))
		clusterlistModel.attributes.select( clusterlistController.model.attributes.generateLocator(this.getAttribute('type'), this.getAttribute('id')));
	}

});