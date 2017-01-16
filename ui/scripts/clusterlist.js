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
				clusterlistModel.attributes.add(type, item, locator, this._selected, contentItem);
			});
		});

		// set the array of elements
		clusterlistModel.set(this._elements);

		// if there was any selected clusterlist item, this is setted again
		if ( _.size(this._locator) > 0 && _.findWhere(this._elements, {locator: this._locator}) != undefined ) {
		 	clusterlistModel.attributes.setSelectedItem(this._locator);
		}
	},

	// add a new item
	add: function(type, item, locator, selected, data) {
		this._elements.push(this.createItem(type, item, locator, selected, data));			
	},

	// createItem jsoned item
	createItem: function(type, item, locator, selected, data) {
		return {type: type, name: item, locator: locator, selected: selected, data: data};
	},

	// set a selected item
	setSelectedItem: function(locator) {
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
	// getSelectedItem
	getSelectedItem: function(){
		return _.where(this._elements, {selected: true});
	},
	// generateLocator
	generateLocator: function (type, item) {
		return "/"+ type + "/" + item;
	}
});

//
// View Object
//-----------------------------------------------------------
var clusterlistView = new View({
	// main class
	parent: '.clusterlist',
	// div classes map
	itemTypeClass: {
		"nodes": "nodeslist",
		"services": "serviceslist"
	},

	// create the div for each type
	initialize: function() {
		//console.log('clusterlistController::initialize');

		// create div to place type's items
		_.each( clusterlistView.itemTypeClass, function(c,t){
			$('<div/>', {
				class: c,
				type: t
			}).hide().appendTo(clusterlistView.parent);
		})		
	},

	// show which type is selected
	showSelected: function() {
		selected = menuModel.attributes.getSelectedItem()[0].name.toLowerCase();
		// show only the selected type
		_.each(clusterlistView.itemTypeClass, function(c,t){
			if (t == selected){
				$('.'+c).show();
			} else {
				$('.'+c).hide();
			}
		})
	},

	// define the actions to do when menu changes
	observeMenu: function(model) {
			// subscribe
			this.on(model.id, this.id, function(model){
			// menuModel.attributes.getSelectedItem()[0].name --> selected menu item
			clusterlistView.showSelected();
		}.bind(this));
	},

	// define the actions to do when clusterlist changes
	observeClusterlist: function(model) {
		// subscribe
		this.on(model.id, this.id, function(model){

			// clear container for not keeping older items
			_.each(clusterlistView.itemTypeClass, function(type){
				$(type).empty();
			});

			_.each(model, function(item){				
				if ( $('.'+clusterlistView.itemTypeClass[item.type]+' .clusterlistitem#'+item.name).length ) {
					$('.'+clusterlistView.itemTypeClass[item.type]+' .clusterlistitem#'+item.name).attr('status',item.data.status);
 

				} else {
					$('<div/>', {
						class: 'clusterlistitem',
						id: item.name,
						type: item.type,
						status: item.data.status,
						text: item.name
					}).appendTo("."+clusterlistView.itemTypeClass[item.type]);
				}
			});
		}.bind(this));
	
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
		".clusterlistitem::click": "selectListItem"
	},
	// initialize the clusterlistCOntroller
	initialize: function(data) {
		// initialize view
		this.view.initialize();
		
		//
		// set listeners
		//
		// subscribe view to model
		this.view.observeClusterlist(this.model);
		// subscribe view to menu model, to be aware on clicks
		this.view.observeMenu(menuController.model);

		// initialitze data for clusterlist
		clusterlistModel.attributes.set(data);
		this.view.showSelected();

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

	selectListItem: function() {
		clusterlistController.model.attributes.setSelectedItem( clusterlistController.model.attributes.generateLocator(this.getAttribute('type'), this.getAttribute('id')));
	}

});
