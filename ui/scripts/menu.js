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
var menuModel = new Model({
	// menu item
	_items:[],

	// add an item to menu
	add: function(i, s) {
		this._items.push({name:i, selected: s});			
		menuModel.set(this._items);
	},

	// get all items
	getItems: function() {
		return this._items;
	},

	// get the selected item
	getSelectedItem: function() {
		return _.where(this._items,{selected: true});
	},

	// set a selected item
	setSelectedItem: function(n) {
		_.each(this._items, function(item){
			if ( item.name == n){
				item.selected = true;	
			} else {
				if (item.selected) {
					item.selected = false;
				}
			}
		});
		// set the changes an notify subscribers
		menuModel.set(this._items);
	}
});

//
// View Object
//-----------------------------------------------------------
var menuView = new View({
	// place to append the items
	parent: ".menu",

	// draw model to view
	render: function(model){	
		_.each(model, function(item){
			if( $('.menuitem#'+item.name).length ) {
				//$('.menuitem#'+item.name).attr('name',item.name);
				$('.menuitem#'+item.name).attr('selected',item.selected);
			} else {
				$('<div/>', {
					class: "menuitem",
					id: item.name,
					selected: item.selected,
					text: item.name
				}).appendTo(menuView.parent);
			}
		});
	},

	// subscribe to model
	observe: function(model){
		this.on(model.id+'update', function(model){
			menuView.render(model);
		}.bind(this));
	}

});

//
// Controller Object
//-----------------------------------------------------------
var menuController = new Controller({
	// menu data management
	model: menuModel,
	// menu view management
	view: menuView,

	// events to be hold on each menu item
	events: {
		".menuitem::click": "selectItem"
	},
	//
	// menu initializes when document is ready
	initialize: function(data) {
		// select the default menu item
		var selected = true;
		
		//
		// set listeners
		//
		// own model: new elements
		this.view.observe(this.model);

		// invoke worker for each menu item received
		_.each(data, function(content, item){
			menuController._initializeWorker(item,selected);
			if (selected) selected = false;
		});

		// set events to menu
		menuController._setEvents();
	},

	// 
	_initializeWorker: function(item, selected) {
		this.model.attributes.add(item, selected);
	},

	// hold an action to each menu item to be managed
	_setEvents: function() {
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

	// select item
	selectItem: function() {
		menuController.model.attributes.setSelectedItem(this['id']);
	}
});
