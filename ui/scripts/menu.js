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
var menuModel = new Model();

menuModel.prototype.getSelectedItem = function() {
	return $(this.attributes[selected=true]);
};

menuModel.prototype.setSelectedItem = function() {
	
};

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
		console.log("hola");
	}
});



