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
	
	add: function(i) {
		console.log(item);
		menuModel.set({ item: i});
	},

	getSelectedItem: function() {
		console.log('getSelectedItem');
	}	
});

//
// View Object
//-----------------------------------------------------------
var menuView = new View();

//
// Controller Object
//-----------------------------------------------------------
var menuController = new Controller({
	model: menuModel,
	view: menuView,

	initialize: function(data){
		_.each(data,function(content, item){
			//console.log(item);
			this.model.add(item);
		});
	}


});



