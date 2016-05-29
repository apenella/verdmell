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
		clusterlistModel.set(this._elements);
	}

});

//
// View Object
//-----------------------------------------------------------
var clusterlistView = new View({
	// main class
	clusterList: 'clusterlist',
	// div classes map
	itemTypeClass: {
		"nodes": "nodeslist",
		"services": "serviceslist"
	},

	initialize: function() {
		// create div to place type's items
		_.each( clusterlistView.itemTypeClass, function(c,t){
			$('<div/>', {
				class: c,
				type: t
			}).hide().appendTo("."+clusterlistView.clusterList);
		})		
	},

	showSelected: function() { 
		selected = menuModel.attributes.getSelectedItem()[0].name;
		_.each(clusterlistView.itemTypeClass, function(c,t){
			if (t == selected){
				$('.'+c).show();
			} else {
				$('.'+c).hide();
			}
		})
	},

	observeMenu: function(model) {
		// subscribe
		this.on(model.id+'update', function(model){
			// menuModel.attributes.getSelectedItem()[0].name --> selected menu item
			// console.log(clusterlistView.itemTypeClass[menuModel.attributes.getSelectedItem()[0].name]);
			clusterlistView.showSelected();
		}.bind(this));
	},

	observeClusterlist: function(model) {
		// subscribe
		this.on(model.id+'update', function(model){		
			_.each(model, function(item){
				console.log("."+clusterlistView.itemTypeClass[item.type]);	
				
				if ( $('.'+clusterlistView.itemTypeClass[item.type]+' .clusternode#'+item.name).length ) {
					$('.'+clusterlistView.itemTypeClass[item.type]+' .clusternode#'+item.name).attr('status',item.data.status);
				} else {
					$('<div/>', {
						class: 'clusternode',
						id: item.name,
						type: item.type,
						status: item.data.status,
						text: item.name
					}).appendTo("."+clusterlistView.itemTypeClass[item.type]);
					}
			});
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
		// initialize view
		this.view.initialize();
		
		//
		// set listeners
		//
		// subscribe view to model
		this.view.observeClusterlist(this.model);
		// subscribe view to menu model, to be aware on clicks
		this.view.observeMenu(menuController.model);

		_.each(data, function(contentType, type){
			_.each(contentType, function(contentItem, item){
				clusterlistController._initializeWorker(type, item, contentItem);
			});
		});
		this.view.showSelected();
	},

	_initializeWorker: function(type, item, content) {
		this.model.attributes.add(type, item, content);
	},

	selectListItem: function() {

	}


});
