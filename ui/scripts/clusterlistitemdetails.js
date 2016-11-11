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
var detailsModel = new Model('detailsModel',{
	_elements: [],

	// set elements
	set: function(data) {
		this._elements = [];

		//
		// iterate all the items from data to generate the clusterlist
		// achive types
		_.each(data, function(contentType, type){
			// achieve item
			_.each(contentType, function(contentItem, item){
				// achieve detail
				_.each(contentItem, function(contentDetailType, detailType){

					// only are required objects where are defined all the details
					if (_.isObject(contentDetailType)){
						// achive moreDetail
						_.each(contentDetailType,function(contentDetailItem, detailItem){
							locator = detailsModel.attributes.generateLocator(type,item,detailType,detailItem);	
							detailsModel.attributes.add(type, item, detailType, detailItem, locator, contentItem.URL, contentDetailItem);	
						});// end achive detail item
					}// en isObject
				});// end achieve detail type
			});// end achieve item
		});// end achieve types
	
		// set the array of elements
		detailsModel.set(this._elements);

	},

	// createItem element
	createItem: function(type, item, detailtype, detailitem, locator, urlbase, content) {
		return {type: type, item: item, detailtype: detailtype, detailitem: detailitem, url: urlbase, locator: locator, content: content};
	},

	// add
	add: function(type, item, detail, moreDetail, locator, urlbase, content) {
		this._elements.push(this.createItem(type, item, detail, moreDetail, locator, urlbase, content));
	},

	// generateLocator
	generateLocator: function(type,item,detail,moreDetail) {
		return "/"+ type + "/" + item + "/" + detail + "/" + moreDetail;
	}

});

//
// View Object
//-----------------------------------------------------------
var detailsView = new View({
	// place to append the items
	parent: ".clusterlistitemdetails",

	render: function(model) {
		_.each(model, function(item){
			 $('<div/>', {
					class: 'itemdetails',
					id: item.detailtype+'_'+item.detailitem,
					base: item.type+'_'+item.item,
					status: item.content.status,
					text: item.detailitem
			}).hide().appendTo(detailsView.parent);
		});
	},

	// subscribe to model
	observe: function(model) {
		// subscribe
		this.on(model.id, this.id, function(model){
			detailsView.render(model);
		}.bind(this));
	},

	// subscribe to model
	observeClusterlist: function(model) {
		// subscribe
		detailsView.on(model.id, detailsView.id, function(){
			selectedItem = model.attributes.getSelectedItem();
			if (selectedItem.length > 0) {
				detailsView.setShowedElements(selectedItem[0].locator);				
			}
		}.bind(this));
	},
	// manage which elements to show
	setShowedElements: function(locator) {	
		locatorSplitted = locator.split('/');
		//console.log('detailsView::setShowedElements',locatorSplitted);

		$('.itemdetails').each(function(index){
			if ( this.getAttribute('base') == locatorSplitted[1] +'_'+ locatorSplitted[2] ) {
				$(this).toggle(); 
			} else {
				$(this).hide();
			}
		});
	},
	// observer for changes on menu
	observeMenu: function(model) {
		this.on(model.id, this.id, function(){
			$('.itemdetails').hide();
		}.bind(this));
	}

});

//
// Controller Object
//-----------------------------------------------------------
var detailsController = new Controller({
	// menu data management
	model: detailsModel,
	// menu view management
	view: detailsView,

	events: {
		".itemdetails::click": "selectListItem"
	},

	//
	// menu initializes when document is ready
	initialize: function(data) {
		//console.log('detailsController::initialize');

		//
		// set listeners
		//
		// subscribe view to model
		this.view.observe(this.model);
		// subscribe for changes to clusterlist
		this.view.observeClusterlist(clusterlistModel);
		// subscribe for changes to menu
		this.view.observeMenu(menuModel);


		// initialitze data for clusterlist
		detailsModel.attributes.set(data);
	},

	_initializeWorker: function() {

	},

	update: function(data) {
		detailsModel.attributes.set(data);
	},

	selectedItem: function() {
		console.log('detailsController::selectedItem');
	}

});