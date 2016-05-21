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
	_items:[],

	add: function(i, s) {
		this._items.push({name:i, selected: s});			
		menuModel.set(this._items);
	},
	getSelectedItem: function() {
		return _.where(this._items,{selected: true});
	},
	setSelectedItem: function(n) {
		_.each(this._items, function(item){
			if ( item.name == n){
				item.selected = true;	
			} else {
				if (item.selected) {
					item.selected = false;
				}
			}
			menuModel.set(item);
		});
	}

});

//
// View Object
//-----------------------------------------------------------
var menuView = new View({
	parent: ".menu",

	render: function(model){
		//$('.menu').empty();
		_.each(model, function(item){
			if( $('.menuitem#'+item.name).length ) {
				$('.menuitem#'+item.name).attr('name',item.name);
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

		if( $(menuView.parent+' #cleaner').length == 0) {
			$('<div/>',{
				id: 'cleaner',
				style: 'clear:both;'
			}).appendTo(menuView.parent);
		}	
	}.bind(this),

	observe: function(model){
		this.on("menu", this.render(model));
	}

});

//
// Controller Object
//-----------------------------------------------------------
var menuController = new Controller({
	model: menuModel,
	view: menuView,

	events: {
		".menuitem::click": "selectItem"
	},

	initialize: function(data) {
		var selected = true;
		_.each(data, function(content, item){
			menuController._initializeWorker(item,selected);
			if (selected) selected = false;
		});

		this.view.observe(this.model.attributes._items);
		menuController._setEvents();
	},

	_initializeWorker: function(item, selected) {
		this.model.attributes.add(item, selected);
	},
	_setEvents: function() {
		var parts, selector, eventType;
		if(this.events){
			_.each(this.events, function(method, eventName){
				parts = eventName.split('::');
				selector = parts[0];
				eventType = parts[1];
				$(selector).on(eventType, this[method]);
			}.bind(this));
		}
	},

	selectItem: function() {
		console.log(menuController.model);
		menuController.model.attributes.setSelectedItem(this['id']);
		menuController.view.observe(menuController.model.attributes._items);
	}
});



