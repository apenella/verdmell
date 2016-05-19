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
	items:[],

	add: function(i) {
		this.items.push({name:i, selected:false});
		menuModel.set(this.items);
	},
	getSelectedItem: function() {
		console.log('getSelectedItem');
	}

});

//
// View Object
//-----------------------------------------------------------
var menuView = new View({
	parent: ".menu",

	render: function(model){
		_.each(model, function(item){
			$('<div/>', {
				class: "menuitem",
				id: item.name,
				selected: item.selected,
				text: item.name
			}).appendTo(menuView.parent);

		});
		$('<div/>',{
			style: 'clear:both;'
		}).appendTo(menuView.parent);
	},

	observe: function(model){
		this.on("menu",this.render(model));
	}

});

//
// Controller Object
//-----------------------------------------------------------
var menuController = new Controller({
	model: menuModel,
	view: menuView,

	initialize: function(data) {
		//
		_.each(data, function(content, item){
			menuController._initializeWorker(item);
		});

		this.view.observe(this.model.attributes.items);
	},
	_initializeWorker: function(item) {
		this.model.attributes.add(item);
	}


});



