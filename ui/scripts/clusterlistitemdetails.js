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
// Detail's Model Object
//-----------------------------------------------------------
var detailsModel = new Model('detailsModel',{
	_elements: [],

	// set elements
	set: function(data) {
		this._elements = [];
		this._locator = locatorModel.attributes.get();
		this._selected = false;

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
							locator = detailsModel.attributes.generateLocator(type,item,detailItem);
								detailsModel.attributes.add(type, item, detailType, detailItem, locator, this._selected, false, contentItem.URL, contentDetailItem);
						});// end achive detail item
					}// end isObject
				});// end achieve detail type
			});// end achieve item
		});// end achieve types
	
		// set the array of elements
		// detailsModel.set(this._elements);

		// if there was any selected clusterlist item, this is setted again
		if ( _.size(this._locator) > 0 && _.findWhere(this._elements, {locator: this._locator}) != undefined ) {
		 	clusterlistModel.attributes.setSelected(this._locator);
		}

		// set the element to show
		if (clusterlistModel.attributes.getSelected().length > 0 ) {
			detailsModel.attributes.show(clusterlistModel.attributes.getSelected()[0].locator);
		}
	},

	// createItem element
	createItem: function(type, item, detailtype, detailitem, locator, selected, show, urlbase, content) {
		return {type: type, item: item, detailtype: detailtype, detailitem: detailitem, selected: selected, show: show, url: urlbase, locator: locator, content: content};
	},

	// add
	add: function(type, item, detail, moreDetail, locator, selected, show, urlbase, content) {
		this._elements.push(this.createItem(type, item, detail, moreDetail, locator, selected, show, urlbase, content));
	},

	show: function(locator) {
		console.log('detailsModel::show',locator);
		// for each detailModel's element
		_.each(this._elements, function(item){
			// console.log('detailsModel::show',locator);
			// console.log('detailsModel::show',detailsModel.attributes.getBase(item.locator));
			// show items which locator base is the same as the selected item on cluster list.
			if ( locator == detailsModel.attributes.getBase(item.locator)) {
				item.show = true;						
				console.log('detailsModel::show',item);
			} else {
				item.show = false;
			}
		});

		detailsModel.set(this._elements);
	},

	getShowed: function() {
		return _.where(this._elements, {show: true});
	}

	// generateLocator
	generateLocator: function(type,item,detail) {
		return '/'+ type + '/' + item + '/' + detail;
	},

	// return to which type of item belongs the detail
	getBase: function(locator) {
		locatorSplitted = locator.split('/');
		if ( locatorSplitted.length > 1 )
			return '/' + locatorSplitted[1] +'/'+ locatorSplitted[2]
		else
			return null;
	},

	// return the detail's id
	getDetail: function(locator) {
		locatorSplitted = locator.split('/');
		if ( locatorSplitted.length > 2 )
			return locatorSplitted[3]
		else
			return null;
	},

	getSelected: function(){
		return _.where(this._elements, {selected: true});
	},

	setSelected: function(locator) {
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
		detailsModel.set(this._elements);
	},

	// subscriptions
	observe: function(model, f) {
		clusterlistModel.on(model.id, this.id, function(){f(model);}.bind(this));
	},
	// observe cluster list
	observeClusterlist: function(model) {
			// get selected item from clusterlist
			selected = model.attributes.getSelected();
	 		// console.log('detailsModel::observeClusterlist',selected[0], selected.length);
			// if selected item exist
			if (selected.length > 0) {
				detailsModel.attributes.show(selected[0].locator);
			}
	}


});

//
// Detail's View Object
//-----------------------------------------------------------
var detailsView = new View({
	// place to append the items
	parent: ".clusterlistitemdetails",
	// titles
	_titles: {
		"nodes": "services",
		"services": "nodes"
	},

	render: function(model) {
		// console.log('detailsView::render', model);
		$(detailsView.parent).empty();

			$('<div/>', {
				class: 'itemdetailstitle',
				id: 'title',
				text: ''
			}).hide().appendTo(detailsView.parent);

		_.each(model, function(item){
			 $('<div/>', {
					class: 'itemdetails',
					id: item.locator,
					status: item.content.status,
					text: item.detailitem
			}).hide().appendTo(detailsView.parent);
		});
		// set elements to show
		detailsView.setShowedElements(locatorModel.attributes.get());
		
		// set events for details
		detailsController.setEvents();
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
			selectedItem = model.attributes.getSelected();
			if (selectedItem.length > 0) {
				// console.log('detailsView::observeClusterlist',selectedItem[0]);
				detailsView.setShowedElements(selectedItem[0].locator);				
			}
		}.bind(this));
	},

	// set title to details list
	setTitle: function(title) {
		// if ( $('.itemdetailstitle').text().length > 0) {
		// 	$('.itemdetailstitle').hide();
		// 	$('.itemdetailstitle').text('');
		// } else {
		// 	console.log('detailsView::setTitle',title);
		// 	$('.itemdetailstitle').text(title);
		// 	$('.itemdetailstitle').show();
		// }
		$('.itemdetailstitle').text(title);

	},
	// manage which elements to show
	setShowedElements: function(locator) {
		
		locatorSplitted = locator.split('/');
		//console.log('detailsView::setShowedElements',locatorSplitted);
		base = detailsModel.attributes.getBase(locator);
		selected = detailsModel.attributes.getSelected();

		// set title
		// if (locatorSplitted.length > 1) {
		// 	console.log('detailsView::setShowedElements',locatorSplitted[1]);
		// 	detailsView.setTitle(this._titles[locatorSplitted[1]].toUpperCase());
		// }
		
		$('.itemdetails').each(function(){
			if ( detailsModel.attributes.getBase(this.getAttribute('id')) == base ) {
				$(this).show();
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
		".itemdetails::click": "selected"
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
		// this.view.observeMenu(menuModel);
		this.model.attributes.observe(clusterlistModel, this.model.attributes.observeClusterlist);

		// initialitze data for clusterlist
		detailsModel.attributes.set(data);

		// set events for details
		//detailsController._setEvents();
	},

	_initializeWorker: function() {

	},

	update: function(data) {
		detailsModel.attributes.set(data);
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

	selected: function() {
		detailsController.model.attributes.setSelected(this.getAttribute('id'));
	}

});